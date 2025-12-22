package vpc

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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunNetResourcesByTag{}
	_ datasource.DataSourceWithConfigure = &ctyunNetResourcesByTag{}
)

type ctyunNetResourcesByTag struct {
	meta *common.CtyunMetadata
}

func NewCtyunNetResourcesByTag() datasource.DataSource {
	return &ctyunNetResourcesByTag{}
}

func (c *ctyunNetResourcesByTag) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_net_resources_by_tag"
}

type CtyunNetResourceModel struct {
	ResourceID   types.String `tfsdk:"resource_id"`
	ResourceType types.String `tfsdk:"resource_type"`
}

type CtyunNetResourcesByTagConfig struct {
	RegionID     types.String            `tfsdk:"region_id"`
	LabelID      types.String            `tfsdk:"label_id"`
	LabelKey     types.String            `tfsdk:"label_key"`
	LabelValue   types.String            `tfsdk:"label_value"`
	PageNumber   types.Int32             `tfsdk:"page_number"`
	PageSize     types.Int32             `tfsdk:"page_size"`
	TotalCount   types.Int32             `tfsdk:"total_count"`
	CurrentCount types.Int32             `tfsdk:"current_count"`
	TotalPage    types.Int32             `tfsdk:"total_page"`
	Resources    []CtyunNetResourceModel `tfsdk:"resources"`
}

func (c *ctyunNetResourcesByTag) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"label_id": schema.StringAttribute{
				Optional:    true,
				Description: "标签ID，label的三个参数至少选填一个",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("label_key"),
						path.MatchRoot("label_value"),
					),
				},
			},
			"label_key": schema.StringAttribute{
				Optional:    true,
				Description: "标签键",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("label_id"),
						path.MatchRoot("label_value"),
					),
				},
			},
			"label_value": schema.StringAttribute{
				Optional:    true,
				Description: "标签值",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("label_id"),
						path.MatchRoot("label_key"),
					),
				},
			},
			"page_number": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码，默认值为1",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "分页查询时每页的行数，最大值为50，默认值为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},
			"total_count": schema.Int32Attribute{
				Computed:    true,
				Description: "列表条目数",
			},
			"current_count": schema.Int32Attribute{
				Computed:    true,
				Description: "分页查询时每页的行数",
			},
			"total_page": schema.Int32Attribute{
				Computed:    true,
				Description: "总页数",
			},
			"resources": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源ID",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "资源类型，支持VPC、VPCE、VPCES、安全组、弹性负载均衡、弹性IP、共享带宽、子网等",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunNetResourcesByTag) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunNetResourcesByTagConfig
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

	params := &ctvpc.CtvpcQueryResourcesByLabelRequest{
		RegionID: regionId,
	}

	// 设置可选参数
	if !config.LabelID.IsNull() {
		labelID := config.LabelID.ValueString()
		params.LabelID = &labelID
	}
	if !config.LabelKey.IsNull() {
		labelKey := config.LabelKey.ValueString()
		params.LabelKey = &labelKey
	}
	if !config.LabelValue.IsNull() {
		labelValue := config.LabelValue.ValueString()
		params.LabelValue = &labelValue
	}
	if !config.PageNumber.IsNull() {
		params.PageNumber = config.PageNumber.ValueInt32()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32()
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcQueryResourcesByLabelApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Resources = []CtyunNetResourceModel{}

	// 处理返回的对象
	returnObj := resp.ReturnObj
	config.TotalCount = types.Int32Value(returnObj.TotalCount)
	config.CurrentCount = types.Int32Value(returnObj.CurrentCount)
	config.TotalPage = types.Int32Value(returnObj.TotalPage)

	// 遍历所有资源
	for _, r := range returnObj.Results {
		item := CtyunNetResourceModel{
			ResourceID:   utils.SecStringValue(r.ResourceID),
			ResourceType: utils.SecStringValue(r.ResourceType),
		}
		config.Resources = append(config.Resources, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunNetResourcesByTag) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
