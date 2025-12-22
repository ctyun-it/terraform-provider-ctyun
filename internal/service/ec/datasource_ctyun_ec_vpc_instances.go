package ec

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ec"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &CtyunExpressConnectionVpcInstances{}
	_ datasource.DataSourceWithConfigure = &CtyunExpressConnectionVpcInstances{}
)

type CtyunExpressConnectionVpcInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunExpressConnectionVpcInstances() datasource.DataSource {
	return &CtyunExpressConnectionVpcInstances{}
}
func (c *CtyunExpressConnectionVpcInstances) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

func (c *CtyunExpressConnectionVpcInstances) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ec_vpc_instances"
}

func (c *CtyunExpressConnectionVpcInstances) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "-> 详细说明请见文档：https://www.ctyun.cn/document/10026763/10038256",
		Attributes: map[string]schema.Attribute{
			"ec_id": schema.StringAttribute{
				Required:    true,
				Description: "云间高速实例ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthBetween(36, 36),
				},
			},
			"cgw_id": schema.StringAttribute{
				Optional:    true,
				Description: "云网关ID",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Description: "VPC实例ID，精确查询",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"query_content": schema.StringAttribute{
				Optional:    true,
				Description: "模糊匹配（支持VPCID、VPCName、dcName）",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Description: "状态（creating：加载中，running：已连接，removing：卸载中，flushing：路由待更新，error：失败）",
				Validators: []validator.String{
					stringvalidator.OneOf("creating", "running", "removing", "flushing", "error"),
				},
			},
			"is_auth": schema.Int32Attribute{
				Optional:    true,
				Description: "是否跨账号实例（0：本账号，1：跨账号）",
				Validators: []validator.Int32{
					int32validator.OneOf(0, 1),
				},
			},
			"is_exclusive": schema.Int32Attribute{
				Optional:    true,
				Description: "是否专属云实例（0：公有云，1：专属云）",
				Validators: []validator.Int32{
					int32validator.OneOf(0, 1),
				},
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
				Description: "每页记录数，默认为10",
				Validators: []validator.Int32{
					int32validator.Between(1, 100),
				},
			},
			"vpc_instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ec_id": schema.StringAttribute{
							Computed:    true,
							Description: "云间高速实例ID",
						},
						"cgw_id": schema.StringAttribute{
							Computed:    true,
							Description: "云网关ID",
						},
						"cgw_name": schema.StringAttribute{
							Computed:    true,
							Description: "云网关名称",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "状态",
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "VPC名称",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC ID",
						},
						"vpc_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "VPC CIDR",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"region_name": schema.StringAttribute{
							Computed:    true,
							Description: "资源池名称",
						},
						"region_type": schema.StringAttribute{
							Computed:    true,
							Description: "资源池类型（CNP/MAZ/OS/CS/PRVT）",
						},
						"rtb_name": schema.StringAttribute{
							Computed:    true,
							Description: "路由表名称",
						},
						"rtb_id": schema.StringAttribute{
							Computed:    true,
							Description: "路由表ID",
						},
						"exclusive_id": schema.StringAttribute{
							Computed:    true,
							Description: "专属云资源池ID",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC实例ID",
						},
						"is_auth": schema.Int32Attribute{
							Computed:    true,
							Description: "是否跨账号实例（0：本账号，1：跨账号）",
						},
						"is_exclusive": schema.Int32Attribute{
							Computed:    true,
							Description: "是否专属云实例（0：公有云，1：专属云）",
						},
						"route_learn": schema.Int32Attribute{
							Computed:    true,
							Description: "路由学习开关（1：学习，0：不学习）",
						},
						"route_sync": schema.Int32Attribute{
							Computed:    true,
							Description: "路由同步开关（1：同步，0：不同步）",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间，为UTC格式",
						},
						"subnets": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"subnet_id": schema.StringAttribute{
										Computed:    true,
										Description: "子网ID",
									},
									"ip_version": schema.StringAttribute{
										Computed:    true,
										Description: "ip版本，ipv4或ipv6",
									},
									"subnet_name": schema.StringAttribute{
										Computed:    true,
										Description: "子网名称",
									},
									"cidr": schema.StringAttribute{
										Computed:    true,
										Description: "子网CIDR",
									},
								},
							},
							Description: "子网列表",
						},
					},
				},
				Description: "VPC实例列表",
			},
		},
	}
}

func (c *CtyunExpressConnectionVpcInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunExpressConnectionVpcInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	params := &ec.EcEcListVPCNetworkRequest{
		EcID: config.EcID.ValueString(),
	}
	if !config.CgwID.IsNull() {
		params.CgwID = config.CgwID.ValueStringPointer()
	}
	if !config.ID.IsNull() {
		params.InstanceID = config.ID.ValueStringPointer()
	}
	if !config.QueryContent.IsNull() {
		params.QueryContent = config.QueryContent.ValueStringPointer()
	}
	if !config.Status.IsNull() {
		params.Status = config.Status.ValueStringPointer()
	}
	if !config.IsAuth.IsNull() {
		params.IsAuth = config.IsAuth.ValueInt32Pointer()
	}
	if !config.IsExclusive.IsNull() {
		params.IsExclusive = config.IsExclusive.ValueInt32Pointer()
	}
	if !config.PageNo.IsNull() {
		params.PageNo = config.PageNo.ValueInt32Pointer()
	}
	if !config.PageSize.IsNull() {
		params.PageSize = config.PageSize.ValueInt32Pointer()
	}

	resp, err := c.meta.Apis.SdkEcApis.EcEcListVPCNetworkApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp == nil {
		err = fmt.Errorf("获取VPC网络实例失败（id=%s），接口返回nil，请联系研发确认问题原因！", config.ID.ValueString())
		return
	} else if *resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf(" API return error. Message: %s Error: %s", *resp.ErrorCode, *resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var vpcInstances []CtyunExpressConnectionVpcInstanceModel
	for _, vpcInstanceItem := range resp.ReturnObj.Results {
		var vpcInstance CtyunExpressConnectionVpcInstanceModel
		vpcInstance.EcID = types.StringValue(*vpcInstanceItem.EcID)
		vpcInstance.CgwID = types.StringValue(*vpcInstanceItem.CgwID)
		vpcInstance.CgwName = types.StringValue(*vpcInstanceItem.CgwName)
		vpcInstance.Status = types.StringValue(*vpcInstanceItem.Status)
		vpcInstance.VpcName = types.StringValue(*vpcInstanceItem.VpcName)
		vpcInstance.VpcID = types.StringValue(*vpcInstanceItem.VpcID)
		vpcInstance.VpcCIDR = types.StringValue(*vpcInstanceItem.VpcCIDR)
		vpcInstance.RegionID = types.StringValue(*vpcInstanceItem.DcID)
		vpcInstance.RegionName = types.StringValue(*vpcInstanceItem.DcName)
		vpcInstance.RegionType = types.StringValue(*vpcInstanceItem.DcType)
		vpcInstance.RtbID = types.StringValue(*vpcInstanceItem.RtbID)
		vpcInstance.RtbName = types.StringValue(*vpcInstanceItem.RtbName)
		vpcInstance.ExclusiveID = types.StringValue(*vpcInstanceItem.ExclusiveID)
		vpcInstance.ID = types.StringValue(*vpcInstanceItem.InstanceID)
		vpcInstance.IsAuth = types.Int32Value(*vpcInstanceItem.IsAuth)
		vpcInstance.IsExclusive = types.Int32Value(*vpcInstanceItem.IsExclusive)
		vpcInstance.RouteLearn = types.Int32Value(*vpcInstanceItem.RouteLearn)
		vpcInstance.RouteSync = types.Int32Value(*vpcInstanceItem.RouteSync)
		vpcInstance.CreateTime = types.StringValue(*vpcInstanceItem.CreateDate)

		var subnets []CtyunExpressConnectionVpcInstanceSubnetModel
		for _, subnetItem := range vpcInstanceItem.SubnetList {
			subnets = append(subnets, CtyunExpressConnectionVpcInstanceSubnetModel{
				SubnetID:   types.StringValue(*subnetItem.SubnetID),
				IpVersion:  types.StringValue(*subnetItem.IPVersion),
				Cidr:       types.StringValue(*subnetItem.CIDR),
				SubnetName: types.StringValue(*subnetItem.SubnetName),
			})
		}
		vpcInstance.Subnets = subnets
		vpcInstances = append(vpcInstances, vpcInstance)
	}
	config.VpcInstances = vpcInstances
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

type CtyunExpressConnectionVpcInstancesConfig struct {
	EcID         types.String                             `tfsdk:"ec_id"`
	CgwID        types.String                             `tfsdk:"cgw_id"`
	QueryContent types.String                             `tfsdk:"query_content"`
	Status       types.String                             `tfsdk:"status"`
	IsAuth       types.Int32                              `tfsdk:"is_auth"`
	IsExclusive  types.Int32                              `tfsdk:"is_exclusive"`
	PageNo       types.Int32                              `tfsdk:"page_no"`
	PageSize     types.Int32                              `tfsdk:"page_size"`
	VpcInstances []CtyunExpressConnectionVpcInstanceModel `tfsdk:"vpc_instances"`
	ID           types.String                             `tfsdk:"id"`
}

type CtyunExpressConnectionVpcInstanceModel struct {
	EcID        types.String                                   `tfsdk:"ec_id"`
	CgwID       types.String                                   `tfsdk:"cgw_id"`
	CgwName     types.String                                   `tfsdk:"cgw_name"`
	Status      types.String                                   `tfsdk:"status"`
	VpcName     types.String                                   `tfsdk:"vpc_name"`
	VpcID       types.String                                   `tfsdk:"vpc_id"`
	VpcCIDR     types.String                                   `tfsdk:"vpc_cidr"`
	RegionID    types.String                                   `tfsdk:"region_id"`
	RegionName  types.String                                   `tfsdk:"region_name"`
	RegionType  types.String                                   `tfsdk:"region_type"`
	RtbName     types.String                                   `tfsdk:"rtb_name"`
	RtbID       types.String                                   `tfsdk:"rtb_id"`
	ExclusiveID types.String                                   `tfsdk:"exclusive_id"`
	ID          types.String                                   `tfsdk:"id"`
	IsAuth      types.Int32                                    `tfsdk:"is_auth"`
	IsExclusive types.Int32                                    `tfsdk:"is_exclusive"`
	RouteLearn  types.Int32                                    `tfsdk:"route_learn"`
	RouteSync   types.Int32                                    `tfsdk:"route_sync"`
	CreateTime  types.String                                   `tfsdk:"create_time"`
	Subnets     []CtyunExpressConnectionVpcInstanceSubnetModel `tfsdk:"subnets"`
}

type CtyunExpressConnectionVpcInstanceSubnetModel struct {
	SubnetID   types.String `tfsdk:"subnet_id"`
	IpVersion  types.String `tfsdk:"ip_version"`
	Cidr       types.String `tfsdk:"cidr"`
	SubnetName types.String `tfsdk:"subnet_name"`
}
