package vpce

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
	_ datasource.DataSource              = &ctyunVpceServiceReverseRules{}
	_ datasource.DataSourceWithConfigure = &ctyunVpceServiceReverseRules{}
)

type ctyunVpceServiceReverseRules struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServiceReverseRules() datasource.DataSource {
	return &ctyunVpceServiceReverseRules{}
}

func (c *ctyunVpceServiceReverseRules) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service_reverse_rules"
}

type CtyunVpceServiceReverseRulesModel struct {
	ID          types.String `tfsdk:"id"`
	EndpointID  types.String `tfsdk:"endpoint_id"`
	TransitIP   types.String `tfsdk:"transit_ip"`
	TransitPort types.Int32  `tfsdk:"transit_port"`
	TargetIP    types.String `tfsdk:"target_ip"`
	TargetPort  types.Int32  `tfsdk:"target_port"`
	Protocol    types.String `tfsdk:"protocol"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

type CtyunVpceServiceReverseRulesConfig struct {
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	RegionID          types.String `tfsdk:"region_id"`
	PageNo            types.Int32  `tfsdk:"page_no"`
	PageSize          types.Int32  `tfsdk:"page_size"`

	CurrentCount types.Int32                         `tfsdk:"current_count"`
	TotalCount   types.Int32                         `tfsdk:"total_count"`
	TotalPage    types.Int32                         `tfsdk:"total_page"`
	Rules        []CtyunVpceServiceReverseRulesModel `tfsdk:"rules"`
}

func (c *ctyunVpceServiceReverseRules) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10042658/10048506**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"endpoint_service_id": schema.StringAttribute{
				Required:    true,
				Description: "终端节点服务id",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小，取值1-50",
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
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "规则ID",
						},
						"endpoint_id": schema.StringAttribute{
							Computed:    true,
							Description: "终端节点id",
						},
						"transit_ip": schema.StringAttribute{
							Computed:    true,
							Description: "中转地址",
						},
						"transit_port": schema.Int32Attribute{
							Computed:    true,
							Description: "中转端口(1-65535)",
						},
						"target_ip": schema.StringAttribute{
							Computed:    true,
							Description: "目标地址",
						},
						"target_port": schema.Int32Attribute{
							Computed:    true,
							Description: "目标端口(1-65535)",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "协议，TCP:TCP协议,UDP:UDP协议",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpceServiceReverseRules) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpceServiceReverseRulesConfig
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
	params := &ctvpc.CtvpcListEndpointServiceReverseRuleRequest{
		RegionID:          regionId,
		EndpointServiceID: config.EndpointServiceID.ValueString(),
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

	// 解析返回值
	config.Rules = []CtyunVpceServiceReverseRulesModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, rule := range resp.ReturnObj.ReverseRules {
		item := CtyunVpceServiceReverseRulesModel{
			ID:          utils.SecStringValue(rule.ID),
			EndpointID:  utils.SecStringValue(rule.EndpointID),
			TransitIP:   utils.SecStringValue(rule.TransitIPAddress),
			TargetIP:    utils.SecStringValue(rule.TargetIPAddress),
			TargetPort:  types.Int32Value(rule.TargetPort),
			TransitPort: types.Int32Value(rule.TransitPort),
			Protocol:    utils.SecStringValue(rule.Protocol),
			CreatedAt:   utils.SecStringValue(rule.CreatedAt),
		}
		config.Rules = append(config.Rules, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpceServiceReverseRules) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
