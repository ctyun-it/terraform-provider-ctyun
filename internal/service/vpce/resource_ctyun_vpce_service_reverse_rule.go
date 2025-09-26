package vpce

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ctyunVpceServiceReverseRule{}
	_ resource.ResourceWithConfigure   = &ctyunVpceServiceReverseRule{}
	_ resource.ResourceWithImportState = &ctyunVpceServiceReverseRule{}
)

type ctyunVpceServiceReverseRule struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServiceReverseRule() resource.Resource {
	return &ctyunVpceServiceReverseRule{}
}

func (c *ctyunVpceServiceReverseRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service_reverse_rule"
}

type CtyunVpceServiceReverseRuleConfig struct {
	ID                types.String `tfsdk:"id"`
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	EndpointID        types.String `tfsdk:"endpoint_id"`
	RegionID          types.String `tfsdk:"region_id"`
	TransitIP         types.String `tfsdk:"transit_ip"`
	TransitPort       types.Int32  `tfsdk:"transit_port"`
	TargetIP          types.String `tfsdk:"target_ip"`
	TargetPort        types.Int32  `tfsdk:"target_port"`
	Protocol          types.String `tfsdk:"protocol"`
}

func (c *ctyunVpceServiceReverseRule) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10042658/10048506`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "规则ID",
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
			"endpoint_service_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点服务ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"endpoint_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"transit_ip": schema.StringAttribute{
				Required:    true,
				Description: "中转IP地址",
				Validators: []validator.String{
					validator2.Ip(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"transit_port": schema.Int32Attribute{
				Required:    true,
				Description: "中转端口(1-65535)",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"target_ip": schema.StringAttribute{
				Required:    true,
				Description: "目标IP地址",
				Validators: []validator.String{
					validator2.Ip(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_port": schema.Int32Attribute{
				Required:    true,
				Description: "目标端口(1-65535)",
				Validators: []validator.Int32{
					int32validator.Between(1, 65535),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "协议，TCP:TCP协议,UDP:UDP协议",
				Validators: []validator.String{
					stringvalidator.OneOf("TCP", "UDP"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (c *ctyunVpceServiceReverseRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpceServiceReverseRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	ruleID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(ruleID)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpceServiceReverseRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceReverseRuleConfig
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

func (c *ctyunVpceServiceReverseRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunVpceServiceReverseRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpceServiceReverseRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
}

func (c *ctyunVpceServiceReverseRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [id],[endpointServiceID],[regionID]
func (c *ctyunVpceServiceReverseRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpceServiceReverseRuleConfig
	var id, endpointServiceID, regionID string
	err = terraform_extend.Split(request.ID, &id, &endpointServiceID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.EndpointServiceID = types.StringValue(endpointServiceID)
	cfg.ID = types.StringValue(id)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpceServiceReverseRule) create(ctx context.Context, plan CtyunVpceServiceReverseRuleConfig) (ruleID string, err error) {
	params := &ctvpc.CtvpcCreateEndpointServiceReverseRuleRequest{
		ClientToken:       uuid.NewString(),
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		EndpointID:        plan.EndpointID.ValueString(),
		TransitIPAddress:  plan.TransitIP.ValueString(),
		TransitPort:       plan.TransitPort.ValueInt32(),
		Protocol:          plan.Protocol.ValueString(),
		TargetIPAddress:   plan.TargetIP.ValueString(),
		TargetPort:        plan.TargetPort.ValueInt32(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateEndpointServiceReverseRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ruleID = resp.ReturnObj.ReverseRuleID
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpceServiceReverseRule) getAndMerge(ctx context.Context, plan *CtyunVpceServiceReverseRuleConfig) (err error) {
	params := &ctvpc.CtvpcListEndpointServiceReverseRuleRequest{
		RegionID:          plan.RegionID.ValueString(),
		EndpointServiceID: plan.EndpointServiceID.ValueString(),
		PageNo:            1,
		PageSize:          50,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListEndpointServiceReverseRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var exist bool
	for _, rule := range resp.ReturnObj.ReverseRules {
		if utils.SecStringValue(rule.ID) == plan.ID {
			plan.EndpointID = utils.SecStringValue(rule.EndpointID)
			plan.TargetIP = utils.SecStringValue(rule.TargetIPAddress)
			plan.TransitIP = utils.SecStringValue(rule.TransitIPAddress)
			plan.TransitPort = types.Int32Value(rule.TransitPort)
			plan.TargetPort = types.Int32Value(rule.TargetPort)
			plan.Protocol = utils.SecStringValue(rule.Protocol)
			exist = true
		}
	}
	if !exist {
		err = common.InvalidReturnObjResultsError
		return
	}
	return
}

// delete 删除
func (c *ctyunVpceServiceReverseRule) delete(ctx context.Context, plan CtyunVpceServiceReverseRuleConfig) (err error) {
	params := &ctvpc.CtvpcDeleteEndpointServiceReverseRuleRequest{
		RegionID:      plan.RegionID.ValueString(),
		ReverseRuleID: plan.ID.ValueString(),
		ClientToken:   uuid.NewString(),
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteEndpointServiceReverseRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
