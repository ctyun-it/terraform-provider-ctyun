package elb

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ctelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunElbLoadBalancers{}
	_ datasource.DataSourceWithConfigure = &ctyunElbLoadBalancers{}
)

type ctyunElbLoadBalancers struct {
	meta *common.CtyunMetadata
}

func NewElbLoadBalancers() datasource.DataSource {
	return &ctyunElbLoadBalancers{}
}

func (c *ctyunElbLoadBalancers) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_elb_loadbalancers"
}

func (c *ctyunElbLoadBalancers) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10026756/10138703**`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID",
			},
			"ids": schema.StringAttribute{
				Optional:    true,
				Description: "负载均衡ID列表，以,分隔",
			},
			"resource_type": schema.StringAttribute{
				Optional:    true,
				Description: "负载均衡类型: external / internal",
				Validators: []validator.String{
					stringvalidator.OneOf(business.LbResourceType...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "名称",
			},
			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "子网ID",
			},
			"elbs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"az_name": schema.StringAttribute{
							Computed:    true,
							Description: "可用区名称",
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡ID",
						},
						"project_id": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "描述",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC ID",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"port_id": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡实例默认创建port ID",
						},
						"private_ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡实例的内网VIP",
						},
						"ipv6_address": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡实例的IPv6地址",
						},
						"sla_name": schema.StringAttribute{
							Computed:    true,
							Description: "规格名称",
						},
						"delete_protection": schema.BoolAttribute{
							Computed:    true,
							Description: "删除保护。开启，不开启",
						},
						"admin_status": schema.StringAttribute{
							Computed:    true,
							Description: "管理状态: DOWN / ACTIVE",
							Validators: []validator.String{
								stringvalidator.OneOf(business.AdminStatusName...),
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡状态: DOWN / ACTIVE",
							Validators: []validator.String{
								stringvalidator.OneOf(business.AdminStatusName...),
							},
						},
						"created_time": schema.StringAttribute{
							Computed:    true,
							Description: "created_time",
						},
						"updated_time": schema.StringAttribute{
							Computed:    true,
							Description: "更新时间，为UTC格式",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "负载均衡类型: external / internal",
							Validators: []validator.String{
								stringvalidator.OneOf(business.LbResourceType...),
							},
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "名称",
						},
						"region_id": schema.StringAttribute{
							Computed:    true,
							Description: "资源池ID",
						},
						"eip_info": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"resource_id": schema.StringAttribute{
										Computed:    true,
										Description: "计费类资源ID",
									},
									"eip_id": schema.StringAttribute{
										Computed:    true,
										Description: "弹性公网IP的ID",
									},
									"bandwidth": schema.Int64Attribute{
										Computed:    true,
										Description: "弹性公网IP的带宽",
									},
									"is_talk_order": schema.BoolAttribute{
										Computed:    true,
										Description: "是否按需资源",
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

func (c *ctyunElbLoadBalancers) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()

	var config CtyunElbLoadBalancersConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)

	params := &ctelb.CtelbListLoadBalancerRequest{
		RegionID: regionId,
	}
	if !config.IDs.IsNull() {
		params.IDs = config.IDs.ValueString()
	}
	if !config.Name.IsNull() {
		params.Name = config.Name.ValueString()
	}
	if !config.ResourceType.IsNull() {
		params.ResourceType = config.ResourceType.ValueString()
	}
	if !config.SubnetID.IsNull() {
		params.SubnetID = config.SubnetID.ValueString()
	}

	resp, err := c.meta.Apis.SdkCtElbApis.CtelbListLoadBalancerApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode == common.ErrorStatusCode {
		err = fmt.Errorf("API return error. Message: %s Description: %s", resp.Message, resp.Description)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	var elbs []CtyunElbLoadBalancersModel
	// 将数据写入到Elbs中
	for _, elbItem := range resp.ReturnObj {
		var elb CtyunElbLoadBalancersModel
		elb.RegionID = types.StringValue(elbItem.RegionID)
		elb.AzName = types.StringValue(elbItem.AzName)
		elb.ID = types.StringValue(elbItem.ID)
		elb.ProjectID = types.StringValue(elbItem.ProjectID)
		elb.Name = types.StringValue(elbItem.Name)
		elb.Description = types.StringValue(elbItem.Description)
		elb.VpcID = types.StringValue(elbItem.VpcID)
		elb.SubnetID = types.StringValue(elbItem.SubnetID)
		elb.PortID = types.StringValue(elbItem.PortID)
		elb.PrivateIpAddress = types.StringValue(elbItem.PrivateIpAddress)
		elb.Ipv6Address = types.StringValue(elbItem.Ipv6Address)
		elb.SlaName = types.StringValue(elbItem.SlaName)
		elb.AdminStatus = types.StringValue(elbItem.AdminStatus)
		elb.Status = types.StringValue(elbItem.Status)
		elb.ResourceType = types.StringValue(elbItem.ResourceType)
		elb.CreatedTime = types.StringValue(elbItem.CreatedTime)
		elb.UpdatedTime = types.StringValue(elbItem.UpdatedTime)
		EipInfoList := elbItem.EipInfo
		var eipInfos []EipInfoModel
		for _, eipItem := range EipInfoList {
			var eipInfo EipInfoModel
			eipInfo.EipID = types.StringValue(eipItem.EipID)
			eipInfo.Bandwidth = types.Float32Value(float32(eipItem.Bandwidth))
			if eipItem.IsTalkOrder != nil {
				eipInfo.IsTalkOrder = types.BoolValue(*eipItem.IsTalkOrder)
			}
			eipInfos = append(eipInfos, eipInfo)
		}
		elb.EipInfo = eipInfos
		if elbItem.DeleteProtection != nil {
			elb.DeleteProtection = types.BoolValue(*elbItem.DeleteProtection)
		}

		elbs = append(elbs, elb)
	}
	config.Elbs = elbs
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (c *ctyunElbLoadBalancers) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

type CtyunElbLoadBalancersConfig struct {
	RegionID     types.String                 `tfsdk:"region_id"`     //区域ID
	IDs          types.String                 `tfsdk:"ids"`           //负载均衡ID列表，以,分隔
	ResourceType types.String                 `tfsdk:"resource_type"` //负载均衡类型: external / internal
	Name         types.String                 `tfsdk:"name"`          //名称
	SubnetID     types.String                 `tfsdk:"subnet_id"`     //子网ID
	Elbs         []CtyunElbLoadBalancersModel `tfsdk:"elbs"`
}

type CtyunElbLoadBalancersModel struct {
	RegionID         types.String   `tfsdk:"region_id"`          //区域ID
	AzName           types.String   `tfsdk:"az_name"`            //可用区名称
	ID               types.String   `tfsdk:"id"`                 //负载均衡ID
	ProjectID        types.String   `tfsdk:"project_id"`         //项目ID
	Name             types.String   `tfsdk:"name"`               //名称
	Description      types.String   `tfsdk:"description"`        //描述
	VpcID            types.String   `tfsdk:"vpc_id"`             //VPC ID
	SubnetID         types.String   `tfsdk:"subnet_id"`          //子网ID
	PortID           types.String   `tfsdk:"port_id"`            //负载均衡实例默认创建port ID
	PrivateIpAddress types.String   `tfsdk:"private_ip_address"` //负载均衡实例的内网VIP
	Ipv6Address      types.String   `tfsdk:"ipv6_address"`       //负载均衡实例的IPv6地址
	EipInfo          []EipInfoModel `tfsdk:"eip_info"`           //弹性公网IP信息
	SlaName          types.String   `tfsdk:"sla_name"`           //规格名称
	DeleteProtection types.Bool     `tfsdk:"delete_protection"`  //删除保护。开启，不开启
	AdminStatus      types.String   `tfsdk:"admin_status"`       //管理状态: DOWN / ACTIVE
	Status           types.String   `tfsdk:"status"`             //负载均衡状态: DOWN / ACTIVE
	ResourceType     types.String   `tfsdk:"resource_type"`      //负载均衡类型: external / internal
	CreatedTime      types.String   `tfsdk:"created_time"`       //创建时间，为UTC格式
	UpdatedTime      types.String   `tfsdk:"updated_time"`       //更新时间，为UTC格式
	// 查询的参数
}
