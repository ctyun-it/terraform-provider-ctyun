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
	_ datasource.DataSource              = &ctyunVpceServiceTransitIPs{}
	_ datasource.DataSourceWithConfigure = &ctyunVpceServiceTransitIPs{}
)

type ctyunVpceServiceTransitIPs struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServiceTransitIPs() datasource.DataSource {
	return &ctyunVpceServiceTransitIPs{}
}

func (c *ctyunVpceServiceTransitIPs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_service_transit_ips"
}

type CtyunVpceServiceTransitIPsModel struct {
	ID        types.String `tfsdk:"id"`
	SubnetID  types.String `tfsdk:"subnet_id"`
	TransitIP types.String `tfsdk:"transit_ip"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type CtyunVpceServiceTransitIPsConfig struct {
	EndpointServiceID types.String `tfsdk:"endpoint_service_id"`
	RegionID          types.String `tfsdk:"region_id"`
	PageNo            types.Int32  `tfsdk:"page_no"`
	PageSize          types.Int32  `tfsdk:"page_size"`

	CurrentCount types.Int32                       `tfsdk:"current_count"`
	TotalCount   types.Int32                       `tfsdk:"total_count"`
	TotalPage    types.Int32                       `tfsdk:"total_page"`
	IPs          []CtyunVpceServiceTransitIPsModel `tfsdk:"ips"`
}

func (c *ctyunVpceServiceTransitIPs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10042658/10048507**`,
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
			"ips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID，这里使用中转地址",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"transit_ip": schema.StringAttribute{
							Computed:    true,
							Description: "中转地址",
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

func (c *ctyunVpceServiceTransitIPs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpceServiceTransitIPsConfig
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
	params := &ctvpc.CtvpcListEndpointServiceTransitIPRequest{
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
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListEndpointServiceTransitIPApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.IPs = []CtyunVpceServiceTransitIPsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, e := range resp.ReturnObj.TransitIPs {
		item := CtyunVpceServiceTransitIPsModel{
			ID:        utils.SecStringValue(e.TransitIP),
			SubnetID:  utils.SecStringValue(e.SubnetID),
			TransitIP: utils.SecStringValue(e.TransitIP),
			CreatedAt: utils.SecStringValue(e.CreatedAt),
		}
		config.IPs = append(config.IPs, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpceServiceTransitIPs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
