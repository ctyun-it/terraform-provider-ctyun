package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
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
	_ resource.Resource                = &CtyunMongodbAssociationEip{}
	_ resource.ResourceWithConfigure   = &CtyunMongodbAssociationEip{}
	_ resource.ResourceWithImportState = &CtyunMongodbAssociationEip{}
)

func (c *CtyunMongodbAssociationEip) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_mongodb_association_eip"
}
func NewCtyunMongodbAssociationEip() resource.Resource {
	return &CtyunMongodbAssociationEip{}
}

type CtyunMongodbAssociationEip struct {
	meta       *common.CtyunMetadata
	eipService *business.EipService
}

func (c *CtyunMongodbAssociationEip) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (c *CtyunMongodbAssociationEip) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.eipService = business.NewEipService(c.meta)
}

func (c *CtyunMongodbAssociationEip) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10034467/10183412`,
		Attributes: map[string]schema.Attribute{
			"eip_id": schema.StringAttribute{
				Required:    true,
				Description: "弹性id",
				Validators: []validator.String{
					validator2.EipValidate(),
				},
			},
			"inst_id": schema.StringAttribute{
				Required:    true,
				Description: "实例id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"host_ip": schema.StringAttribute{
				Required:    true,
				Description: "主机ip",
				Validators: []validator.String{
					validator2.Ip(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
				Validators: []validator.String{
					validator2.Project(),
				},
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池Id",
				Default:     defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"eip_address": schema.StringAttribute{
				Computed:    true,
				Description: "弹性ip对应的地址",
			},
		},
	}
}

func (c *CtyunMongodbAssociationEip) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan MongodbAssociationEipConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 实例绑定弹性IP
	err = c.MongodbBindEip(ctx, &plan)
	if err != nil {
		return
	}
	// 轮询查看绑定状态
	err = c.BindLoop(ctx, &plan)
	if err != nil {
		return
	}
	// 查询实例详情，确认是否绑定成功
	err = c.getAndMergeBindEip(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (c *CtyunMongodbAssociationEip) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state MongodbAssociationEipConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// 查询远端
	err = c.getAndMergeBindEip(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "is not found") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *CtyunMongodbAssociationEip) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunMongodbAssociationEip) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var state MongodbAssociationEipConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	unbindParams := &mongodb.MongodbUnbindEipRequest{
		EipID:  state.EipID.ValueString(),
		Eip:    state.EipAddress.ValueString(),
		InstID: state.InstID.ValueString(),
	}
	unbindHeader := &mongodb.MongodbUnbindEipRequestHeader{}
	if state.ProjectID.ValueString() != "" {
		unbindHeader.ProjectID = state.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbUnbindEipApi.Do(ctx, c.meta.Credential, unbindParams, unbindHeader)
	if err != nil {
		return
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 轮询确定解绑成功
	err = c.unBindLoop(ctx, &state)
	if err != nil {
		return
	}
}

func (c *CtyunMongodbAssociationEip) MongodbBindEip(ctx context.Context, config *MongodbAssociationEipConfig) (err error) {

	eip, err := c.eipService.GetEipAddressByEipID(ctx, config.EipID.ValueString(), config.RegionID.ValueString())
	if err != nil {
		return err
	}
	config.EipAddress = types.StringValue(*eip.EipAddress)

	bindParams := &mongodb.MongodbBindEipRequest{
		EipID:  config.EipID.ValueString(),
		Eip:    config.EipAddress.ValueString(),
		InstID: config.InstID.ValueString(),
		HostIp: config.HostIP.ValueString(),
	}
	bindHeader := &mongodb.MongodbBindEipRequestHeader{}
	if config.ProjectID.ValueString() != "" {
		bindHeader.ProjectID = config.ProjectID.ValueStringPointer()
	}
	resp, err := c.meta.Apis.SdkMongodbApis.MongodbBindEipApi.Do(ctx, c.meta.Credential, bindParams, bindHeader)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	return
}

func (c *CtyunMongodbAssociationEip) BindLoop(ctx context.Context, config *MongodbAssociationEipConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: config.InstID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				ProjectID: nil,
				RegionID:  config.RegionID.ValueString(),
			}
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			nodeInfoVos := resp.ReturnObj.NodeInfoVOS
			for _, vos := range nodeInfoVos {
				if vos.OuterElasticIpId == config.EipID.ValueString() && vos.ElasticIp == config.EipAddress.ValueString() {
					return false
				}
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，eip仍未绑定/解绑成功！")
	}
	return
}

func (c *CtyunMongodbAssociationEip) unBindLoop(ctx context.Context, config *MongodbAssociationEipConfig, loopCount ...int) (err error) {
	count := 60
	if len(loopCount) > 0 {
		count = loopCount[0]
	}
	retryer, err := business.NewRetryer(time.Second*10, count)
	if err != nil {
		return
	}
	result := retryer.Start(
		func(currentTime int) bool {
			detailParams := &mongodb.MongodbQueryDetailRequest{
				ProdInstId: config.InstID.ValueString(),
			}
			detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
				ProjectID: nil,
				RegionID:  config.RegionID.ValueString(),
			}
			resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
			if err2 != nil {
				err = err2
				return false
			} else if resp.StatusCode != 800 {
				err = fmt.Errorf("API return error. Message: %s", *resp.Message)
				return false
			} else if resp.ReturnObj == nil {
				err = common.InvalidReturnObjError
				return false
			}
			nodeInfoVos := resp.ReturnObj.NodeInfoVOS
			for _, vos := range nodeInfoVos {
				if vos.OuterElasticIpId == "" && vos.ElasticIp == "" {
					return false
				}
			}
			return true
		})
	if result.ReturnReason == business.ReachMaxLoopTime {
		return errors.New("轮询已达最大次数，eip仍未绑定/解绑成功！")
	}
	return
}

func (c *CtyunMongodbAssociationEip) getAndMergeBindEip(ctx context.Context, config *MongodbAssociationEipConfig) (err error) {
	detailParams := &mongodb.MongodbQueryDetailRequest{
		ProdInstId: config.InstID.ValueString(),
	}
	detailHeader := &mongodb.MongodbQueryDetailRequestHeaders{
		ProjectID: nil,
		RegionID:  config.RegionID.ValueString(),
	}
	resp, err2 := c.meta.Apis.SdkMongodbApis.MongodbQueryDetailApi.Do(ctx, c.meta.Credential, detailParams, detailHeader)
	if err2 != nil {
		err = err2
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	nodeinfoVos := resp.ReturnObj.NodeInfoVOS[0]
	config.EipID = types.StringValue(nodeinfoVos.OuterElasticIpId)
	config.EipAddress = types.StringValue(nodeinfoVos.ElasticIp)
	config.HostIP = types.StringValue(resp.ReturnObj.Host)
	return
}

type MongodbAssociationEipConfig struct {
	EipID      types.String `tfsdk:"eip_id"`      // 弹性ip id
	InstID     types.String `tfsdk:"inst_id"`     // 实例id
	HostIP     types.String `tfsdk:"host_ip"`     // 主机ip
	ProjectID  types.String `tfsdk:"project_id"`  // 项目id
	RegionID   types.String `tfsdk:"region_id"`   // 资源池id
	EipAddress types.String `tfsdk:"eip_address"` // eip地址
}
