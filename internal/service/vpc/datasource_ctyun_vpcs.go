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
	_ datasource.DataSource              = &ctyunVpcs{}
	_ datasource.DataSourceWithConfigure = &ctyunVpcs{}
)

type ctyunVpcs struct {
	meta *common.CtyunMetadata
}

func NewCtyunVpcs() datasource.DataSource {
	return &ctyunVpcs{}
}

func (c *ctyunVpcs) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpcs"
}

type CtyunVpcsModel struct {
	VpcID          types.String `tfsdk:"vpc_id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	CIDR           types.String `tfsdk:"cidr"`
	Ipv6Enabled    types.Bool   `tfsdk:"ipv6_enabled"`
	EnableIpv6     types.Bool   `tfsdk:"enable_ipv6"`
	Ipv6CIDRS      []string     `tfsdk:"ipv6_cidrs"`
	SubnetIDs      []string     `tfsdk:"subnet_ids"`
	NatGatewayIDs  []string     `tfsdk:"nat_gateway_ids"`
	SecondaryCIDRS []string     `tfsdk:"secondary_cidrs"`
	ProjectID      types.String `tfsdk:"project_id"`
}

type CtyunVpcsConfig struct {
	RegionID  types.String `tfsdk:"region_id"`
	VpcID     types.String `tfsdk:"vpc_id"`
	PageNo    types.Int32  `tfsdk:"page_no"`
	PageSize  types.Int32  `tfsdk:"page_size"`
	ProjectID types.String `tfsdk:"project_id"`

	CurrentCount types.Int32      `tfsdk:"current_count"`
	TotalCount   types.Int32      `tfsdk:"total_count"`
	TotalPage    types.Int32      `tfsdk:"total_page"`
	Vpcs         []CtyunVpcsModel `tfsdk:"vpcs"`
}

func (c *ctyunVpcs) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10026755/10028487`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池id，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "多个VPC的ID之间用半角逗号（,）隔开。",
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
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "企业项目ID，默认为`0`",
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
							Description: "vpc示例ID",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"cidr": schema.StringAttribute{
							Computed:    true,
							Description: "子网",
						},
						"ipv6_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启ipv6",
						},
						"enable_ipv6": schema.BoolAttribute{
							Computed:    true,
							Description: "是否开启ipv6",
						},
						"ipv6_cidrs": schema.ListAttribute{
							Computed:    true,
							Description: "ipv6子网列表",
							ElementType: types.StringType,
						},
						"subnet_ids": schema.ListAttribute{
							Computed:    true,
							Description: "子网id列表",
							ElementType: types.StringType,
						},
						"nat_gateway_ids": schema.ListAttribute{
							Computed:    true,
							Description: "网关id列表",
							ElementType: types.StringType,
						},
						"secondary_cidrs": schema.ListAttribute{
							Computed:    true,
							Description: "附加网段",
							ElementType: types.StringType,
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "企业项目ID，默认为`0`",
						},
					},
				},
			},
		}}
}

func (c *ctyunVpcs) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunVpcsConfig
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
	params := &ctvpc.CtvpcNewVpcListRequest{
		RegionID: regionId,
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	projectId := c.meta.GetExtraIfEmpty(config.ProjectID.ValueString(), common.ExtraProjectId)
	vpcId := config.VpcID.ValueString()
	if pageNo > 0 {
		params.PageNo = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	if projectId != "" {
		params.ProjectID = &projectId
	}
	if vpcId != "" {
		params.VpcID = &vpcId
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCtVpcApis.CtvpcNewVpcListApi.Do(ctx, c.meta.SdkCredential, params)
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
	config.Vpcs = []CtyunVpcsModel{}
	config.TotalPage = types.Int32Value(resp.ReturnObj.TotalPage)
	config.TotalCount = types.Int32Value(resp.ReturnObj.TotalCount)
	config.CurrentCount = types.Int32Value(resp.ReturnObj.CurrentCount)
	for _, v := range resp.ReturnObj.Vpcs {
		item := CtyunVpcsModel{
			VpcID:          utils.SecStringValue(v.VpcID),
			Name:           utils.SecStringValue(v.Name),
			Description:    utils.SecStringValue(v.Description),
			CIDR:           utils.SecStringValue(v.CIDR),
			Ipv6Enabled:    utils.SecBoolValue(v.Ipv6Enabled),
			EnableIpv6:     utils.SecBoolValue(v.EnableIpv6),
			Ipv6CIDRS:      utils.StrPointerArrayToStrArray(v.Ipv6CIDRS),
			SubnetIDs:      utils.StrPointerArrayToStrArray(v.SubnetIDs),
			NatGatewayIDs:  utils.StrPointerArrayToStrArray(v.NatGatewayIDs),
			SecondaryCIDRS: utils.StrPointerArrayToStrArray(v.SecondaryCIDRS),
			ProjectID:      utils.SecStringValue(v.ProjectID),
		}
		config.Vpcs = append(config.Vpcs, item)
	}
	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunVpcs) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
