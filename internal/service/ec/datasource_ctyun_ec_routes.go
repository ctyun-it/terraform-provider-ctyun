package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunExpressConnectionRoutes{}
	_ datasource.DataSourceWithConfigure = &CtyunExpressConnectionRoutes{}
)

type CtyunExpressConnectionRoutes struct {
	meta *common.CtyunMetadata
}

func NewCtyunExpressConnectionRoutes() datasource.DataSource {
	return &CtyunExpressConnectionRoutes{}
}
func (c *CtyunExpressConnectionRoutes) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunExpressConnectionRoutes) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_routes"
}

func (c *CtyunExpressConnectionRoutes) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10132372",
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cgw_id": schema.StringAttribute{
				Required:    true,
				Description: "云网关ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rtb_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "路由ID，精确查询",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊匹配CIDR",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "运行状态（运行中-running，未启用-stop）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcRouteStatusRunning, business.EcRouteStatusStop),
				},
			},
			"route_type": schema.StringAttribute{
				Optional:    true,
				Description: "路由类型（auto-自动学习，default-自定义（默认））",
				Validators: []validator.String{
					stringvalidator.OneOf("auto", "default"),
				},
			},
			"next_hop_type": schema.StringAttribute{
				Optional:    true,
				Description: "下一跳实例类型（取值范围：vpc-虚拟私有云，cda-云专线，vpn-vpn网关，cross-跨域连接）",
				Validators: []validator.String{
					stringvalidator.OneOf(business.EcNextHopTypeVPC, business.EcNextHopCDA, business.EcNextHopVPN, business.EcNextHopBlackCross),
				},
			},
			"routes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速实例ID",
						},
						"cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关ID",
						},
						"rtb_id": schema.StringAttribute{
							Computed:    true,
							Description: "路由表ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "路由ID",
						},
						"route_type": schema.StringAttribute{
							Computed:    true,
							Description: "路由类型（1：自动学习，2：自定义）",
						},
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "子网信息（CIDR格式）",
						},
						"next_hop_type": schema.StringAttribute{
							Computed:    true,
							Description: "下一跳实例类型",
						},
						"next_hop_id": schema.StringAttribute{
							Computed:    true,
							Description: "目的实例ID",
						},
						"ip_version": schema.StringAttribute{
							Computed:    true,
							Description: "子网类型（1：IPv4，2：IPv6）",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "路由描述信息",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "运行状态（1：正常，2：异常）",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
					},
				},
				Description: "云专线路由列表",
			},
		},
	}
}

func (c *CtyunExpressConnectionRoutes) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunExpressConnectRoutesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	params := &ec.EcEcListRouteRequest{
		EcID:  config.EcID.ValueString(),
		CgwID: config.CgwID.ValueString(),
		RtbID: config.RtbID.ValueString(),
	}
	if !config.ID.IsNull() {
		params.RouteID = config.ID.ValueStringPointer()
	}
	if !config.QueryContent.IsNull() {
		params.QueryContent = config.QueryContent.ValueStringPointer()
	}
	if !config.Status.IsNull() {
		status := business.EcRouteStatusMap[config.Status.ValueString()]
		params.Status = &status
	}
	if !config.RouteType.IsNull() {
		routeType := business.EcRouteTypeMap[config.RouteType.ValueString()]
		params.RouteType = &routeType
	}
	if !config.NextHopType.IsNull() {
		nextHopType := business.EcNextHopTypeMap[config.NextHopType.ValueString()]
		params.NexthopType = &nextHopType
	}
	resp, err := c.meta.Apis.SdkEcApis.EcEcListRouteApi.Do(ctx, c.meta.SdkCredential, params)

	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询云间高速(id=%s)的路由表(id=%s)中路由信息失败，接口返回nil，请联系研发确认问题原因！",
			config.EcID.ValueString(), config.CgwID.ValueString())
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var routes []CtyunExpressConnectionRouteModel
	for _, routeItem := range resp.ReturnObj.Results {
		var route CtyunExpressConnectionRouteModel
		route.EcID = types.StringValue(*routeItem.EcID)
		route.CgwID = types.StringValue(*routeItem.CgwID)
		route.RtbID = types.StringValue(*routeItem.RtbID)
		route.ID = types.StringValue(*routeItem.RouteID)
		route.RouteType = types.StringValue(business.EcRouteTypeRevMap[*routeItem.RouteType])
		route.Cidr = types.StringValue(*routeItem.RouteCIDR)
		route.NextHopType = types.StringValue(business.EcNextHopTypeRevMap[*routeItem.NexthopType])
		route.NextHopID = types.StringValue(*routeItem.NexthopID)
		route.IpVersion = types.StringValue(business.EcIpVersionRevMap[*routeItem.IPVersion])
		route.Description = types.StringValue(*routeItem.RouteDescription)
		route.Status = types.StringValue(business.EcRouteStatusRevMap[*routeItem.Status])
		route.CreateTime = types.StringValue(*routeItem.CreateDate)
		routes = append(routes, route)
	}

	config.Routes = routes
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunExpressConnectionRouteModel struct {
	EcID        types.String `tfsdk:"ec_id"`
	CgwID       types.String `tfsdk:"cgw_id"`
	RtbID       types.String `tfsdk:"rtb_id"`
	ID          types.String `tfsdk:"id"`
	RouteType   types.String `tfsdk:"route_type"`
	Cidr        types.String `tfsdk:"cidr"`
	NextHopType types.String `tfsdk:"next_hop_type"`
	NextHopID   types.String `tfsdk:"next_hop_id"`
	IpVersion   types.String `tfsdk:"ip_version"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
	CreateTime  types.String `tfsdk:"create_time"`
}

type CtyunExpressConnectRoutesConfig struct {
	EcID         types.String                       `tfsdk:"ec_id"`
	CgwID        types.String                       `tfsdk:"cgw_id"`
	RtbID        types.String                       `tfsdk:"rtb_id"`
	ID           types.String                       `tfsdk:"id"`
	QueryContent types.String                       `tfsdk:"query_content"`
	Status       types.String                       `tfsdk:"status"`
	RouteType    types.String                       `tfsdk:"route_type"`
	NextHopType  types.String                       `tfsdk:"next_hop_type"`
	Routes       []CtyunExpressConnectionRouteModel `tfsdk:"routes"`
}
