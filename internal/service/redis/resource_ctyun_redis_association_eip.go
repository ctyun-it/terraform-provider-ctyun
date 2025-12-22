package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunRedisAssociationEip{}
	_ resource.ResourceWithConfigure   = &ctyunRedisAssociationEip{}
	_ resource.ResourceWithImportState = &ctyunRedisAssociationEip{}
)

type ctyunRedisAssociationEip struct {
	meta         *common.CtyunMetadata
	redisService *business.RedisService
	eipService   *business.EipService
}

func NewCtyunRedisAssociationEip() resource.Resource {
	return &ctyunRedisAssociationEip{}
}

func (c *ctyunRedisAssociationEip) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_association_eip"
}

type CtyunRedisAssociationEipConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.String `tfsdk:"instance_id"`
	RegionID   types.String `tfsdk:"region_id"`
	EipID      types.String `tfsdk:"eip_id"`
	eipAddress string
}

func (c *ctyunRedisAssociationEip) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10132173`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed:    true,
				Description: "ID",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性IP的ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.EipValidate(),
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "Redis实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(32, 32),
				},
			},
		},
	}
}

func (c *ctyunRedisAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeAssociate(ctx, &plan)
	if err != nil {
		return
	}
	// 创建
	err = c.associate(ctx, plan)
	if err != nil {
		return
	}
	// 创建后检查
	err = c.checkAfterAssociate(ctx, plan)
	if err != nil {
		return
	}

	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "can't find") || strings.Contains(err.Error(), "已退订") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisAssociationEip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *ctyunRedisAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.dissociate(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDissociate(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunRedisAssociationEip) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.redisService = business.NewRedisService(meta)
	c.eipService = business.NewEipService(meta)
}

func (c *ctyunRedisAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [instanceID],[eipID],[regionID]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var cfg CtyunRedisAssociationEipConfig
	var instanceID, eipID, regionID string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) == 1 {
		regionID = c.meta.GetExtraIfEmpty(regionID, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &instanceID, &eipID)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &instanceID, &eipID, &regionID)
		if err != nil {
			return
		}
	}

	if instanceID == "" {
		err = fmt.Errorf("instanceID不能为空")
		return
	}
	if eipID == "" {
		err = fmt.Errorf("eipID不能为空")
		return
	}
	if regionID == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}

	cfg.RegionID = types.StringValue(regionID)
	cfg.InstanceID = types.StringValue(instanceID)
	cfg.EipID = types.StringValue(eipID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// checkBeforeAssociate 绑定前检查
func (c *ctyunRedisAssociationEip) checkBeforeAssociate(ctx context.Context, plan *CtyunRedisAssociationEipConfig) (err error) {
	regionID, instanceID, eipID := plan.RegionID.ValueString(), plan.InstanceID.ValueString(), plan.EipID.ValueString()
	eip, err := c.eipService.GetEipAddressByEipID(ctx, eipID, regionID)
	if err != nil {
		return
	}
	plan.eipAddress = utils.SecString(eip.EipAddress)
	_, err = c.redisService.GetRedisByID(ctx, instanceID, regionID)
	if err != nil {
		return
	}
	return
}

// associate 绑定
func (c *ctyunRedisAssociationEip) associate(ctx context.Context, plan CtyunRedisAssociationEipConfig) (err error) {
	params := &dcs2.Dcs2BindElasticIPRequest{
		RegionId:    plan.RegionID.ValueString(),
		BindObjType: 4,
		ProdInstId:  plan.InstanceID.ValueString(),
		ElasticIp:   plan.eipAddress,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2BindElasticIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// getAndMerge 从远端查询
func (c *ctyunRedisAssociationEip) getAndMerge(ctx context.Context, plan *CtyunRedisAssociationEipConfig) (err error) {
	regionID, instanceID, eipID := plan.RegionID.ValueString(), plan.InstanceID.ValueString(), plan.EipID.ValueString()
	eip, err := c.eipService.GetEipAddressByEipID(ctx, eipID, regionID)
	if err != nil {
		return
	}
	plan.eipAddress = utils.SecString(eip.EipAddress)

	instance, err := c.redisService.GetRedisByID(ctx, instanceID, regionID)
	if err != nil {
		return
	}
	if instance.UserInfo.ElasticIpBind != 1 || instance.UserInfo.ElasticIp != plan.eipAddress {
		err = fmt.Errorf("Redis实例 %s 和弹性IP %s 未绑定", instanceID, eipID)
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", instanceID, eipID))
	return
}

// dissociate 解绑
func (c *ctyunRedisAssociationEip) dissociate(ctx context.Context, plan CtyunRedisAssociationEipConfig) (err error) {
	regionID, instanceID, eipID := plan.RegionID.ValueString(), plan.InstanceID.ValueString(), plan.EipID.ValueString()
	eip, err := c.eipService.GetEipAddressByEipID(ctx, eipID, regionID)
	if err != nil {
		return
	}
	plan.eipAddress = utils.SecString(eip.EipAddress)
	params := &dcs2.Dcs2UnBindElasticIPRequest{
		RegionId:    regionID,
		BindObjType: 4,
		ProdInstId:  instanceID,
		ElasticIp:   plan.eipAddress,
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2UnBindElasticIPApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	return
}

// checkAfterAssociate 绑定后检查
func (c *ctyunRedisAssociationEip) checkAfterAssociate(ctx context.Context, plan CtyunRedisAssociationEipConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeInstancesOverviewReturnObjResponse
			instance, err = c.redisService.GetRedisByID(ctx, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
			if err != nil {
				return false
			}
			if instance.UserInfo.ElasticIpBind != 1 {
				return true
			}

			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("绑定时间过长")
	}
	return
}

// checkAfterDissociate 删除后检查
func (c *ctyunRedisAssociationEip) checkAfterDissociate(ctx context.Context, plan CtyunRedisAssociationEipConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 60)
	retryer.Start(
		func(currentTime int) bool {
			var instance *dcs2.Dcs2DescribeInstancesOverviewReturnObjResponse
			instance, err = c.redisService.GetRedisByID(ctx, plan.InstanceID.ValueString(), plan.RegionID.ValueString())
			if err != nil {
				return false
			}
			if instance.UserInfo.ElasticIpBind != 0 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("解绑时间过长")
	}
	return
}
