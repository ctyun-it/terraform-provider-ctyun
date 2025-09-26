package ccse

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ctyunCcseClusters{}
	_ datasource.DataSourceWithConfigure = &ctyunCcseClusters{}
)

type ctyunCcseClusters struct {
	meta *common.CtyunMetadata
}

func NewCtyunCcseClusters() datasource.DataSource {
	return &ctyunCcseClusters{}
}

func (c *ctyunCcseClusters) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ccse_clusters"
}

type CtyunCcseClustersModel struct {
	ID               types.String `tfsdk:"id"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	ClusterVersion   types.String `tfsdk:"cluster_version"`
	DeployMode       types.String `tfsdk:"deploy_mode"`
	PodCidr          types.String `tfsdk:"pod_cidr"`
	VpcID            types.String `tfsdk:"vpc_id"`
	SubnetID         types.String `tfsdk:"subnet_id"`
	NetworkPlugin    types.String `tfsdk:"network_plugin"`
	ContainerRuntime types.String `tfsdk:"container_runtime"`
	Timezone         types.String `tfsdk:"timezone"`
	ClusterSeries    types.String `tfsdk:"cluster_series"`
	KubeProxy        types.String `tfsdk:"kube_proxy"`
	StartPort        types.Int32  `tfsdk:"start_port"`
	EndPort          types.Int32  `tfsdk:"end_port"`
	ClusterStatus    types.String `tfsdk:"cluster_status"`
	BizState         types.Int32  `tfsdk:"biz_state"`
}

type CtyunCcseClustersConfig struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	RegionID    types.String `tfsdk:"region_id"`
	PageNo      types.Int32  `tfsdk:"page_no"`
	PageSize    types.Int32  `tfsdk:"page_size"`

	Records []CtyunCcseClustersModel `tfsdk:"records"`
}

func (c *ctyunCcseClusters) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `-> 详细说明请见文档：https://www.ctyun.cn/document/10083472/10656137`,
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池ID",
			},
			"cluster_name": schema.StringAttribute{
				Optional:    true,
				Description: "集群名称",
			},
			"page_no": schema.Int32Attribute{
				Optional:    true,
				Description: "列表的页码",
			},
			"page_size": schema.Int32Attribute{
				Optional:    true,
				Description: "每页数据量大小，1-50",
				Validators: []validator.Int32{
					int32validator.Between(1, 50),
				},
			},

			"records": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "集群ID",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "虚拟私有云ID",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "子网ID",
						},
						"cluster_name": schema.StringAttribute{
							Computed:    true,
							Description: "集群名字",
						},
						"network_plugin": schema.StringAttribute{
							Computed:    true,
							Description: "网络插件",
						},
						"start_port": schema.Int32Attribute{
							Computed:    true,
							Description: "节点服务开始端口，范围30000-65535",
						},
						"end_port": schema.Int32Attribute{
							Computed:    true,
							Description: "节点服务终止端口，范围30000-65535",
						},
						"pod_cidr": schema.StringAttribute{
							Computed:    true,
							Description: "pod网络cidr",
						},
						"container_runtime": schema.StringAttribute{
							Computed:    true,
							Description: "容器运行时,可选containerd、docker",
						},
						"timezone": schema.StringAttribute{
							Computed:    true,
							Description: "时区，例如Asia/Shanghai (UTC+08:00)",
						},
						"cluster_version": schema.StringAttribute{
							Computed:    true,
							Description: "集群版本，支持1.23.3 ，1.25.6 ，1.27.8，1.29.3",
						},
						"deploy_mode": schema.StringAttribute{
							Computed:    true,
							Description: "部署模式，单可用区为single，多可用区为multi",
						},
						"kube_proxy": schema.StringAttribute{
							Computed:    true,
							Description: "kubeProxy类型：iptables或ipvs。您可查看<a href=\"https://www.ctyun.cn/document/10083472/10915725\">iptables与IPVS如何选择</a>",
						},
						"cluster_series": schema.StringAttribute{
							Computed:    true,
							Description: "集群系列，cce.standard，cce.managed，您可查看<a href=\"https://www.ctyun.cn/document/10083472/10892150\">产品定义</a>选择",
						},
						"cluster_status": schema.StringAttribute{
							Computed:    true,
							Description: "集群状态：creating：创建中。abnormal：异常。normal：正常。create_fail：创建失败。adjust：规模调整中。updating：升级中。suspend：暂停。deleting：删除中。deleted：已删除。delete_fail：删除失败。resetting：节点重置中。resettled：节点已重置。reset_fail：节点重置失败。upgrading：集群升级中。upgrade_fail：集群升级失败。",
						},
						"biz_state": schema.Int32Attribute{
							Computed:    true,
							Description: "业务状态，1：运行中，2：已停止，3：已注销，4：已退订，5：扩容中，6：开通中，7：已取消，9：重启中，10：节点重置中，11：升级中，13：缩容中，14：已过期(冻结、过期)，15：节点升规格中，17：创建失败，18：退订中，19：控制面升配中，20：休眠中，21：唤醒中，22：转订购模式中",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunCcseClusters) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunCcseClustersConfig
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
	params := &ccse2.CcseListClustersRequest{
		RegionId:    regionId,
		ResPoolId:   regionId,
		ClusterName: config.ClusterName.ValueString(),
	}
	pageNo := config.PageNo.ValueInt32()
	pageSize := config.PageSize.ValueInt32()
	if pageNo > 0 {
		params.PageNow = pageNo
	}
	if pageSize > 0 {
		params.PageSize = pageSize
	}
	// 调用API
	resp, err := c.meta.Apis.SdkCcseApis.CcseListClustersApi.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.ReturnObj == nil {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	}
	// 解析返回值
	config.Records = []CtyunCcseClustersModel{}
	for _, cluster := range resp.ReturnObj.Records {
		item := CtyunCcseClustersModel{
			ID:               types.StringValue(cluster.Id),
			ClusterName:      types.StringValue(cluster.ClusterName),
			ClusterVersion:   types.StringValue(cluster.ClusterVersion),
			DeployMode:       types.StringValue(cluster.DeployMode),
			PodCidr:          types.StringValue(cluster.PodCidr),
			VpcID:            types.StringValue(cluster.VpcId),
			SubnetID:         types.StringValue(cluster.SubnetUuid),
			NetworkPlugin:    types.StringValue(cluster.NetworkPlugin),
			ContainerRuntime: types.StringValue(cluster.ContainerRuntime),
			Timezone:         types.StringValue(cluster.Timezone),
			KubeProxy:        types.StringValue(cluster.KubeProxyPattern),
			StartPort:        types.Int32Value(cluster.StartPort),
			EndPort:          types.Int32Value(cluster.EndPort),
			ClusterStatus:    types.StringValue(cluster.ClusterStatus),
			BizState:         types.Int32Value(cluster.BizState),
		}
		switch cluster.ClusterType {
		case 0:
			item.ClusterSeries = types.StringValue("cce.standard")
		case 2:
			item.ClusterSeries = types.StringValue("cce.managed")
		}
		config.Records = append(config.Records, item)
	}

	// 保存到state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunCcseClusters) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}
