package vpc

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunVips{}
	_ datasource.DataSourceWithConfigure = &CtyunVips{}
)

type CtyunVips struct {
	meta *common.CtyunMetadata
}

func NewCtyunVips() datasource.DataSource {
	return &CtyunVips{}
}

func (c *CtyunVips) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vips"
}

type CtyunVipsModel struct {
	ID        types.String `tfsdk:"id"`
	Ipv4      types.String `tfsdk:"ipv4"`
	VpcID     types.String `tfsdk:"vpc_id"`
	SubnetID  types.String `tfsdk:"subnet_id"`
	ProjectID types.String `tfsdk:"project_id"`
}

type CtyunVipsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	ProjectID types.String `tfsdk:"project_id"`
	PageNo    types.Int32  `tfsdk:"page_no"`
	PageSize  types.Int32  `tfsdk:"page_size"`

	CurrentCount types.Int32      `tfsdk:"current_count"`
	TotalCount   types.Int32      `tfsdk:"total_count"`
	TotalPage    types.Int32      `tfsdk:"total_page"`
	Vips         []CtyunVipsModel `tfsdk:"vips"`
}

func (c *CtyunVips) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，默认为`0`",
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
			"vips": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "高可用虚IP的ID",
						},
						"ipv4": schema.StringAttribute{
							Computed:    true,
							Description: "IPv4地址",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云的ID",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID",
						},
					},
				},
			},
		}}
}

func (c *CtyunVips) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVipsConfig
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
	params := &ctvpc.CtvpcListHavipRequest{
		ClientToken: uuid.NewString(),
		RegionID:    regionId,
	}

	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)

	if pageNo > 0 {
		// 注意：ListHavip API 不支持分页参数，需要在客户端处理
	}
	if pageSize > 0 {
		// 注意：ListHavip API 不支持分页参数，需要在客户端处理
	}
	if projectId != "" {
		params.ProjectID = &projectId
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcListHavipApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != 800 {
		err = fmt.Errorf("API return error. Message: %s Description: %s", *resp.Message, *resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return empty response object")
		return
	}

	// 解析返回值
	config.Vips = []CtyunVipsModel{}

	// 由于API不支持分页，我们手动设置分页信息
	config.TotalCount = types.Int32Value(int32(len(resp.ReturnObj)))
	config.CurrentCount = types.Int32Value(int32(len(resp.ReturnObj)))
	config.TotalPage = types.Int32Value(1)

	for _, v := range resp.ReturnObj {
		item := CtyunVipsModel{
			ID:       utils.SecStringValue(v.Id),
			Ipv4:     utils.SecStringValue(v.Ipv4),
			VpcID:    utils.SecStringValue(v.VpcID),
			SubnetID: utils.SecStringValue(v.SubnetID),
		}
		if projectId != "" {
			item.ProjectID = types.StringValue(projectId)
		} else {
			item.ProjectID = types.StringValue("0")
		}
		config.Vips = append(config.Vips, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *CtyunVips) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
