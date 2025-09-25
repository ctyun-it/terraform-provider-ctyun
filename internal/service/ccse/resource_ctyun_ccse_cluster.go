package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/business"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	terraform_extend "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/defaults"
	validator2 "github.com/ctyun-it/terraform-provider-ctyun/internal/extend/terraform/validator"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &ctyunCcseCluster{}
	_ resource.ResourceWithConfigure   = &ctyunCcseCluster{}
	_ resource.ResourceWithImportState = &ctyunCcseCluster{}
)

type ctyunCcseCluster struct {
	meta          *common.CtyunMetadata
	ecsService    *business.EcsService
	ebmService    *business.EbmService
	vpcService    *business.VpcService
	regionService *business.RegionService
}

func NewCtyunCcseCluster() resource.Resource {
	return &ctyunCcseCluster{}
}

func (c *ctyunCcseCluster) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_cluster"
}

type CtyunCcseClusterConfig struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	MasterOrderID      types.String             `tfsdk:"master_order_id"`
	RegionID           types.String             `tfsdk:"region_id"`
	BaseInfo           CtyunCcseClusterBaseInfo `tfsdk:"base_info"`
	SlaveHost          CtyunCcseClusterSlave    `tfsdk:"slave_host"`
	MasterHost         *CtyunCcseClusterMaster  `tfsdk:"master_host"`
	InternalKubeConfig types.String             `tfsdk:"internal_kube_config"`
	ExternalKubeConfig types.String             `tfsdk:"external_kube_config"`

	totalNodeNum int32
}

type CtyunCcseClusterBaseInfo struct {
	ProjectID             types.String `tfsdk:"project_id"`
	VpcID                 types.String `tfsdk:"vpc_id"`
	SubnetID              types.String `tfsdk:"subnet_id"`
	SecurityGroupID       types.String `tfsdk:"security_group_id"`
	ClusterName           types.String `tfsdk:"cluster_name"`
	ClusterDomain         types.String `tfsdk:"cluster_domain"`
	NetworkPlugin         types.String `tfsdk:"network_plugin"`
	StartPort             types.Int32  `tfsdk:"start_port"`
	EndPort               types.Int32  `tfsdk:"end_port"`
	ElbProdCode           types.String `tfsdk:"elb_prod_code"`
	PodCidr               types.String `tfsdk:"pod_cidr"`
	ServiceCidr           types.String `tfsdk:"service_cidr"`
	PodSubnetIdList       []string     `tfsdk:"pod_subnet_id_list"`
	CycleType             types.String `tfsdk:"cycle_type"`
	CycleCount            types.Int64  `tfsdk:"cycle_count"`
	ContainerRuntime      types.String `tfsdk:"container_runtime"`
	Timezone              types.String `tfsdk:"timezone"`
	ClusterVersion        types.String `tfsdk:"cluster_version"`
	DeployType            types.String `tfsdk:"deploy_type"`
	KubeProxy             types.String `tfsdk:"kube_proxy"`
	ClusterSeries         types.String `tfsdk:"cluster_series"`
	SeriesType            types.String `tfsdk:"series_type"`
	AutoRenew             types.Bool   `tfsdk:"auto_renew"`            // 自动续订
	EnableApiServerEip    types.Bool   `tfsdk:"enable_api_server_eip"` // 是否开启ApiServerEip，默认false，若开启将自动创建按需计费类型的eip。
	EnableSnat            types.Bool   `tfsdk:"enable_snat"`           // 是否开启nat网关，默认false，若开启将自动创建按需计费类型的nat网关。
	NatGatewaySpec        types.String `tfsdk:"nat_gateway_spec"`
	InstallAlsCubeEvent   types.Bool   `tfsdk:"install_als_cube_event"`
	InstallAls            types.Bool   `tfsdk:"install_als"`
	InstallCcseMonitor    types.Bool   `tfsdk:"install_ccse_monitor"`
	InstallNginxIngress   types.Bool   `tfsdk:"install_nginx_ingress"`
	NginxIngressLBSpec    types.String `tfsdk:"nginx_ingress_lb_spec"`
	NginxIngressLBNetWork types.String `tfsdk:"nginx_ingress_network"`
	IpVlan                types.Bool   `tfsdk:"ip_vlan"`
	NetworkPolicy         types.Bool   `tfsdk:"network_policy"`
}

type CtyunCcseClusterAzInfo struct {
	AzName types.String `tfsdk:"az_name"`
	Size   types.Int32  `tfsdk:"size"`
}
type CtyunCcseClusterMaster struct {
	ItemDefName types.String             `tfsdk:"item_def_name"`
	SysDisk     *CtyunCcseClusterDisk    `tfsdk:"sys_disk"`
	DataDisks   []CtyunCcseClusterDisk   `tfsdk:"data_disks"`
	AzInfos     []CtyunCcseClusterAzInfo `tfsdk:"az_infos"`
}
type CtyunCcseClusterSlave struct {
	ItemDefName  types.String             `tfsdk:"item_def_name"`
	AzInfos      []CtyunCcseClusterAzInfo `tfsdk:"az_infos"`
	SysDisk      *CtyunCcseClusterDisk    `tfsdk:"sys_disk"`
	DataDisks    []CtyunCcseClusterDisk   `tfsdk:"data_disks"`
	InstanceType types.String             `tfsdk:"instance_type"`
	MirrorID     types.String             `tfsdk:"mirror_id"`
	MirrorName   types.String             `tfsdk:"mirror_name"`
	MirrorType   types.Int32              `tfsdk:"mirror_type"`
}

type CtyunCcseClusterDisk struct {
	Type types.String `tfsdk:"type"`
	Size types.Int32  `tfsdk:"size"`
}

func (c *ctyunCcseCluster) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `**详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137**`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
				Computed:      true,
				Description:   "ID",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "名称",
			},
			"master_order_id": schema.StringAttribute{
				Computed:    true,
				Description: "主订单号",
			},
			"internal_kube_config": schema.StringAttribute{
				Computed:    true,
				Description: "内网连接信息",
			},
			"external_kube_config": schema.StringAttribute{
				Computed:    true,
				Description: "外网连接信息",
			},
			"region_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "资源池ID，如果不填则默认使用provider ctyun中的region_id或环境变量中的CTYUN_REGION_ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtLeast(1),
				},
				Default: defaults.AcquireFromGlobalString(common.ExtraRegionId, true),
			},
			"base_info": schema.SingleNestedAttribute{
				Required:    true,
				Description: "集群基础信息",
				Attributes: map[string]schema.Attribute{
					"project_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "企业项目ID，如果不填则默认使用provider ctyun中的project_id或环境变量中的CTYUN_PROJECT_ID",
						Default:     defaults.AcquireFromGlobalString(common.ExtraProjectId, false),
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							validator2.Project(),
						},
					},
					"vpc_id": schema.StringAttribute{
						Required:    true,
						Description: "虚拟私有云ID",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							validator2.VpcValidate(),
						},
					},
					"subnet_id": schema.StringAttribute{
						Required:    true,
						Description: "子网ID",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							validator2.SubnetValidate(),
						},
					},
					"security_group_id": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "安全组ID，需属于所选vpc。使用自定义安全组时，需要配置如下规则，参考<a href=\"https://www.ctyun.cn/document/10083472/10915714\">集群安全组规则配置</a>",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							validator2.SecurityGroupValidate(),
						},
					},
					"cluster_name": schema.StringAttribute{
						Required:    true,
						Description: "集群名字",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"cluster_domain": schema.StringAttribute{
						Required:    true,
						Description: "集群本地域名",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"network_plugin": schema.StringAttribute{
						Required:    true,
						Description: "网络插件，可选calico和cubecni，calico需要申请白名单。您可查看<a href=\"https://www.ctyun.cn/document/10083472/10520760\">容器网络插件说明</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcsePluginCubecni, business.CcsePluginCalico),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"start_port": schema.Int32Attribute{
						Required:    true,
						Description: "节点服务开始端口，可选范围30000-65535",
						Validators: []validator.Int32{
							int32validator.Between(30000, 65535),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.RequiresReplace(),
						},
					},
					"end_port": schema.Int32Attribute{
						Required:    true,
						Description: "节点服务终止端口，可选范围30000-65535，startPort到endPort范围需大于20",
						Validators: []validator.Int32{
							int32validator.Between(30000, 65535),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.RequiresReplace(),
						},
					},
					"elb_prod_code": schema.StringAttribute{
						Required:    true,
						Description: "ApiServer的ELB类型，支持standardI（标准I型），standardII（标准II型），enhancedI（增强I型），enhancedII（增强II型），higherI（高阶I型），您可查看<a href=\"https://www.ctyun.cn/document/10026756/10032048\">ELB类型规格说明</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseApiServerElbSpecs...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"pod_cidr": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "pod网络cidr，使用cubecni作为网络插件时，podCidr不填，服务端会取vpcCidr。使用calico作为网络插件时，podCidr与vpcCidr和serviceCidr不能重叠。",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("network_plugin"),
								types.StringValue(business.CcsePluginCalico),
							),
							validator2.ConflictsWithEqualString(
								path.MatchRoot("base_info").AtName("network_plugin"),
								types.StringValue(business.CcsePluginCubecni),
							),
							validator2.Cidr(),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"service_cidr": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "服务cidr，默认10.96.0.0/16。网络插件为calico时，podCidr与vpcCidr与serviceCidr不能重叠。选择cubecni时，podCidr（vpcCidr）与serviceCidr不能重叠。",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Default: stringdefault.StaticString("10.96.0.0/16"),
						Validators: []validator.String{
							validator2.Cidr(),
						},
					},
					"pod_subnet_id_list": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "pod子网ID列表，网络插件选择cubecni必传，需要属于所选VPC，最多支持10个子网",
						Validators: []validator.Set{
							validator2.AlsoRequiresEqualSet(
								path.MatchRoot("base_info").AtName("network_plugin"),
								types.StringValue(business.CcsePluginCubecni),
							),
							setvalidator.SizeAtMost(10),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.RequiresReplace(),
						},
					},
					"enable_api_server_eip": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否开启ApiServerEip，默认false，若开启将自动创建按需计费类型的eip。",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"enable_snat": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否开启nat网关，默认false，若开启将自动创建按需计费类型的nat网关。",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"nat_gateway_spec": schema.StringAttribute{
						Optional:    true,
						Description: "当enable_snat=true时填写，nat网关规格：small，medium，large，xlarge，可参考<a href=\"https://www.ctyun.cn/document/10026759/10043996\">产品规格说明</a>",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("small", "medium", "large", "xlarge"),
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("enable_snat"),
								types.BoolValue(true),
							),
						},
					},
					"install_als_cube_event": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否安装事件采集插件，默认false",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"install_als": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否安装日志插件，默认false",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"install_ccse_monitor": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否安装监控插件，默认false",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"install_nginx_ingress": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否安装nginx_ingress插件，默认false",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"nginx_ingress_lb_spec": schema.StringAttribute{
						Optional:    true,
						Description: "install_nginx_ingress=true必填，支持规格：standardI（标准I型） ,standardII（标准II型）, enhancedI（增强I型）, enhancedII（增强II型） , higherI（高阶I型），可参考<a href=\"https://www.ctyun.cn/document/10026756/10032048\">规格详情</a>",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("standardI", "standardII", "enhancedI", "enhancedII", "higherI"),
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("install_nginx_ingress"),
								types.BoolValue(true),
							),
						},
					},
					"nginx_ingress_network": schema.StringAttribute{
						Optional:    true,
						Description: "install_nginx_ingress=true必填，nginx ingress访问方式：external（公网），internal（内网），当选择公网时将自动创建eip额外产生eip相关费用",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("external", "internal"),
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("install_nginx_ingress"),
								types.BoolValue(true),
							),
						},
					},
					"ip_vlan": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "基于IPVLAN做弹性网卡共享，默认false，当指定为true时，主机镜像只有使用CtyunOS系统才能生效",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"network_policy": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "是否提供基于策略的网络访问控制，默认false",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"auto_renew": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "是否自动续订，默认非自动续订，当cycle_type不等于on_demand时才可填写，按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年。",
						Default:     booldefault.StaticBool(false),
						Validators: []validator.Bool{
							validator2.ConflictsWithEqualBool(
								path.MatchRoot("base_info").AtName("cycle_type"),
								types.StringValue(business.OrderCycleTypeOnDemand),
							),
						},
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.RequiresReplace(),
						},
					},
					"cycle_type": schema.StringAttribute{
						Required:    true,
						Description: "订购周期类型，取值范围：month：按月，year：按年、on_demand：按需。当此值为month或者year时，cycle_count为必填",
						Validators: []validator.String{
							stringvalidator.OneOf(business.OrderCycleTypes...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"cycle_count": schema.Int64Attribute{
						Optional:    true,
						Description: "订购时长，该参数在cycle_type为month或year时才生效，当cycle_type=month，支持订购1-11个月；当cycle_type=year，支持订购1-3年",
						Validators: []validator.Int64{
							validator2.AlsoRequiresEqualInt64(
								path.MatchRoot("base_info").AtName("cycle_type"),
								types.StringValue(business.OrderCycleTypeMonth),
								types.StringValue(business.OrderCycleTypeYear),
							),
							validator2.ConflictsWithEqualInt64(
								path.MatchRoot("base_info").AtName("cycle_type"),
								types.StringValue(business.OrderCycleTypeOnDemand),
							),
							validator2.CycleCount(1, 11, 1, 3),
						},
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"container_runtime": schema.StringAttribute{
						Required:    true,
						Description: "容器运行时,可选containerd、docker，您可查看<a href=\"https://www.ctyun.cn/document/10083472/10902208\">容器运行时说明</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseContainerRuntimeContainerd, business.CcseContainerRuntimeDocker),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"timezone": schema.StringAttribute{
						Required:    true,
						Description: "时区，例如Asia/Shanghai (UTC+08:00)",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"cluster_version": schema.StringAttribute{
						Required:    true,
						Description: "集群版本，支持1.31.6，1.29.3，1.27.8，您可查看<a href=\"https://www.ctyun.cn/document/10083472/10650447\">集群版本说明</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseClusterVersions...),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"deploy_type": schema.StringAttribute{
						Required:    true,
						Description: "部署模式，单可用区为single，多可用区为multi",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseDeployTypeSingle, business.CcseDeployTypeMulti),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"kube_proxy": schema.StringAttribute{
						Required:    true,
						Description: "kubeProxy类型：iptables或ipvs。您可查看<a href=\"https://www.ctyun.cn/document/10083472/10915725\">iptables与IPVS如何选择</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseKubeProxyIptables, business.CcseKubeProxyIpvs),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"cluster_series": schema.StringAttribute{
						Required:    true,
						Description: "集群系列，支持cce.standard（专有版），cce.managed（托管版），您可查看<a href=\"https://www.ctyun.cn/document/10083472/10892150\">产品定义</a>",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseClusterSeriesStandard, business.CcseClusterSeriesManaged, business.CcseClusterSeriesIcce),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"series_type": schema.StringAttribute{
						Optional:    true,
						Description: "托管版集群规格，托管版集群必填。支持managedbase（单实例），managedpro（多实例）。单/多实例指控制面是否高可用，生产环境建议使用多实例",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseSeriesTypeManagedbase, business.CcseSeriesTypeManagedpro),
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("base_info").AtName("cluster_series"),
								types.StringValue(business.CcseClusterSeriesManaged),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"master_host": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "master节点基本信息，专有版必填，托管版时不传",
				Attributes: map[string]schema.Attribute{
					"item_def_name": schema.StringAttribute{
						Required:    true,
						Description: "实例规格名称，使用至少4C8G以上的规格，仅支持云主机，可通过ctyun_ecs_flavors查询",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"sys_disk": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "系统盘信息",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "系统盘类型，支持SATA、SAS、SSD",
								Validators: []validator.String{
									stringvalidator.OneOf(business.CcseDiskTypes...),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"size": schema.Int32Attribute{
								Required:    true,
								Description: "系统盘大小，单位为G，支持范围80-2040",
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.RequiresReplace(),
								},
								Validators: []validator.Int32{
									int32validator.Between(80, 2040),
								},
							},
						},
					},
					"data_disks": schema.ListNestedAttribute{
						Optional:    true,
						Description: "数据盘信息",
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "数据盘类型，支持SATA、SAS、SSD",
									Validators: []validator.String{
										stringvalidator.OneOf(business.CcseDiskTypes...),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "数据盘大小，单位为G，支持范围10-20000",
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.Between(10, 20000),
									},
								},
							},
						},
					},
					"az_infos": schema.ListNestedAttribute{
						Required: true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
						Description: "可用区信息，包括可用区编码和该可用区下master节点数量，支持的可用区可通过ctyun_regions查询",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"az_name": schema.StringAttribute{
									Required:    true,
									Description: "master可用区编码",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "该可用区下master节点数量",
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.AtLeast(1),
									},
								},
							},
						},
					},
				},
				Validators: []validator.Object{
					validator2.AlsoRequiresEqualObject(
						path.MatchRoot("base_info").AtName("cluster_series"),
						types.StringValue(business.CcseClusterSeriesStandard),
					),
					validator2.ConflictsWithEqualObject(
						path.MatchRoot("base_info").AtName("cluster_series"),
						types.StringValue(business.CcseClusterSeriesManaged),
						types.StringValue(business.CcseClusterSeriesIcce),
					),
				},
			},
			"slave_host": schema.SingleNestedAttribute{
				Required:    true,
				Description: "slave节点基本信息",
				Attributes: map[string]schema.Attribute{
					"instance_type": schema.StringAttribute{
						Required:    true,
						Description: "实例类型，支持ecs（云主机）、ebm（裸金属）",
						Validators: []validator.String{
							stringvalidator.OneOf(business.CcseSlaveInstanceTypeEcs, business.CcseSlaveInstanceTypeEbm),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"mirror_id": schema.StringAttribute{
						Optional:    true,
						Description: "镜像id，worker节点为ecs类型必填，可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEcs),
							),
							validator2.ConflictsWithEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEbm),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"mirror_name": schema.StringAttribute{
						Optional:    true,
						Description: "镜像名称，worker节点为ebm类型必填，可查看<a href=\"https://www.ctyun.cn/document/10083472/11004475\">节点规格和节点镜像</a>",
						Validators: []validator.String{
							validator2.AlsoRequiresEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEbm),
							),
							validator2.ConflictsWithEqualString(
								path.MatchRoot("slave_host").AtName("instance_type"),
								types.StringValue(business.CcseSlaveInstanceTypeEcs),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"mirror_type": schema.Int32Attribute{
						Required:    true,
						Description: "镜像类型，支持传0（私有），1（公有），可查看<a href=\"https://www.ctyun.cn/document/10026730/10030151\">镜像概述</a>",
						Validators: []validator.Int32{
							int32validator.Between(0, 1),
						},
						PlanModifiers: []planmodifier.Int32{
							int32planmodifier.RequiresReplace(),
						},
					},
					"item_def_name": schema.StringAttribute{
						Required:    true,
						Description: "实例规格名称，使用至少4C8G以上的规格，云主机规格通过ctyun_ecs_flavors查询，裸金属规格通过ctyun_ebm_device_types查询",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"az_infos": schema.ListNestedAttribute{
						Required: true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
						Description: "可用区信息，包括可用区编码和该可用区下worker节点数量，支持的可用区可通过ctyun_regions查询",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"az_name": schema.StringAttribute{
									Required:    true,
									Description: "worker可用区编码",
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "该可用区下worker节点数量",
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.AtLeast(1),
									},
								},
							},
						},
					},
					"sys_disk": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "系统盘信息",
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required:    true,
								Description: "系统盘类型，支持SATA、SAS、SSD",
								Validators: []validator.String{
									stringvalidator.OneOf(business.CcseDiskTypes...),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"size": schema.Int32Attribute{
								Required:    true,
								Description: "系统盘大小，单位为G，支持范围80-2040",
								PlanModifiers: []planmodifier.Int32{
									int32planmodifier.RequiresReplace(),
								},
								Validators: []validator.Int32{
									int32validator.Between(80, 2040),
								},
							},
						},
					},
					"data_disks": schema.ListNestedAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.RequiresReplace(),
						},
						Description: "数据盘信息",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "数据盘类型，支持SATA、SAS、SSD",
									Validators: []validator.String{
										stringvalidator.OneOf(business.CcseDiskTypes...),
									},
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.RequiresReplace(),
									},
								},
								"size": schema.Int32Attribute{
									Required:    true,
									Description: "数据盘大小，单位为G，支持范围10-20000",
									PlanModifiers: []planmodifier.Int32{
										int32planmodifier.RequiresReplace(),
									},
									Validators: []validator.Int32{
										int32validator.Between(10, 20000),
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

func (c *ctyunCcseCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var plan CtyunCcseClusterConfig
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 创建前检查
	err = c.checkBeforeCreate(ctx, plan)
	if err != nil {
		return
	}
	// 创建
	masterOrderID, err := c.create(ctx, &plan)
	if err != nil {
		return
	}
	plan.MasterOrderID = types.StringValue(masterOrderID)
	// 创建后检查
	id, err := c.checkAfterCreate(ctx, plan)
	if err != nil {
		if strings.Contains(err.Error(), "初始节点创建时间过长") {
			plan.ID = types.StringValue(id)
			response.Diagnostics.Append(response.State.Set(ctx, plan)...)
		}
		return
	}
	plan.ID = types.StringValue(id)
	// 反查信息
	err = c.getAndMerge(ctx, &plan)
	if err != nil {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, plan)...)
}

func (c *ctyunCcseCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseClusterConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 查询远端
	err = c.getAndMerge(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "退订状态") {
			err = nil
			response.State.RemoveResource(ctx)
		}
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (c *ctyunCcseCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (c *ctyunCcseCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var state CtyunCcseClusterConfig
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	// 删除
	err = c.delete(ctx, state)
	if err != nil {
		return
	}
	err = c.checkAfterDelete(ctx, state)
	if err != nil {
		return
	}
	response.Diagnostics.AddWarning("删除CCSE集群成功", "集群退订后，若立即删除子网或安全组可能会失败，需要等待底层资源释放")
}

// 导入命令：terraform import [配置标识].[导入配置名称],[clusterID],[regionID]
func (c *ctyunCcseCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var cfg CtyunCcseClusterConfig
	var clusterID, regionID string
	err = terraform_extend.Split(request.ID, &clusterID, &regionID)
	if err != nil {
		return
	}
	cfg.RegionID = types.StringValue(regionID)
	cfg.ID = types.StringValue(clusterID)
	// 查询远端
	err = c.getAndMerge(ctx, &cfg)
	if err != nil {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, cfg)...)
}

func (c *ctyunCcseCluster) Configure(_ context.Context, request resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
	c.ecsService = business.NewEcsService(c.meta)
	c.ebmService = business.NewEbmService(c.meta)
	c.vpcService = business.NewVpcService(c.meta)
	c.regionService = business.NewRegionService(c.meta)
}

// checkBeforeCreate 创建前检查
func (c *ctyunCcseCluster) checkBeforeCreate(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	// 确保当前虚拟私有云存在，且子网与虚拟私有云存在对应关系
	vpc, regionID, projectID := plan.BaseInfo.VpcID.ValueString(), plan.RegionID.ValueString(), plan.BaseInfo.ProjectID.ValueString()
	subnets, err := c.vpcService.GetVpcSubnet(ctx, vpc, regionID, projectID)
	if err != nil {
		return err
	}
	if _, ok := subnets[plan.BaseInfo.SubnetID.ValueString()]; !ok {
		return fmt.Errorf("子网 %s 不在 %s 内", plan.BaseInfo.SubnetID.ValueString(), vpc)
	}
	for s, _ := range subnets {
		if _, ok := subnets[s]; !ok {
			return fmt.Errorf("子网 %s 不在 %s 内", plan.BaseInfo.SubnetID.ValueString(), vpc)
		}
	}
	return
}

// create 创建
func (c *ctyunCcseCluster) create(ctx context.Context, plan *CtyunCcseClusterConfig) (masterOrderID string, err error) {
	params := &ccse2.CcseCreateClusterRequest{
		RegionId:  plan.RegionID.ValueString(),
		ResPoolId: plan.RegionID.ValueString(),
	}
	auto := true
	// 处理 clusterBaseInfo
	clusterBaseInfo := ccse2.CcseCreateClusterClusterBaseInfoRequest{
		ProjectId:                 plan.BaseInfo.ProjectID.ValueString(),
		VpcUuid:                   plan.BaseInfo.VpcID.ValueString(),
		SubnetUuid:                plan.BaseInfo.SubnetID.ValueString(),
		AutoGenerateSecurityGroup: &auto,
		ClusterName:               plan.BaseInfo.ClusterName.ValueString(),
		ClusterDomain:             plan.BaseInfo.ClusterDomain.ValueString(),
		ClusterVersion:            plan.BaseInfo.ClusterVersion.ValueString(),
		ClusterSeries:             plan.BaseInfo.ClusterSeries.ValueString(),
		NetworkPlugin:             plan.BaseInfo.NetworkPlugin.ValueString(),
		StartPort:                 int64(plan.BaseInfo.StartPort.ValueInt32()),
		EndPort:                   int64(plan.BaseInfo.EndPort.ValueInt32()),
		ElbProdCode:               plan.BaseInfo.ElbProdCode.ValueString(),
		PodSubnetUuidList:         plan.BaseInfo.PodSubnetIdList,
		PodCidr:                   plan.BaseInfo.PodCidr.ValueString(),
		ServiceCidr:               plan.BaseInfo.ServiceCidr.ValueString(),
		ContainerRuntime:          plan.BaseInfo.ContainerRuntime.ValueString(),
		Timezone:                  plan.BaseInfo.Timezone.ValueString(),
		DeployType:                plan.BaseInfo.DeployType.ValueString(),
		KubeProxy:                 plan.BaseInfo.KubeProxy.ValueString(),
		SeriesType:                plan.BaseInfo.SeriesType.ValueString(),
		EnableApiServerEip:        plan.BaseInfo.EnableApiServerEip.ValueBoolPointer(),
		EnableSnat:                plan.BaseInfo.EnableSnat.ValueBoolPointer(),
		NatGatewaySpec:            plan.BaseInfo.NatGatewaySpec.ValueString(),
		EnableAlsCubeEventer:      plan.BaseInfo.InstallAlsCubeEvent.ValueBoolPointer(),
		EnableAls:                 plan.BaseInfo.InstallAls.ValueBoolPointer(),
		PluginCcseMonitorEnabled:  plan.BaseInfo.InstallCcseMonitor.ValueBoolPointer(),
		InstallNginxIngress:       plan.BaseInfo.InstallNginxIngress.ValueBoolPointer(),
		NginxIngressLBSpec:        plan.BaseInfo.NginxIngressLBSpec.ValueString(),
		Ipvlan:                    plan.BaseInfo.IpVlan.ValueBoolPointer(),
		NetworkPolicy:             plan.BaseInfo.NetworkPolicy.ValueBoolPointer(),
		NginxIngressLBNetWork:     plan.BaseInfo.NginxIngressLBNetWork.ValueString(),
	}
	if plan.BaseInfo.SeriesType.ValueString() == business.CcseSeriesTypeManagedpro {
		clusterBaseInfo.NodeScale = "50"
	}
	if plan.BaseInfo.SecurityGroupID.ValueString() != "" {
		f := false
		clusterBaseInfo.AutoGenerateSecurityGroup = &f
		clusterBaseInfo.SecurityGroupUuid = plan.BaseInfo.SecurityGroupID.ValueString()
	}

	switch plan.BaseInfo.CycleType.ValueString() {
	case business.OnDemandCycleType:
		clusterBaseInfo.BillMode = "2"
	case business.MonthCycleType:
		clusterBaseInfo.AutoRenewStatus = plan.BaseInfo.AutoRenew.ValueBoolPointer()
		clusterBaseInfo.BillMode = "1"
		clusterBaseInfo.CycleType = "3"
		clusterBaseInfo.CycleCnt = int32(plan.BaseInfo.CycleCount.ValueInt64())
	case business.YearCycleType:
		clusterBaseInfo.AutoRenewStatus = plan.BaseInfo.AutoRenew.ValueBoolPointer()
		clusterBaseInfo.BillMode = "1"
		clusterBaseInfo.CycleType = fmt.Sprintf("%d", plan.BaseInfo.CycleCount.ValueInt64()+4) // 1年传5，2年传6，3年传7
		clusterBaseInfo.CycleCnt = 1
	}

	// 处理masterHost
	if plan.MasterHost != nil {
		flavorName := plan.MasterHost.ItemDefName.ValueString()
		var flavor ctecs.EcsFlavorListFlavorListResponse
		flavor, err = c.ecsService.GetFlavorByName(ctx, flavorName, plan.RegionID.ValueString())
		if err != nil {
			return
		}
		if flavor.FlavorCpu < 4 || flavor.FlavorRam < 8 {
			err = fmt.Errorf("master节点的规格至少需要4c8g")
		}
		var totalSize int32
		for _, az := range plan.MasterHost.AzInfos {
			clusterBaseInfo.AzInfos = append(clusterBaseInfo.AzInfos, &ccse2.CcseCreateClusterClusterBaseInfoAzInfosRequest{
				AzName: az.AzName.ValueString(),
				Size:   az.Size.ValueInt32(),
			})
			totalSize += az.Size.ValueInt32()
		}
		plan.totalNodeNum += totalSize

		masterHost := ccse2.CcseCreateClusterMasterHostRequest{
			Cpu:         int32(flavor.FlavorCpu),
			Mem:         int32(flavor.FlavorRam),
			ItemDefName: flavorName,
			ItemDefType: flavor.FlavorType,
			Size:        totalSize,
		}
		if plan.MasterHost.SysDisk != nil {
			masterHost.SysDisk = &ccse2.CcseCreateClusterMasterHostSysDiskRequest{
				ItemDefName: plan.MasterHost.SysDisk.Type.ValueString(),
				Size:        plan.MasterHost.SysDisk.Size.ValueInt32(),
			}
		}
		for _, disk := range plan.MasterHost.DataDisks {
			masterHost.DataDisks = append(masterHost.DataDisks, &ccse2.CcseCreateClusterMasterHostDataDisksRequest{
				ItemDefName: disk.Type.ValueString(),
				Size:        disk.Size.ValueInt32(),
			})
		}
		params.MasterHost = &masterHost
	} else {
		// 通过资源池查所有可用区名称，然后组装azInfos
		zones, err := c.regionService.GetZonesByRegionID(ctx, plan.RegionID.ValueString())
		if err != nil {
			return "", err
		}
		for _, az := range zones {
			clusterBaseInfo.AzInfos = append(clusterBaseInfo.AzInfos, &ccse2.CcseCreateClusterClusterBaseInfoAzInfosRequest{
				AzName: az,
			})
		}
	}

	// 处理slaveHost

	slaveHost := ccse2.CcseCreateClusterSlaveHostRequest{
		Size:       0,
		MirrorType: plan.SlaveHost.MirrorType.ValueInt32(),
	}
	if plan.SlaveHost.SysDisk != nil {
		slaveHost.SysDisk = &ccse2.CcseCreateClusterSlaveHostSysDiskRequest{
			ItemDefName: plan.SlaveHost.SysDisk.Type.ValueString(),
			Size:        plan.SlaveHost.SysDisk.Size.ValueInt32(),
		}
	}

	for _, disk := range plan.SlaveHost.DataDisks {
		slaveHost.DataDisks = append(slaveHost.DataDisks, &ccse2.CcseCreateClusterSlaveHostDataDisksRequest{
			ItemDefName: disk.Type.ValueString(),
			Size:        disk.Size.ValueInt32(),
		})
	}

	var azName string
	for _, az := range plan.SlaveHost.AzInfos {
		if azName == "" {
			azName = az.AzName.ValueString()
		}
		slaveHost.AzInfos = append(slaveHost.AzInfos, &ccse2.CcseCreateClusterSlaveHostAzInfosRequest{
			AzName: az.AzName.ValueString(),
			Size:   az.Size.ValueInt32(),
		})
		slaveHost.Size += az.Size.ValueInt32()
	}
	plan.totalNodeNum += slaveHost.Size
	switch plan.SlaveHost.InstanceType.ValueString() {
	case business.CcseSlaveInstanceTypeEcs:
		slaveHost.ForeignMirrorId = plan.SlaveHost.MirrorID.ValueString()
		flavorName := plan.SlaveHost.ItemDefName.ValueString()
		flavor, err := c.ecsService.GetFlavorByName(ctx, flavorName, plan.RegionID.ValueString())
		if err != nil {
			return "", err
		}
		slaveHost.ItemDefName = flavorName
		slaveHost.ItemDefType = flavor.FlavorType
	case business.CcseSlaveInstanceTypeEbm:
		slaveHost.MirrorName = plan.SlaveHost.MirrorName.ValueString()
		deviceType := plan.SlaveHost.ItemDefName.ValueString()
		flavor, err := c.ebmService.GetDeviceType(ctx, deviceType, plan.RegionID.ValueString(), azName)
		if err != nil {
			return "", err
		}
		slaveHost.ItemDefName = deviceType
		slaveHost.ItemDefType = deviceType

		if !utils.SecBool(flavor.CloudBoot) && slaveHost.SysDisk != nil {
			return "", fmt.Errorf("裸金属规格 %s 不支持自定义系统盘", deviceType)
		}
		if !utils.SecBool(flavor.SupportCloud) && len(slaveHost.DataDisks) > 0 {
			return "", fmt.Errorf("裸金属规格 %s 不支持自定义数据盘", deviceType)
		}
	}

	params.ClusterBaseInfo = &clusterBaseInfo
	params.SlaveHost = &slaveHost

	resp, err := c.meta.Apis.SdkCcseApis.CcseCreateClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	masterOrderID = resp.ReturnObj.OrderId
	return
}

// getAndMerge 从远端查询
func (c *ctyunCcseCluster) getAndMerge(ctx context.Context, plan *CtyunCcseClusterConfig) (err error) {
	params := &ccse2.CcseGetClusterRequest{
		RegionId:  plan.RegionID.ValueString(),
		ClusterId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}

	instance := resp.ReturnObj
	if instance.BizState == business.CcseRefundedBizState || instance.BizState == business.CcseRefundingBizState {
		return fmt.Errorf("集群 %s 处于退订状态", plan.ID.ValueString())
	}
	plan.Name = types.StringValue(instance.ClusterName)
	plan.BaseInfo.VpcID = types.StringValue(instance.VpcId)
	plan.BaseInfo.SecurityGroupID = types.StringValue(instance.SecurityGroupId)
	plan.BaseInfo.SubnetID = types.StringValue(instance.SubnetUuid)
	plan.BaseInfo.NetworkPlugin = types.StringValue(instance.NetworkPlugin)
	plan.BaseInfo.PodCidr = types.StringValue(instance.PodCidr)
	plan.BaseInfo.ServiceCidr = types.StringValue(instance.ServiceCidr)
	plan.BaseInfo.Timezone = types.StringValue(instance.Timezone)
	plan.BaseInfo.ClusterVersion = types.StringValue(instance.ClusterVersion)
	plan.BaseInfo.KubeProxy = types.StringValue(instance.KubeProxyPattern)
	switch instance.ClusterType {
	case 0:
		plan.BaseInfo.ClusterSeries = types.StringValue(business.CcseClusterSeriesStandard)
	case 2:
		plan.BaseInfo.ClusterSeries = types.StringValue(business.CcseClusterSeriesManaged)
	case 4:
		plan.BaseInfo.ClusterSeries = types.StringValue(business.CcseClusterSeriesIcce)
	}
	plan.BaseInfo.StartPort = types.Int32Value(instance.StartPort)
	plan.BaseInfo.EndPort = types.Int32Value(instance.EndPort)

	internalKubeConfig, err := c.getKubeConfig(ctx, *plan, false)
	if err != nil {
		return
	}
	externalKubeConfig, err := c.getKubeConfig(ctx, *plan, true)
	if err != nil {
		return
	}
	plan.InternalKubeConfig = types.StringValue(internalKubeConfig)
	plan.ExternalKubeConfig = types.StringValue(externalKubeConfig)
	return
}

// delete 删除
func (c *ctyunCcseCluster) delete(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	params := &ccse2.CcseDeleteClusterRequest{
		RegionId:   plan.RegionID.ValueString(),
		ResPoolId:  plan.RegionID.ValueString(),
		ProdInstId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseDeleteClusterApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s RequestId: %s", resp.Message, resp.RequestId)
		return
	}
	return
}

// listByName 根据名称查询集群
func (c *ctyunCcseCluster) listByName(ctx context.Context, plan CtyunCcseClusterConfig) (clusters []*ccse2.CcseListClustersReturnObjRecordsResponse, err error) {
	params := &ccse2.CcseListClustersRequest{
		RegionId:    plan.RegionID.ValueString(),
		ClusterName: plan.BaseInfo.ClusterName.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListClustersApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	clusters = resp.ReturnObj.Records
	return
}

// checkAfterCreate 创建后检查
func (c *ctyunCcseCluster) checkAfterCreate(ctx context.Context, plan CtyunCcseClusterConfig) (id string, err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var clusters []*ccse2.CcseListClustersReturnObjRecordsResponse
			clusters, err = c.listByName(ctx, plan)
			if err != nil {
				return false
			}
			if len(clusters) == 0 || clusters[0].BizState != 1 || clusters[0].Id == "" {
				return true
			}

			id = clusters[0].Id
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("创建时间过长")
		return
	}
	plan.ID = types.StringValue(id)
	// 集群创建成功，还需要检查一下节点是否都成功
	err = c.checkNodeStatus(ctx, plan)
	if err != nil {
		return
	}
	return
}

// checkNodeStatus 检查节点状态
func (c *ctyunCcseCluster) checkNodeStatus(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var nodes []*ccse2.CcseListClusterNodesReturnObjResponse
			nodes, err = c.listNode(ctx, plan)
			if err != nil {
				return false
			}
			if len(nodes) < int(plan.totalNodeNum) {
				return true
			}
			for _, n := range nodes {
				if n.NodeStatus != "normal" {
					return true
				}
			}

			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("初始节点创建时间过长")
		return
	}
	return
}

// listNode 获取节点列表
func (c *ctyunCcseCluster) listNode(ctx context.Context, plan CtyunCcseClusterConfig) (nodes []*ccse2.CcseListClusterNodesReturnObjResponse, err error) {
	params := &ccse2.CcseListClusterNodesRequest{
		RegionId:  plan.RegionID.ValueString(),
		ClusterId: plan.ID.ValueString(),
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseListClusterNodesApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	nodes = resp.ReturnObj
	return
}

// checkAfterDelete 删除后检查
func (c *ctyunCcseCluster) checkAfterDelete(ctx context.Context, plan CtyunCcseClusterConfig) (err error) {
	var executeSuccessFlag bool
	retryer, _ := business.NewRetryer(time.Second*10, 180)
	retryer.Start(
		func(currentTime int) bool {
			var clusters []*ccse2.CcseListClustersReturnObjRecordsResponse
			clusters, err = c.listByName(ctx, plan)
			if err != nil {
				return false
			}
			if len(clusters) != 0 && clusters[0].BizState != 4 {
				return true
			}
			executeSuccessFlag = true
			return false
		})
	if err != nil {
		return
	}
	if !executeSuccessFlag {
		err = fmt.Errorf("删除时间过长")
	}
	return
}

func (c *ctyunCcseCluster) getKubeConfig(ctx context.Context, plan CtyunCcseClusterConfig, isPublic bool) (config string, err error) {
	params := &ccse2.CcseGetClusterKubeConfigRequest{
		RegionId:  plan.RegionID.ValueString(),
		ClusterId: plan.ID.ValueString(),
		IsPublic:  isPublic,
	}
	resp, err := c.meta.Apis.SdkCcseApis.CcseGetClusterKubeConfigApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config = resp.ReturnObj.KubeConfig
	return
}
