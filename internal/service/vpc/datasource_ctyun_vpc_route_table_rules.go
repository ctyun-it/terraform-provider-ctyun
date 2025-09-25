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
	_ datasource.DataSource              = &ctyunVpcRouteTableRules{}
	_ datasource.DataSourceWithConfigure = &ctyunVpcRouteTableRules{}
)

type ctyunVpcRouteTableRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcRouteTableRules() datasource.DataSource {
	return &ctyunVpcRouteTableRules{}
}

func (c *ctyunVpcRouteTableRules) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_route_table_rules"
}

type CtyunVpcRouteTableRulesModel struct {
	NextHopID   types.String `tfsdk:"next_hop_id"`
	NextHopType types.String `tfsdk:"next_hop_type"`
	Destination types.String `tfsdk:"destination"`
	IpVersion   types.Int32  `tfsdk:"ip_version"`
	Description types.String `tfsdk:"description"`
	RuleID      types.String `tfsdk:"rule_id"`
}

type CtyunVpcRouteTableRulesConfig struct {
	RegionID     types.String                   `tfsdk:"region_id"`
	RouteTableID types.String                   `tfsdk:"route_table_id"`
	PageNo       types.Int32                    `tfsdk:"page_no"`
	PageSize     types.Int32                    `tfsdk:"page_size"`
	CurrentCount types.Int32                    `tfsdk:"current_count"`
	TotalCount   types.Int32                    `tfsdk:"total_count"`
	TotalPage    types.Int32                    `tfsdk:"total_page"`
	Rules        []CtyunVpcRouteTableRulesModel `tfsdk:"rules"`
}

func (c *ctyunVpcRouteTableRules) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026755/10171000**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"route_table_id": schema.StringAttribute{
				Required:    true,
				Description: "路由表ID",
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
			"rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "规则id",
						},
						"next_hop_id": schema.StringAttribute{
							Computed:    true,
							Description: "下一跳设备id",
						},
						"next_hop_type": schema.StringAttribute{
							Computed:    true,
							Description: "下一跳设备类型",
						},
						"destination": schema.StringAttribute{
							Computed:    true,
							Description: "无类别域间路由，例如：192.168.0.1/32",
						},
						"ip_version": schema.Int32Attribute{
							Computed:    true,
							Description: "4标识ipv4,6标识ipv6",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "规则描述",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpcRouteTableRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpcRouteTableRulesConfig
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
	params := &ctvpc.CtvpcNewRouteRulesListRequest{
		RegionID:     regionId,
		RouteTableID: config.RouteTableID.ValueString(),
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewRouteRulesListApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	// 解析返回值
	config.Rules = []CtyunVpcRouteTableRulesModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, v := range resp.ReturnObj.RouteRules {
		item := CtyunVpcRouteTableRulesModel{
			Description: utils.SecStringValue(v.Description),
			Destination: utils.SecStringValue(v.Destination),
			NextHopType: utils.SecStringValue(v.NextHopType),
			NextHopID:   utils.SecStringValue(v.NextHopID),
			RuleID:      utils.SecStringValue(v.RouteRuleID),
			IpVersion:   types.Int32Value(v.IpVersion),
		}
		config.Rules = append(config.Rules, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpcRouteTableRules) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
