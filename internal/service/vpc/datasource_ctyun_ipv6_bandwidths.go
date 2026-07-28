package vpc

import (
	"context"
	"fmt"

	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunIPv6Bandwidths{}
	_ datasource.DataSourceWithConfigure = &ctyunIPv6Bandwidths{}
)

type ctyunIPv6Bandwidths struct {
	meta *common.CtyunMetadata
}

func NewCtyunIPv6Bandwidths() datasource.DataSource {
	return &ctyunIPv6Bandwidths{}
}

func (c *ctyunIPv6Bandwidths) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ipv6_bandwidths"
}

type CtyunIPv6BandwidthsModel struct {
	ID            types.String `tfsdk:"id"`
	Status        types.String `tfsdk:"status"`
	Name          types.String `tfsdk:"name"`
	Bandwidth     types.Int64  `tfsdk:"bandwidth"`
	Expired       types.Bool   `tfsdk:"expired"`
	ResourceSpec  types.String `tfsdk:"resource_spec"`
	CreatedAt     types.String `tfsdk:"create_time"`
	ExpiredAt     types.String `tfsdk:"expire_time"`
	IpAddress     types.String `tfsdk:"ip_address"`
	IPv6GatewayID types.String `tfsdk:"ipv6_gateway_id"`
	PaymentType   types.String `tfsdk:"payment_type"`
}

type CtyunIPv6BandwidthsConfig struct {
	ID           types.String               `tfsdk:"id"`
	RegionID     types.String               `tfsdk:"region_id"`
	QueryContent types.String               `tfsdk:"query_content"`
	PageNo       types.Int64                `tfsdk:"page_no"`
	PageSize     types.Int64                `tfsdk:"page_size"`
	Bandwidths   []CtyunIPv6BandwidthsModel `tfsdk:"bandwidths"`
	CurrentCnt   types.Int32                `tfsdk:"current_count"`
	TotalCnt     types.Int32                `tfsdk:"total_count"`
	TotalPage    types.Int32                `tfsdk:"total_page"`
}

func (c *ctyunIPv6Bandwidths) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询IPv6带宽列表", "弹性IP（Elastic IP，EIP）", "https://www.ctyun.cn/document/10026753/10037269"),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "IPv6带宽ID",
				Optional:    true,
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"query_content": schema.StringAttribute{
				Description: "【模糊查询】 IPv6 带宽实例名称 / 带宽 ID",
				Optional:    true,
			},
			"page_no": schema.Int64Attribute{
				Description: "列表的页码，默认值为 1",
				Optional:    true,
			},
			"page_size": schema.Int64Attribute{
				Description: "分页查询时每页的行数，最大值为 50，默认值为 10",
				Optional:    true,
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
			"bandwidths": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "IPv6带宽ID",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "IPv6带宽状态，ACTIVE（正常） / EXPIRED（过期） / FREEZING（冻结） /CREATEING（创建中）",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "eip名称",
						},
						"bandwidth": schema.Int64Attribute{
							Description: "带宽大小，默认值为1(Mbps)",
							Computed:    true,
						},
						"expired": schema.BoolAttribute{
							Computed:    true,
							Description: "是否过期",
						},
						"resource_spec": schema.StringAttribute{
							Computed:    true,
							Description: "独享/共享模式",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"expire_time": schema.StringAttribute{
							Computed:    true,
							Description: "过期时间",
						},
						"ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "IP地址",
						},
						"ipv6_gateway_id": schema.StringAttribute{
							Computed:    true,
							Description: "IPv6网关",
						},
						"payment_type": schema.StringAttribute{
							Computed:    true,
							Description: "计费类型",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunIPv6Bandwidths) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunIPv6BandwidthsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id不能为空")
		return
	}
	config.RegionID = types.StringValue(regionId)

	params := &ctvpc.CtvpcNewIPv6BandwidthListRequest{
		RegionID:     regionId,
		QueryContent: config.QueryContent.ValueStringPointer(),
		BandwidthID:  config.ID.ValueStringPointer(),
		PageNo:       int32(config.PageNo.ValueInt64()),
		PageSize:     int32(config.PageSize.ValueInt64()),
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewIPv6BandwidthListApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCnt = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCnt = types.Int32Value(resp.ReturnObj.CurrentCount)
	config.Bandwidths = []CtyunIPv6BandwidthsModel{}

	for _, e := range resp.ReturnObj.Objs {
		item := CtyunIPv6BandwidthsModel{
			ID:            utils.SecStringValue(e.ID),
			Status:        utils.SecLowerStringValue(e.Status),
			Name:          utils.SecStringValue(e.Name),
			Bandwidth:     types.Int64Value(int64(e.Bandwidth)),
			Expired:       types.BoolValue(e.Expired),
			ResourceSpec:  utils.SecStringValue(e.ResourceSpec),
			CreatedAt:     utils.SecStringValue(e.CreatedTime),
			ExpiredAt:     utils.SecStringValue(e.ExpiredTime),
			IpAddress:     utils.SecStringValue(e.IpAddress),
			IPv6GatewayID: utils.SecStringValue(e.Ipv6GatewayID),
			PaymentType:   utils.SecStringValue(e.PaymentType),
		}
		config.Bandwidths = append(config.Bandwidths, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunIPv6Bandwidths) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
