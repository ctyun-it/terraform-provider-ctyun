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
	_ datasource.DataSource              = &ctyunDhcpOptionSetAssociationVpcs{}
	_ datasource.DataSourceWithConfigure = &ctyunDhcpOptionSetAssociationVpcs{}
)

type ctyunDhcpOptionSetAssociationVpcs struct {
	meta *common.CtyunMetadata
}

func NewCtyunDhcpOptionSetAssociationVpcs() datasource.DataSource {
	return &ctyunDhcpOptionSetAssociationVpcs{}
}

func (c *ctyunDhcpOptionSetAssociationVpcs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dhcpoptionset_association_vpcs"
}

type CtyunDhcpOptionSetAssociationVpcsModel struct {
	VpcID          types.String `tfsdk:"vpc_id"`
	Name           types.String `tfsdk:"name"`
	Cidr           types.String `tfsdk:"cidr"`
	SecondaryCidrs []string     `tfsdk:"secondary_cidrs"`
	Status         types.String `tfsdk:"status"`
}

type CtyunDhcpOptionSetAssociationVpcsConfig struct {
	RegionID         types.String `tfsdk:"region_id"`
	DhcpOptionSetsID types.String `tfsdk:"dhcp_option_sets_id"`
	PageNo           types.Int32  `tfsdk:"page_no"`
	PageSize         types.Int32  `tfsdk:"page_size"`

	CurrentCount types.Int32                              `tfsdk:"current_count"`
	TotalCount   types.Int32                              `tfsdk:"total_count"`
	TotalPage    types.Int32                              `tfsdk:"total_page"`
	Vpcs         []CtyunDhcpOptionSetAssociationVpcsModel `tfsdk:"vpcs"`
}

func (c *ctyunDhcpOptionSetAssociationVpcs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028310`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"dhcp_option_sets_id": schema.StringAttribute{
				Required:    true,
				Description: "DHCP选项集ID",
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
			"vpcs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "VPC名称",
						},
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "VPC CIDR",
						},
						"secondary_cidrs": schema.ListAttribute{
							Computed:    true,
							Description: "扩展网段列表",
							ElementType: types.StringType,
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态",
						},
					},
				},
			},
		}}
}

func (c *ctyunDhcpOptionSetAssociationVpcs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunDhcpOptionSetAssociationVpcsConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("regionId不能为空")
		return
	}
	// 组装请求体
	params := &ctvpc.CtvpcDhcplistvpcRequest{
		RegionID:         regionId,
		DhcpOptionSetsID: config.DhcpOptionSetsID.ValueString(),
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
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcDhcplistvpcApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Vpcs = []CtyunDhcpOptionSetAssociationVpcsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)

	for _, v := range resp.ReturnObj.Results {
		item := CtyunDhcpOptionSetAssociationVpcsModel{
			VpcID:          utils.SecStringValue(v.VpcID),
			Name:           utils.SecStringValue(v.Name),
			Cidr:           utils.SecStringValue(v.Cidr),
			SecondaryCidrs: utils.StrPointerArrayToStrArray(v.SecondaryCidrs),
			Status:         utils.SecStringValue(v.Status),
		}
		config.Vpcs = append(config.Vpcs, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunDhcpOptionSetAssociationVpcs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
