package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &ctyunVpcRouteTableRule{}
	_ resource.ResourceWithConfigure   = &ctyunVpcRouteTableRule{}
	_ resource.ResourceWithImportState = &ctyunVpcRouteTableRule{}
)

type ctyunVpcRouteTableRule struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcRouteTableRule() resource.Resource {
	return &ctyunVpcRouteTableRule{}
}

func (c *ctyunVpcRouteTableRule) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_route_table_rule"
}

type CtyunVpcRouteTableRuleConfig struct {
	ID           types.String `tfsdk:"id"`
	RuleID       types.String `tfsdk:"rule_id"`
	RegionID     types.String `tfsdk:"region_id"`
	RouteTableID types.String `tfsdk:"route_table_id"`
	NextHopID    types.String `tfsdk:"next_hop_id"`
	NextHopType  types.String `tfsdk:"next_hop_type"`
	Destination  types.String `tfsdk:"destination"`
	IpVersion    types.Int32  `tfsdk:"ip_version"`
	Description  types.String `tfsdk:"description"`
}

func (c *ctyunVpcRouteTableRule) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10171000`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "规则id",
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
			"route_table_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"next_hop_id": schema.StringAttribute{
				Required:    true,
				Description: "下一跳设备id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"next_hop_type": schema.StringAttribute{
				Required:    true,
				Description: "下一跳设备类型，支持vpcpeering、havip、bm、vm、natgw、igw、igw6、dc、ticc、vpngw、enic",
				Validators: []validator.String{
					stringvalidator.OneOf("vpcpeering", "havip", "bm", "vm", "natgw", "igw", "igw6", "dc", "ticc", "vpngw", "enic"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "无类别域间路由，例如：192.168.0.1/32",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"ip_version": schema.Int32Attribute{
				Required:    true,
				Description: "4标识ipv4,6标识ipv6",
				Validators: []validator.Int32{
					int32validator.OneOf(4, 6),
				},
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "规则描述，支持更新",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
		},
	}
}

func (c *ctyunVpcRouteTableRule) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunVpcRouteTableRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建
	ruleID, err := c.create(ctx, plan)
	if err != nil {
		return
	}
	plan.RuleID = types.StringValue(ruleID)
	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunVpcRouteTableRule) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcRouteTableRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "未找到") {
			response.State.RemoveResource(ctx)
			err = nil
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpcRouteTableRule) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	// tf文件中的
	var plan CtyunVpcRouteTableRuleConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// state中的
	var state CtyunVpcRouteTableRuleConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 更新
	err = c.update(ctx, plan, state)
	if err != nil {
		return
	}
	// 查询远端信息
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunVpcRouteTableRule) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunVpcRouteTableRuleConfig
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

func (c *ctyunVpcRouteTableRule) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// 导入命令：terraform import [配置标识].[导入配置名称] [ruleID],[routeTableID],[regionID]
func (c *ctyunVpcRouteTableRule) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunVpcRouteTableRuleConfig
	var ruleID, routeTableID, regionID string
	err = terraform_extend.Split(request.ID, &ruleID, &routeTableID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.RouteTableID = types.StringValue(routeTableID)
	cfg.RuleID = types.StringValue(ruleID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

// create 创建
func (c *ctyunVpcRouteTableRule) create(ctx context.Context, plan CtyunVpcRouteTableRuleConfig) (ruleID string, err error) {

	routeTableID, regionID := plan.RouteTableID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcCreateRouteRuleRequest{
		RegionID:     regionID,
		RouteTableID: routeTableID,
		NextHopID:    plan.NextHopID.ValueString(),
		NextHopType:  plan.NextHopType.ValueString(),
		Destination:  plan.Destination.ValueString(),
		IpVersion:    plan.IpVersion.ValueInt32(),
		Description:  plan.Description.ValueStringPointer(),
	}

	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcCreateRouteRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	ruleID = *resp.ReturnObj.RouteRuleID
	return
}

// getAndMerge 从远端查询
func (c *ctyunVpcRouteTableRule) getAndMerge(ctx context.Context, plan *CtyunVpcRouteTableRuleConfig) (err error) {
	ruleID, routeTableID, regionID := plan.RuleID.ValueString(), plan.RouteTableID.ValueString(), plan.RegionID.ValueString()
	var rules []*ctvpc.CtvpcNewRouteRulesListReturnObjRouteRulesResponse
	pageNo := 1
	for {
		params := &ctvpc.CtvpcNewRouteRulesListRequest{
			RegionID:     regionID,
			RouteTableID: routeTableID,
			PageSize:     50,
			PageNo:       int32(pageNo),
		}
		var resp *ctvpc.CtvpcNewRouteRulesListResponse
		resp, err = c.meta.Apis.SdkCtVpcApis.CtvpcNewRouteRulesListApi.Do(ctx, c.meta.SdkCredential, params)
		if err != nil {
			return
		} else if resp.StatusCode == common.ErrorStatusCode {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
			return
		} else if resp.ReturnObj == nil {
			err = common.InvalidReturnObjError
			return
		}
		rules = append(rules, resp.ReturnObj.RouteRules...)
		if int32(pageNo) >= resp.ReturnObj.TotalPage {
			break
		}
		pageNo++
	}
	if len(rules) == 0 {
		err = common.InvalidReturnObjResultsError
		return
	}
	var exist bool
	for _, r := range rules {
		if utils.SecString(r.RouteRuleID) == ruleID {
			exist = true
			plan.NextHopID = utils.SecStringValue(r.NextHopID)
			plan.NextHopType = utils.SecStringValue(r.NextHopType)
			plan.Description = utils.SecStringValue(r.Description)
			plan.Destination = utils.SecStringValue(r.Destination)
			plan.IpVersion = types.Int32Value(r.IpVersion)
			plan.ID = plan.RuleID
		}
	}
	if !exist {
		err = fmt.Errorf("未找到路由表 %s 下的规则 %s", routeTableID, ruleID)
	}
	return
}

// update 更新
func (c *ctyunVpcRouteTableRule) update(ctx context.Context, plan, state CtyunVpcRouteTableRuleConfig) (err error) {
	if plan.Description.Equal(state.Description) {
		return
	}
	ruleID, regionID, description := state.RuleID.ValueString(), state.RegionID.ValueString(), plan.Description.ValueString()
	params := &ctvpc.CtvpcModifyRouteRuleRequest{
		RegionID:    regionID,
		RouteRuleID: ruleID,
		Description: &description,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcModifyRouteRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}

// delete 删除
func (c *ctyunVpcRouteTableRule) delete(ctx context.Context, plan CtyunVpcRouteTableRuleConfig) (err error) {
	ruleID, regionID := plan.RuleID.ValueString(), plan.RegionID.ValueString()
	params := &ctvpc.CtvpcDeleteRouteRuleRequest{
		RegionID:    regionID,
		RouteRuleID: ruleID,
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDeleteRouteRuleApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	}
	return
}
