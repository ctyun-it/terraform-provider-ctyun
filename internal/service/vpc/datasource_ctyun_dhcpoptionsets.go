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
	_ datasource.DataSource              = &ctyunDhcpOptionSets{}
	_ datasource.DataSourceWithConfigure = &ctyunDhcpOptionSets{}
)

type ctyunDhcpOptionSets struct {
	meta *common.CtyunMetadata
}

func NewCtyunDhcpOptionSets() datasource.DataSource {
	return &ctyunDhcpOptionSets{}
}

func (c *ctyunDhcpOptionSets) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dhcpoptionsets"
}

type CtyunDhcpOptionSetsModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DomainName  types.String `tfsdk:"domain_name"`
	DnsList     []string     `tfsdk:"dns_list"`
	VpcList     []string     `tfsdk:"vpc_list"`
	CreatedAt   types.String `tfsdk:"create_time"`
	UpdatedAt   types.String `tfsdk:"update_time"`
}

type CtyunDhcpOptionSetsConfig struct {
	RegionID     types.String `tfsdk:"region_id"`
	QueryContent types.String `tfsdk:"query_content"`
	PageNo       types.Int32  `tfsdk:"page_no"`
	PageSize     types.Int32  `tfsdk:"page_size"`

	CurrentCount   types.Int32                `tfsdk:"current_count"`
	TotalCount     types.Int32                `tfsdk:"total_count"`
	TotalPage      types.Int32                `tfsdk:"total_page"`
	DhcpOptionSets []CtyunDhcpOptionSetsModel `tfsdk:"dhcpoptionsets"`
}

func (c *ctyunDhcpOptionSets) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊查询内容",
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
			"dhcpoptionsets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "DHCP选项集ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"domain_name": schema.StringAttribute{
							Computed:    true,
							Description: "域名列表",
						},
						"dns_list": schema.ListAttribute{
							Computed:    true,
							Description: "DNS服务器地址列表",
							ElementType: types.StringType,
						},
						"vpc_list": schema.ListAttribute{
							Computed:    true,
							Description: "关联的VPC列表",
							ElementType: types.StringType,
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
					},
				},
			},
		}}
}

func (c *ctyunDhcpOptionSets) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunDhcpOptionSetsConfig
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
	params := &ctvpc.CtvpcDhcpoptionsetsqueryRequest{
		RegionID: regionId,
	}

	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	queryContent := config.QueryContent.ValueString()

	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if queryContent != "" {
		params.QueryContent = &queryContent
	}

	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcpoptionsetsqueryApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.DhcpOptionSets = []CtyunDhcpOptionSetsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)

	for _, v := range resp.ReturnObj.Results {
		item := CtyunDhcpOptionSetsModel{
			ID:          utils.SecStringValue(v.DhcpOptionSetsID),
			Name:        utils.SecStringValue(v.Name),
			Description: utils.SecStringValue(v.Description),
			DomainName:  utils.SecStringValue(v.DomainName),
			DnsList:     utils.StrPointerArrayToStrArray(v.DnsList),
			VpcList:     utils.StrPointerArrayToStrArray(v.VpcList),
			CreatedAt:   utils.SecStringValue(v.CreatedAt),
			UpdatedAt:   utils.SecStringValue(v.UpdatedAt),
		}
		config.DhcpOptionSets = append(config.DhcpOptionSets, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunDhcpOptionSets) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
