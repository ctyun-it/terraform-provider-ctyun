package dns

import (
	"context"
	"errors"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunPrivateZones{}
	_ datasource.DataSourceWithConfigure = &CtyunPrivateZones{}
)

type CtyunPrivateZones struct {
	meta *common.CtyunMetadata
}

func NewCtyunPrivateZones() datasource.DataSource {
	return &CtyunPrivateZones{}
}
func (c *CtyunPrivateZones) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunPrivateZones) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_private_zones"
}

func (c *CtyunPrivateZones) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026757/10033667",
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，默认使用provider配置",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "内网DNS ID，精确查询",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "内网DNS名称，精确匹配",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "分页页码，默认为1",
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
				},
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页记录数，默认为10，最大50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"zones": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "内网DNS ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "内网DNS名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "内网DNS描述",
						},
						"proxy_pattern": schema.StringAttribute{
							Computed:    true,
							Description: "代理模式，zone：当前可用区不进行递归解析。 record：不完全劫持，进行递归解析代理",
						},
						"ttl": schema.Int64Attribute{
							Computed:    true,
							Description: "TTL值（秒）",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
						"vpc_associations": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vpc_id": schema.StringAttribute{
										Computed:    true,
										Description: "VPC ID",
									},
									"vpc_name": schema.StringAttribute{
										Computed:    true,
										Description: "VPC名称",
									},
								},
							},
							Description: "关联的VPC列表",
						},
					},
				},
				Description: "私有域列表",
			},
		},
	}
}

func (c *CtyunPrivateZones) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunPrivateZonesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = errors.New("region ID不能为空！")
		return
	}
	config.RegionID = types.StringValue(regionId)
	params := &ctvpc.CtvpcListPrivateZoneRequest{
		RegionID: config.RegionID.ValueString(),
		PageNo:   1,
		PageSize: 10,
	}
	if !config.ID.IsNull() {
		params.ZoneID = config.ID.ValueStringPointer()
	}
	if !config.Name.IsNull() {
		params.ZoneName = config.Name.ValueStringPointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListPrivateZoneApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("查询内网DNS列表失败！接口返回nil,请联系研发确认问题原因！")
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	var zones []CtyunPrivateZoneInfoModel
	for _, zoneItem := range resp.ReturnObj.Zones {
		var zone CtyunPrivateZoneInfoModel
		zone.Name = types.StringValue(*zoneItem.Name)
		zone.ID = types.StringValue(*zoneItem.ZoneID)
		zone.Description = types.StringValue(*zoneItem.Description)
		zone.ProxyPattern = types.StringValue(*zoneItem.ProxyPattern)
		zone.TTL = types.Int32Value(zoneItem.TTL)
		zone.CreateTime = types.StringValue(*zoneItem.CreatedAt)
		zone.UpdateTime = types.StringValue(*zoneItem.UpdatedAt)
		// 处理vpc ids
		var vpcAssociations []CtyunPrivateZoneVpcModel
		for _, vpcItem := range zoneItem.VpcAssociations {
			vpcAssociations = append(vpcAssociations, CtyunPrivateZoneVpcModel{
				VpcID:   types.StringValue(*vpcItem.VpcID),
				VpcName: types.StringValue(*vpcItem.VpcName),
			})
		}
		zone.VpcAssociations = vpcAssociations
		zones = append(zones, zone)
	}
	config.Zones = zones
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunPrivateZoneVpcModel struct {
	VpcID   types.String `tfsdk:"vpc_id"`
	VpcName types.String `tfsdk:"vpc_name"`
}

type CtyunPrivateZoneInfoModel struct {
	ID              types.String               `tfsdk:"id"`
	Name            types.String               `tfsdk:"name"`
	Description     types.String               `tfsdk:"description"`
	ProxyPattern    types.String               `tfsdk:"proxy_pattern"`
	TTL             types.Int32                `tfsdk:"ttl"`
	VpcAssociations []CtyunPrivateZoneVpcModel `tfsdk:"vpc_associations"`
	CreateTime      types.String               `tfsdk:"create_time"`
	UpdateTime      types.String               `tfsdk:"update_time"`
}
type CtyunPrivateZonesConfig struct {
	RegionID types.String                `tfsdk:"region_id"`
	ID       types.String                `tfsdk:"id"`
	Name     types.String                `tfsdk:"name"`
	PageNo   types.Int32                 `tfsdk:"page_no"`
	PageSize types.Int32                 `tfsdk:"page_size"`
	Zones    []CtyunPrivateZoneInfoModel `tfsdk:"zones"`
}
