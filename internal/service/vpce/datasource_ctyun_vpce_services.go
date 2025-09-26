package vpce

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunVpceServices{}
	_ datasource.DataSourceWithConfigure = &ctyunVpceServices{}
)

type ctyunVpceServices struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpceServices() datasource.DataSource {
	return &ctyunVpceServices{}
}

func (c *ctyunVpceServices) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpce_services"
}

type CtyunVpceServicesRule struct {
	Protocol     types.String `tfsdk:"protocol"`
	ServerPort   types.Int32  `tfsdk:"server_port"`
	EndpointPort types.Int32  `tfsdk:"endpoint_port"`
}
type CtyunVpceServicesModel struct {
	ID             types.String            `tfsdk:"id"`
	Name           types.String            `tfsdk:"name"`
	VpcID          types.String            `tfsdk:"vpc_id"`
	Description    types.String            `tfsdk:"description"`
	Type           types.String            `tfsdk:"type"`
	AutoConnection types.Bool              `tfsdk:"auto_connection"`
	Rules          []CtyunVpceServicesRule `tfsdk:"rules"`
	InstanceType   types.String            `tfsdk:"instance_type"`
	InstanceID     types.String            `tfsdk:"instance_id"`
	CreatedAt      types.String            `tfsdk:"created_at"`
	UpdatedAt      types.String            `tfsdk:"updated_at"`
}

type CtyunVpceServicesConfig struct {
	RegionID          types.String             `tfsdk:"region_id"`
	PageNo            types.Int32              `tfsdk:"page_no"`
	PageSize          types.Int32              `tfsdk:"page_size"`
	EndpointServiceID types.String             `tfsdk:"endpoint_service_id"`
	CurrentCount      types.Int32              `tfsdk:"current_count"`
	TotalCount        types.Int32              `tfsdk:"total_count"`
	TotalPage         types.Int32              `tfsdk:"total_page"`
	VpceServices      []CtyunVpceServicesModel `tfsdk:"vpce_services"`
}

func (c *ctyunVpceServices) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10042658/10217013`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"endpoint_service_id": schema.StringAttribute{
				Optional:    true,
				Description: "终端节点服务id",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1,推荐使用该字段",
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
			"vpce_services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "终端节点服务ID",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "终端节点服务名称",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "所属的专有网络id",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "接口还是反向，interface:接口，reverse:反向",
						},
						"auto_connection": schema.BoolAttribute{
							Computed:    true,
							Description: "是否自动连接",
						},
						"rules": schema.ListNestedAttribute{
							Computed:    true,
							Description: "接口规则数据",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"protocol": schema.StringAttribute{
										Computed:    true,
										Description: "协议，TCP:TCP协议,UDP:UDP协议",
										Validators: []validator.String{
											stringvalidator.OneOf("TCP", "UDP"),
										},
									},
									"server_port": schema.Int32Attribute{
										Computed:    true,
										Description: "服务端口",
									},
									"endpoint_port": schema.Int32Attribute{
										Computed:    true,
										Description: "节点端口",
									}},
							},
						},
						"instance_type": schema.StringAttribute{
							Computed:    true,
							Description: "服务后端实例类型，vm:虚机类型,bm:物理机,vip:vip类型,lb:负载均衡类型,当type为interface时，必填",
						},
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "服务后端实例id,当type为interface时，必填",
						},
						"created_at": schema.StringAttribute{
							Required:    true,
							Description: "创建时间",
						},
						"updated_at": schema.StringAttribute{
							Required:    true,
							Description: "更新时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunVpceServices) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpceServicesConfig
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
	params := &ctvpc.CtvpcNewEndpointServicesListRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	id := config.EndpointServiceID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if id != "" {
		params.Id = &id
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewEndpointServicesListApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.VpceServices = []CtyunVpceServicesModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, e := range resp.ReturnObj.EndpointServices {
		item := CtyunVpceServicesModel{
			ID:             utils.SecStringValue(e.ID),
			AutoConnection: utils.SecBoolValue(e.AutoConnection),
			Type:           utils.SecStringValue(e.RawType),
			VpcID:          utils.SecStringValue(e.VpcID),
			Name:           utils.SecStringValue(e.Name),
			CreatedAt:      utils.SecStringValue(e.CreatedAt),
			UpdatedAt:      utils.SecStringValue(e.UpdatedAt),
		}
		if len(e.Backends) > 0 {
			item.InstanceType = utils.SecStringValue(e.Backends[0].InstanceType)
			item.InstanceID = utils.SecStringValue(e.Backends[0].InstanceID)
		}

		for _, r := range e.Rules {
			rule := CtyunVpceServicesRule{
				Protocol:     utils.SecStringValue(r.Protocol),
				EndpointPort: types.Int32Value(r.EndpointPort),
				ServerPort:   types.Int32Value(r.ServerPort),
			}
			item.Rules = append(item.Rules, rule)
		}

		config.VpceServices = append(config.VpceServices, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpceServices) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
