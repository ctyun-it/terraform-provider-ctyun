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
	_ datasource.DataSource              = &ctyunBandwidths{}
	_ datasource.DataSourceWithConfigure = &ctyunBandwidths{}
)

type ctyunBandwidths struct {
	meta *common.CtyunMetadata
}

func NewCtyunBandwidths() datasource.DataSource {
	return &ctyunBandwidths{}
}

func (c *ctyunBandwidths) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_bandwidths"
}

type CtyunBandwidthsEip struct {
	EipID      types.String `tfsdk:"eip_id"`
	EipAddress types.String `tfsdk:"eip_address"`
}
type CtyunBandwidthsModel struct {
	ID        types.String         `tfsdk:"id"`
	Name      types.String         `tfsdk:"name"`
	Bandwidth types.Int32          `tfsdk:"bandwidth"`
	Status    types.String         `tfsdk:"status"`
	Eips      []CtyunBandwidthsEip `tfsdk:"eips"`
	CreatedAt types.String         `tfsdk:"created_at"`
	ExpiredAt types.String         `tfsdk:"expired_at"`
}

type CtyunBandwidthsConfig struct {
	RegionID     types.String           `tfsdk:"region_id"`
	ProjectID    types.String           `tfsdk:"project_id"`
	PageNo       types.Int32            `tfsdk:"page_no"`
	PageSize     types.Int32            `tfsdk:"page_size"`
	BandwidthID  types.String           `tfsdk:"bandwidth_id"`
	CurrentCount types.Int32            `tfsdk:"current_count"`
	TotalCount   types.Int32            `tfsdk:"total_count"`
	TotalPage    types.Int32            `tfsdk:"total_page"`
	Bandwidths   []CtyunBandwidthsModel `tfsdk:"bandwidths"`
}

func (c *ctyunBandwidths) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026761`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"bandwidth_id": schema.StringAttribute{
				Optional:    true,
				Description: "带宽ID",
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
							Computed:    true,
							Description: "带宽ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "eip名称",
						},
						"bandwidth": schema.Int64Attribute{
							Computed:    true,
							Description: "带宽峰值大小，单位Mb",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "弹性ip状态，取值范围：active：有效，freezing：冻结中，expired：已过期",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间",
						},
						"expired_at": schema.StringAttribute{
							Computed:    true,
							Description: "到期时间",
						},
						"eips": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"eip_id": schema.StringAttribute{
										Computed:    true,
										Description: "eipID",
									},
									"eip_address": schema.StringAttribute{
										Computed:    true,
										Description: "eip地址",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (c *ctyunBandwidths) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunBandwidthsConfig
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
	if config.BandwidthID.IsNull() {
		err = c.list(ctx, &config)
	} else {
		err = c.show(ctx, &config)
	}
	if err != nil {
		return
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunBandwidths) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// show 查询单个
func (c *ctyunBandwidths) show(ctx context.Context, config *CtyunBandwidthsConfig) (err error) {
	// 组装请求体
	params := &ctvpc.CtvpcShowBandwidthRequest{
		RegionID:    config.RegionID.ValueString(),
		BandwidthID: config.BandwidthID.ValueString(),
	}
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	if projectId != "" {
		config.ProjectID = types.StringValue(projectId)
		params.ProjectID = &projectId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcShowBandwidthApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		if utils.SecString(resp.ErrorCode) == common.OpenapiSharedbandwidthNotFound {
			config.Bandwidths = []CtyunBandwidthsModel{}
			config.TotalPage = types.Int32Value(0)
			config.TotalCount = types.Int32Value(0)
			config.CurrentCount = types.Int32Value(0)
		} else {
			err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		}
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	// 解析返回值
	config.Bandwidths = []CtyunBandwidthsModel{}
	config.TotalPage = types.Int32Value(1)
	config.TotalCount = types.Int32Value(1)
	config.CurrentCount = types.Int32Value(1)
	b := resp.ReturnObj
	item := CtyunBandwidthsModel{
		ID:        utils.SecStringValue(b.Id),
		Name:      utils.SecStringValue(b.Name),
		Bandwidth: types.Int32Value(b.Bandwidth),
		Status:    utils.SecLowerStringValue(b.Status),
		CreatedAt: utils.SecStringValue(b.CreatedAt),
		ExpiredAt: utils.SecStringValue(b.ExpireAt),
	}
	item.Eips = []CtyunBandwidthsEip{}
	for _, e := range b.Eips {
		t := CtyunBandwidthsEip{
			EipID:      utils.SecStringValue(e.EipID),
			EipAddress: utils.SecStringValue(e.Ip),
		}
		item.Eips = append(item.Eips, t)
	}
	config.Bandwidths = append(config.Bandwidths, item)
	return
}

// list 列表
func (c *ctyunBandwidths) list(ctx context.Context, config *CtyunBandwidthsConfig) (err error) {
	// 组装请求体
	params := &ctvpc.CtvpcNewBandwidthListRequest{
		RegionID: config.RegionID.ValueString(),
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)

	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if projectId != "" {
		config.ProjectID = types.StringValue(projectId)
		params.ProjectID = &projectId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewBandwidthListApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Bandwidths = []CtyunBandwidthsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, b := range resp.ReturnObj.Bandwidths {
		item := CtyunBandwidthsModel{
			ID:        utils.SecStringValue(b.Id),
			Name:      utils.SecStringValue(b.Name),
			Bandwidth: types.Int32Value(b.Bandwidth),
			Status:    utils.SecLowerStringValue(b.Status),
			CreatedAt: utils.SecStringValue(b.CreatedAt),
			ExpiredAt: utils.SecStringValue(b.ExpireAt),
		}
		item.Eips = []CtyunBandwidthsEip{}
		for _, e := range b.Eips {
			t := CtyunBandwidthsEip{
				EipID:      utils.SecStringValue(e.EipID),
				EipAddress: utils.SecStringValue(e.Ip),
			}
			item.Eips = append(item.Eips, t)
		}
		config.Bandwidths = append(config.Bandwidths, item)
	}
	return
}
