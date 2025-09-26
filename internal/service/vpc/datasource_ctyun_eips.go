package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &ctyunEips{}
	_ datasource.DataSourceWithConfigure = &ctyunEips{}
)

type ctyunEips struct {
	meta *common.CtyunMetadata
}

func NewCtyunEips() datasource.DataSource {
	return &ctyunEips{}
}

func (c *ctyunEips) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_eips"
}

type CtyunEipsModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	EipAddress       types.String `tfsdk:"eip_address"`
	AssociationID    types.String `tfsdk:"association_id"`
	AssociationType  types.String `tfsdk:"association_type"`
	PrivateIpAddress types.String `tfsdk:"private_ip_address"`
	Bandwidth        types.Int32  `tfsdk:"bandwidth"`
	BandwidthID      types.String `tfsdk:"bandwidth_id"`
	BandwidthType    types.String `tfsdk:"bandwidth_type"`
	Status           types.String `tfsdk:"status"`
	Tags             types.String `tfsdk:"tags"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	ExpiredAt        types.String `tfsdk:"expired_at"`
}

type CtyunEipsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	ProjectID types.String `tfsdk:"project_id"`
	PageNo    types.Int32  `tfsdk:"page_no"`
	PageSize  types.Int32  `tfsdk:"page_size"`
	Ids       types.String `tfsdk:"ids"`
	Status    types.String `tfsdk:"status"`
	EipType   types.String `tfsdk:"eip_type"`
	Ip        types.String `tfsdk:"ip"`

	CurrentCount types.Int32      `tfsdk:"current_count"`
	TotalCount   types.Int32      `tfsdk:"total_count"`
	TotalPage    types.Int32      `tfsdk:"total_page"`
	Eips         []CtyunEipsModel `tfsdk:"eips"`
}

func (c *ctyunEips) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026753`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，默认为`0`",
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
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "使用,连接",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "弹性ip状态，支持ACTIVE（已绑定）/ DOWN（未绑定）/ FREEZING（已冻结）/ EXPIRED（已过期）",
				Validators: []validator.String{
					stringvalidator.OneOf("ACTIVE", "DOWN", "FREEZING", "EXPIRED"),
				},
			},
			"eip_type": schema.StringAttribute{
				Optional:    true,
				Description: "eip类型normal/cn2",
				Validators: []validator.String{
					stringvalidator.OneOf("normal", "cn2"),
				},
			},
			"ip": schema.StringAttribute{
				Optional:    true,
				Description: "弹性IP的ip地址",
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
			"eips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "eipID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "eip名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"eip_address": schema.StringAttribute{
							Computed:    true,
							Description: "eip地址",
						},
						"association_id": schema.StringAttribute{
							Computed:    true,
							Description: "当前绑定的实例的ID",
						},
						"association_type": schema.StringAttribute{
							Computed:    true,
							Description: "当前绑定的实例类型:loadbalancer/instance/portforwording/vip/physical_instance",
						},
						"private_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "交换机网段内的一个IP地址",
						},
						"bandwidth": schema.Int64Attribute{
							Computed:    true,
							Description: "带宽峰值大小，单位Mb",
						},
						"bandwidth_id": schema.StringAttribute{
							Computed:    true,
							Description: "绑定的共享带宽ID",
						},
						"bandwidth_type": schema.StringAttribute{
							Computed:    true,
							Description: "eip带宽规格：standalone/upflowc",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "弹性ip状态，取值范围：active：有效，down：未绑定，error：出错，updating：更新中，banding_or_unbangding：绑定解绑中，deleting：删除中，deleted：已删除，expired：已过期",
						},
						"tags": schema.StringAttribute{
							Computed:    true,
							Description: "EIP的标签集合",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间",
						},
						"expired_at": schema.StringAttribute{
							Computed:    true,
							Description: "到期时间",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunEips) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunEipsConfig
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
	params := &ctvpc.CtvpcNewEipListRequest{
		RegionID:    regionId,
		Status:      config.Status.ValueStringPointer(),
		EipType:     config.EipType.ValueStringPointer(),
		Ip:          config.Ip.ValueStringPointer(),
		ClientToken: uuid.NewString(),
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	ids := config.Ids.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if projectId != "" {
		params.ProjectID = &projectId
		config.ProjectID = types.StringValue(projectId)
	}
	if len(ids) > 0 {
		params.Ids = utils.StrArrayToStrPointerArray(strings.Split(ids, ","))
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewEipListApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Eips = []CtyunEipsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, e := range resp.ReturnObj.Eips {
		item := CtyunEipsModel{
			ID:               utils.SecStringValue(e.ID),
			Name:             utils.SecStringValue(e.Name),
			Description:      utils.SecStringValue(e.Description),
			EipAddress:       utils.SecStringValue(e.EipAddress),
			AssociationID:    utils.SecStringValue(e.AssociationID),
			AssociationType:  utils.SecStringValue(e.AssociationType),
			PrivateIpAddress: utils.SecStringValue(e.PrivateIpAddress),
			Bandwidth:        types.Int32Value(e.Bandwidth),
			BandwidthID:      utils.SecStringValue(e.BandwidthID),
			BandwidthType:    utils.SecStringValue(e.BandwidthType),
			Status:           utils.SecLowerStringValue(e.Status),
			Tags:             utils.SecStringValue(e.Tags),
			CreatedAt:        utils.SecStringValue(e.CreatedAt),
			UpdatedAt:        utils.SecStringValue(e.UpdatedAt),
			ExpiredAt:        utils.SecStringValue(e.ExpiredAt),
		}
		config.Eips = append(config.Eips, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunEips) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
