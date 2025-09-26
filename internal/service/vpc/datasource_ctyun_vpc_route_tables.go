package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunVpcRouteTables{}
	_ datasource.DataSourceWithConfigure = &ctyunVpcRouteTables{}
)

type ctyunVpcRouteTables struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcRouteTables() datasource.DataSource {
	return &ctyunVpcRouteTables{}
}

func (c *ctyunVpcRouteTables) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_route_tables"
}

type CtyunVpcRouteTablesModel struct {
	Name            types.String `tfsdk:"name"`
	VpcID           types.String `tfsdk:"vpc_id"`
	RouteTableID    types.String `tfsdk:"route_table_id"`
	Freezing        types.Bool   `tfsdk:"freezing"`
	RouteRulesCount types.Int32  `tfsdk:"route_rules_count"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	Type            types.Int32  `tfsdk:"type"`
	Origin          types.String `tfsdk:"origin"`
}

type CtyunVpcRouteTablesConfig struct {
	RegionID     types.String `tfsdk:"region_id"`
	VpcID        types.String `tfsdk:"vpc_id"`
	PageNo       types.Int32  `tfsdk:"page_no"`
	PageSize     types.Int32  `tfsdk:"page_size"`
	RouteTableID types.String `tfsdk:"route_table_id"`

	CurrentCount types.Int32                `tfsdk:"current_count"`
	TotalCount   types.Int32                `tfsdk:"total_count"`
	TotalPage    types.Int32                `tfsdk:"total_page"`
	RouteTables  []CtyunVpcRouteTablesModel `tfsdk:"route_tables"`
}

func (c *ctyunVpcRouteTables) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10105078`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "虚拟私有云ID",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10。",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"route_table_id": schema.StringAttribute{
				Optional:    true,
				Description: "路由表ID",
			},
			"current_count": schema.Int32Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数。",
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "总数。",
			},
			"total_page": schema.Int32Attribute{
				Computed:    true,
				Description: "总页数。",
			},
			"route_tables": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "路由表名字",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云id",
						},
						"route_table_id": schema.StringAttribute{
							Computed:    true,
							Description: "路由id",
						},
						"freezing": schema.BoolAttribute{
							Computed:    true,
							Description: "是否冻结",
						},
						"route_rules_count": schema.Int64Attribute{
							Computed:    true,
							Description: "路由表中的路由数",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"type": schema.Int64Attribute{
							Computed:    true,
							Description: "路由表类型:0-子网路由表，2-网关路由表",
						},
						"origin": schema.StringAttribute{
							Computed:    true,
							Description: "路由表来源：default-系统默认;user-用户创建",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpcRouteTables) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpcRouteTablesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)
	// 组装请求体
	params := &ctvpc.CtvpcNewRouteTableListRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	routeTableId := config.RouteTableID.ValueString()
	vpcId := config.VpcID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if routeTableId != "" {
		params.RouteTableID = &routeTableId
	}
	if vpcId != "" {
		params.VpcID = &vpcId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewRouteTableListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		if *resp.ErrorCode == common.OpenapiRouterTableAccessFailed {
			config.RouteTables = []CtyunVpcRouteTablesModel{}
			config.TotalPage = types.Int32Value(0)
			config.TotalCount = types.Int32Value(0)
			config.CurrentCount = types.Int32Value(0)
			response.Diagnostics.Append(response.State.Set(ctx, &config)...)
			return
		}
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	config.RouteTables = []CtyunVpcRouteTablesModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, r := range resp.ReturnObj.RouteTables {
		item := CtyunVpcRouteTablesModel{
			Name:            utils.SecStringValue(r.Name),
			VpcID:           utils.SecStringValue(r.VpcID),
			RouteTableID:    utils.SecStringValue(r.Id),
			Freezing:        utils.SecBoolValue(r.Freezing),
			RouteRulesCount: types.Int32Value(r.RouteRulesCount),
			CreatedAt:       utils.SecStringValue(r.CreatedAt),
			UpdatedAt:       utils.SecStringValue(r.UpdatedAt),
			Type:            types.Int32Value(r.RawType),
			Origin:          utils.SecStringValue(r.Origin),
		}
		config.RouteTables = append(config.RouteTables, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpcRouteTables) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
