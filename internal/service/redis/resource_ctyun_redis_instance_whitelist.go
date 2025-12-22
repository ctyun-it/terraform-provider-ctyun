package redis

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctgdcs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
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
	_ resource.Resource                = &ctyunRedisInstanceWhitelist{}
	_ resource.ResourceWithConfigure   = &ctyunRedisInstanceWhitelist{}
	_ resource.ResourceWithImportState = &ctyunRedisInstanceWhitelist{}
)

type ctyunRedisInstanceWhitelist struct {
	meta *common.CtyunMetadata
}

func NewCtyunRedisInstanceWhitelist() resource.Resource {
	return &ctyunRedisInstanceWhitelist{}
}

func (c *ctyunRedisInstanceWhitelist) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_redis_instance_whitelist"
}

type CtyunRedisInstanceWhitelistConfig struct {
	ID         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	RegionId   types.String `tfsdk:"region_id"`
	Name       types.String `tfsdk:"name"`
	Ip         types.String `tfsdk:"ip"`
}

func (c *ctyunRedisInstanceWhitelist) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10029420/10398174`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "资源唯一标识符",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "实例ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "白名单分组名",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"ip": schema.StringAttribute{
				Optional:    true,
				Description: "白名单列表，可填写IP地址(如192.168.1.1)或IP段(如192.168.1.0/24)，多个IP用英文逗号隔开。当mode=delete时此参数为空 支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *ctyunRedisInstanceWhitelist) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunRedisInstanceWhitelistConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 创建白名单
	err = c.createWhitelist(ctx, plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunRedisInstanceWhitelist) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisInstanceWhitelistConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisInstanceWhitelist) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunRedisInstanceWhitelistConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// state中的
	var state CtyunRedisInstanceWhitelistConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新白名单
	err = c.updateWhitelist(ctx, plan, state)
	if err != nil {
		return
	}
	err = c.checkAfterUpdate(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询备份信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunRedisInstanceWhitelist) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunRedisInstanceWhitelistConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 删除白名单
	err = c.deleteWhitelist(ctx, state)
	if err != nil {
		return
	}
	response.State.RemoveResource(ctx)
}

func (c *ctyunRedisInstanceWhitelist) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *ctyunRedisInstanceWhitelist) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [分组名称],[实例ID],[region_id]"
			response.Diagnostics.AddError(title, detail)
		}
	}()

	var cfg CtyunRedisInstanceWhitelistConfig

	var name, instanceId, regionId string
	// 根据分隔符数量判断是否输入了regionID
	if strings.Count(request.ID, common.ImportSeparator) < 2 {
		regionId = c.meta.GetExtraIfEmpty(regionId, common.ExtraRegionId)
		err = terraform_extend.Split(request.ID, &name, &instanceId)
		if err != nil {
			return
		}
	} else {
		err = terraform_extend.Split(request.ID, &name, &instanceId, &regionId)
		if err != nil {
			return
		}
	}

	if name == "" {
		err = fmt.Errorf("分组名称不能为空")
		return
	}
	if instanceId == "" {
		err = fmt.Errorf("实例ID不能为空")
		return
	}
	if regionId == "" {
		err = fmt.Errorf("regionID不能为空")
		return
	}
	cfg.InstanceId = types.StringValue(instanceId)
	cfg.RegionId = types.StringValue(regionId)
	cfg.Name = types.StringValue(name)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// createWhitelist 创建白名单
func (c *ctyunRedisInstanceWhitelist) createWhitelist(ctx context.Context, plan CtyunRedisInstanceWhitelistConfig) (err error) {
	params := &ctgdcs2.Dcs2ModifySecurityIpsRequest{
		RegionId:   plan.RegionId.ValueString(),
		ProdInstId: plan.InstanceId.ValueString(),
		Group:      plan.Name.ValueString(),
		Mode:       "append",
		Ip:         plan.Ip.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifySecurityIpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// updateWhitelist 更新白名单
func (c *ctyunRedisInstanceWhitelist) updateWhitelist(ctx context.Context, plan, state CtyunRedisInstanceWhitelistConfig) (err error) {
	params := &ctgdcs2.Dcs2ModifySecurityIpsRequest{
		RegionId:   state.RegionId.ValueString(),
		ProdInstId: state.InstanceId.ValueString(),
		Group:      state.Name.ValueString(),
		Mode:       "cover",
		Ip:         plan.Ip.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifySecurityIpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}

	return
}

// checkAfterUpdate 更新后检查  更新需要时间，直接查仍是原来的
func (c *ctyunRedisInstanceWhitelist) checkAfterUpdate(ctx context.Context, plan, state CtyunRedisInstanceWhitelistConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			err := c.getAndMerge(ctx, &state)
			if err != nil {
				return false
			}
			if state.Ip != plan.Ip {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if !executeSuccessFlag {
		err = fmt.Errorf("更新时间过长")
	}
	return
}

// deleteWhitelist 删除白名单
func (c *ctyunRedisInstanceWhitelist) deleteWhitelist(ctx context.Context, state CtyunRedisInstanceWhitelistConfig) (err error) {
	params := &ctgdcs2.Dcs2ModifySecurityIpsRequest{
		RegionId:   state.RegionId.ValueString(),
		ProdInstId: state.InstanceId.ValueString(),
		Group:      state.Name.ValueString(),
		Mode:       "delete",
		Ip:         "",
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2ModifySecurityIpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

// getAndMerge 从远端查询备份信息
func (c *ctyunRedisInstanceWhitelist) getAndMerge(ctx context.Context, state *CtyunRedisInstanceWhitelistConfig) (err error) {
	// 调用API查询白名单信息
	params := &ctgdcs2.Dcs2DescribeSecurityIpsRequest{
		RegionId:   state.RegionId.ValueString(),
		ProdInstId: state.InstanceId.ValueString(),
	}

	resp, err := c.meta.Apis.SdkDcs2Apis.Dcs2DescribeSecurityIpsApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 查找匹配的白名单组
	var whitelistData *ctgdcs2.Dcs2DescribeSecurityIpsReturnObjRowsResponse
	name := state.Name.ValueString()

	for _, whitelist := range resp.ReturnObj.Rows {
		if whitelist.Group == name {
			whitelistData = whitelist
			break
		}
	}

	// 如果找不到对应的白名单组
	if whitelistData == nil {
		err = fmt.Errorf("whitelist group %s not found", name)
		return
	}

	// 更新state中的IP信息
	state.Ip = types.StringValue(whitelistData.Ip)
	state.Name = types.StringValue(whitelistData.Group)

	// 设置ID
	state.ID = types.StringValue(fmt.Sprintf("%s,%s,%s", state.Name.ValueString(), state.InstanceId.ValueString(), state.RegionId.ValueString()))

	return
}
