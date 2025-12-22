package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CtyunExpressConnectRoute struct {
	meta *common.CtyunMetadata
}

func (c *CtyunExpressConnectRoute) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_route"
}

func (c *CtyunExpressConnectRoute) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta

}

func NewCtyunExpressConnectRoute() resource.Resource {
	return &CtyunExpressConnectRoute{}
}

func (c *CtyunExpressConnectRoute) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			title := "导入失败：" + err.Error()
			detail := "导入命令：terraform import [配置标识].[导入配置名称] [ID],[ecId],[cgwId],[rtbId],[nextHopId]"
			response.Diagnostics.AddError(title, detail)
		}
	}()
	var config CtyunExpressConnectRouteConfig

	var ID, ecId, cgwId, rtbId, nextHopId string

	err = terraform_extend.Split(request.ID, &ID, &ecId, &cgwId, &rtbId, &nextHopId)
	if err != nil {
		return
	}

	if ID == "" {
		err = fmt.Errorf("ID不能为空")
		return
	}
	if ecId == "" {
		err = fmt.Errorf("ecId不能为空")
		return
	}
	if cgwId == "" {
		err = fmt.Errorf("cgwId不能为空")
		return
	}
	if rtbId == "" {
		err = fmt.Errorf("rtbId不能为空")
		return
	}
	if nextHopId == "" {
		err = fmt.Errorf("nextHopId不能为空")
		return
	}

	config.ID = types.StringValue(ID)
	config.EcID = types.StringValue(ecId)
	config.CgwID = types.StringValue(cgwId)
	config.RtbID = types.StringValue(rtbId)
	config.NextHopID = types.StringValue(nextHopId)

	err = c.getAndMerge(ctx, &config)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, config)...)
}

func (c *CtyunExpressConnectRoute) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10132372",
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rtb_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表id",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			//"route_type": schema.StringAttribute{
			//	Optional:    true,
			//	Computed:    true,
			//	Default:     stringdefault.StaticString(business.ECRouteTypeCustom),
			//	Description: "路由类型，取值范围：auto-自动学习，custom-自定义",
			//	Validators: []validator.String{
			//		stringvalidator.OneOf(business.ECRouteTypeAuto, business.ECRouteTypeCustom),
			//	},
			//	PlanModifiers: []planmodifier.String{
			//		stringplanmodifier.RequiresReplace(),
			//	},
			//},
			"cidr": schema.StringAttribute{
				Required:    true,
				Description: "子网信息",
				Validators: []validator.String{
					validator2.Cidr(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"next_hop_type": schema.StringAttribute{
				Optional:    true,
				Description: "下一跳实例的类型，如不是黑洞路由则必填。取值范围：vpc-虚拟私有云，cda-云专线，vpn-vpn网关，cross-跨域连接",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcNextHopTypeVPC, business.EcNextHopCDA, business.EcNextHopVPN, business.EcNextHopBlackCross),
					validator2.ConflictsWithEqualString(
						path.MatchRoot("is_black_hole_route"),
						types.BoolValue(true),
					),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("is_black_hole_route"),
						types.BoolValue(false),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"next_hop_id": schema.StringAttribute{
				Optional:    true,
				Description: "目的实例ID/跨域连接ID，如不是黑洞路由则必填",
				Validators: []validator.String{
					validator2.ConflictsWithEqualString(
						path.MatchRoot("is_black_hole_route"),
						types.BoolValue(true),
					),
					validator2.AlsoRequiresEqualString(
						path.MatchRoot("is_black_hole_route"),
						types.BoolValue(false),
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "路由描述信息",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validator2.Desc(),
				},
			},
			"ip_version": schema.StringAttribute{
				Required:    true,
				Description: "子网类型。取值范围:ipv4和ipv6",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcIpVersionIpv4, business.EcIpVersionIpv6),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_black_hole_route": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "是否是黑洞路由, 如果选择true，next_hop_type、next_hop_id字段可不填写",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "路由规则id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "创建时间，为UTC格式",
			},
		},
	}
}

func (c *CtyunExpressConnectRoute) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var plan CtyunExpressConnectRouteConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	err = c.create(ctx, &plan)
	if err != nil {
		return
	}
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunExpressConnectRoute) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunExpressConnectRouteConfig
	// 读取state状态
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		response.State.RemoveResource(ctx)
		err = nil
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *CtyunExpressConnectRoute) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	return
}

func (c *CtyunExpressConnectRoute) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	// 获取state
	var config CtyunExpressConnectRouteConfig
	response.Diagnostics.Append(request.State.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	err = c.delete(ctx, config)
	if err != nil {
		return
	}
}

func (c *CtyunExpressConnectRoute) create(ctx context.Context, config *CtyunExpressConnectRouteConfig) error {
	params := &ec.EcEcCreateRouteRequest{
		EcID:             config.EcID.ValueString(),
		CgwID:            config.CgwID.ValueString(),
		RtbID:            config.RtbID.ValueString(),
		RouteType:        "2",
		RouteCIDR:        config.CIDR.ValueString(),
		IPVersion:        business.EcIpVersionMap[config.IPVersion.ValueString()],
		IsBlackholeRoute: config.IsBlackHoleRoute.ValueBoolPointer(),
	}
	if !config.IsBlackHoleRoute.ValueBool() {
		nextHopType := business.EcNextHopTypeMap[config.NextHopType.ValueString()]
		params.NexthopType = &nextHopType
		params.NexthopID = config.NextHopID.ValueStringPointer()
	}
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params.RouteDescription = config.Description.ValueStringPointer()
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcCreateRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("创建云间高速路由失败，接口返回nil，请联系研发确认问题原因！")
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	config.ID = types.StringValue(*resp.ReturnObj.RouteID)
	return nil

}

func (c *CtyunExpressConnectRoute) getAndMerge(ctx context.Context, config *CtyunExpressConnectRouteConfig) error {
	params := &ec.EcEcListRouteRequest{
		EcID:  config.EcID.ValueString(),
		CgwID: config.CgwID.ValueString(),
		RtbID: config.RtbID.ValueString(),
		//RouteID: config.ID.ValueStringPointer(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcListRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("查询云间高速(id=%s)的路由表(id=%s)中路由(id=%s)详情失败，接口返回nil，请联系研发确认问题原因！",
			config.EcID.ValueString(), config.CgwID.ValueString(), config.ID)
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return err
	}
	for _, routeObj := range resp.ReturnObj.Results {
		if *routeObj.RouteID == config.ID.ValueString() {
			//config.RouteType = types.StringValue(business.EcRouteTypeRevMap[*routeObj.RouteType])
			config.CIDR = types.StringValue(*routeObj.RouteCIDR)
			if !config.IsBlackHoleRoute.ValueBool() {
				config.NextHopType = types.StringValue(business.EcNextHopTypeRevMap[*routeObj.NexthopType])
			} else {
				config.NextHopType = types.StringNull()
				config.NextHopID = types.StringNull()
			}
			config.IPVersion = types.StringValue(business.EcIpVersionRevMap[*routeObj.IPVersion])
			config.Description = types.StringValue(*routeObj.RouteDescription)
			config.CreateTime = types.StringValue(*routeObj.CreateDate)
		}
	}
	return nil
}

func (c *CtyunExpressConnectRoute) delete(ctx context.Context, config CtyunExpressConnectRouteConfig) error {
	params := &ec.EcEcDeleteRouteRequest{
		RouteID: config.ID.ValueString(),
		RtbID:   config.RtbID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcDeleteRouteApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return err
	} else if resp == nil {
		err = fmt.Errorf("删除路由(id=%s)失败，接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return err
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return err
	}
	return nil
}

type CtyunExpressConnectRouteConfig struct {
	EcID             types.String `tfsdk:"ec_id"`
	CgwID            types.String `tfsdk:"cgw_id"`
	RtbID            types.String `tfsdk:"rtb_id"`
	CIDR             types.String `tfsdk:"cidr"`
	NextHopType      types.String `tfsdk:"next_hop_type"`
	NextHopID        types.String `tfsdk:"next_hop_id"`
	Description      types.String `tfsdk:"description"`
	IPVersion        types.String `tfsdk:"ip_version"`
	IsBlackHoleRoute types.Bool   `tfsdk:"is_black_hole_route"`
	ID               types.String `tfsdk:"id"`
	CreateTime       types.String `tfsdk:"create_time"`
	//RouteType        types.String `tfsdk:"route_type"`

}
