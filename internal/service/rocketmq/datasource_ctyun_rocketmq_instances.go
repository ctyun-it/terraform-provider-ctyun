package rocketmq

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/common"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/rocketmq"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

var (
	_ datasource.DataSource              = &ctyunRocketmqInstances{}
	_ datasource.DataSourceWithConfigure = &ctyunRocketmqInstances{}
)

type ctyunRocketmqInstances struct {
	meta *common.CtyunMetadata
}

func NewCtyunRocketmqInstances() datasource.DataSource {
	return &ctyunRocketmqInstances{}
}

func (c *ctyunRocketmqInstances) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_rocketmq_instances"
}

type CtyunRocketmqInstancesModel struct {
	ID            types.String `tfsdk:"id"`
	MachineSpec   types.String `tfsdk:"machine_spec"`   // 机器规格名
	ProdType      types.String `tfsdk:"prod_type"`      // 产品类型
	BillMode      types.String `tfsdk:"bill_mode"`      // 计费模式 1：包周期；2：按需
	InstanceName  types.String `tfsdk:"instance_name"`  // MQ 实例名称
	RunningState  types.String `tfsdk:"running_state"`  // 实例运行状态
	State         types.Int32  `tfsdk:"state"`          // 实例状态编码
	StatusDesc    types.String `tfsdk:"status_desc"`    // 状态描述
	DiskSpace     types.String `tfsdk:"disk_space"`     // 磁盘空间大小
	DiskType      types.String `tfsdk:"disk_type"`      // 磁盘类型
	NodeSize      types.Int32  `tfsdk:"node_size"`      // broker 节点数量
	ClusterType   types.Int32  `tfsdk:"cluster_type"`   // 集群类型编码 1-单机版 2-集群版
	Version       types.String `tfsdk:"version"`        // 版本号
	EngineType    types.String `tfsdk:"engine_type"`    // 引擎类型
	VpcId         types.String `tfsdk:"vpc_id"`         // VPC ID
	NetName       types.String `tfsdk:"net_name"`       // 网络名称（VPC 名称）
	Subnet        types.String `tfsdk:"subnet"`         // 子网名称
	SecurityGroup types.String `tfsdk:"security_group"` // 安全组 ID
	Vip           types.String `tfsdk:"vip"`            // 实例 VIP 地址
	CrtTime       types.String `tfsdk:"crt_time"`       // 创建时间（UTC 格式）
	ModTime       types.String `tfsdk:"mod_time"`       // 最后修改时间（UTC 格式）
}

type CtyunRocketmqInstancesConfig struct {
	RegionID  types.String                  `tfsdk:"region_id"`
	Instances []CtyunRocketmqInstancesModel `tfsdk:"instances"`
}

func (c *ctyunRocketmqInstances) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: utils.FormatDesc("查询 RocketMQ 实例", "分布式消息服务RocketMQ", "https://www.ctyun.cn/document/10000118/10001967"),
		Attributes: map[string]schema.Attribute{
			"region_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "资源池 ID",
			},
			"instances": schema.ListNestedAttribute{
				Description: "实例列表",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "实例 id",
						},
						"machine_spec": schema.StringAttribute{
							Computed:    true,
							Description: "机器规格名",
						},
						"prod_type": schema.StringAttribute{
							Computed:    true,
							Description: "产品类型",
						},
						"bill_mode": schema.StringAttribute{
							Computed:    true,
							Description: "计费模式 1：包周期；2：按需",
						},
						"instance_name": schema.StringAttribute{
							Computed:    true,
							Description: "MQ 实例名称",
						},
						"running_state": schema.StringAttribute{
							Computed:    true,
							Description: "实例运行状态",
						},
						"state": schema.Int32Attribute{
							Computed:    true,
							Description: "实例状态编码 1.运行中 2.已过期 3.已注销 4.已退订 5.扩容中 6.开通中 7.已取消 8.缩容中 9.重启中 10.网络变更中 11.运维恢复 12.运维停止 13.异常中 15.已欠费 -1.变更 101.开通失败",
						},
						"status_desc": schema.StringAttribute{
							Computed:    true,
							Description: "状态描述",
						},
						"disk_space": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘空间大小",
						},
						"disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "磁盘类型",
						},
						"node_size": schema.Int32Attribute{
							Computed:    true,
							Description: "broker 节点数量",
						},
						"cluster_type": schema.Int32Attribute{
							Computed:    true,
							Description: "集群类型编码 1-单机版 2-集群版",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "版本号",
						},
						"engine_type": schema.StringAttribute{
							Computed:    true,
							Description: "引擎类型",
						},
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "VPC ID",
						},
						"net_name": schema.StringAttribute{
							Computed:    true,
							Description: "网络名称（VPC 名称）",
						},
						"subnet": schema.StringAttribute{
							Computed:    true,
							Description: "子网名称",
						},
						"security_group": schema.StringAttribute{
							Computed:    true,
							Description: "安全组 ID",
						},
						"vip": schema.StringAttribute{
							Computed:    true,
							Description: "实例 VIP 地址",
						},
						"crt_time": schema.StringAttribute{
							Computed:    true,
							Description: "创建时间（UTC 格式）",
						},
						"mod_time": schema.StringAttribute{
							Computed:    true,
							Description: "最后修改时间（UTC 格式）",
						},
					},
				},
			},
		},
	}
}

func (c *ctyunRocketmqInstances) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var err error
	defer func() {
		if err != nil {
			response.Diagnostics.AddError(err.Error(), err.Error())
		}
	}()
	var config CtyunRocketmqInstancesConfig
	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}
	regionId := c.meta.GetExtraIfEmpty(config.RegionID.ValueString(), common.ExtraRegionId)
	if regionId == "" {
		err = fmt.Errorf("region_id 不能为空")
		return
	}

	config.RegionID = types.StringValue(regionId)
	err = c.getAll(ctx, &config)
	if err != nil {
		return
	}

	// 保存到 state
	response.Diagnostics.Append(response.State.Set(ctx, &config)...)
}

func (c *ctyunRocketmqInstances) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	meta := request.ProviderData.(*common.CtyunMetadata)
	c.meta = meta
}

// getAll 查询所有实例
func (c *ctyunRocketmqInstances) getAll(ctx context.Context, config *CtyunRocketmqInstancesConfig) (err error) {
	// 组装请求体
	params := &rocketmq.RocketmqInstQueryV3Request{
		RegionId: config.RegionID.ValueString(),
	}
	resp, err := c.meta.Apis.RocketmqApis.RocketmqInstQueryV3Api.Do(ctx, c.meta.SdkCredential, params)
	if err != nil {
		return
	} else if resp.StatusCode != common.NormalStatusCode {
		err = fmt.Errorf("API return error. Message: %s", resp.Message)
		return
	} else if resp.ReturnObj == nil {
		err = common.InvalidReturnObjError
		return
	}
	config.Instances = []CtyunRocketmqInstancesModel{}
	// 解析返回值
	for _, r := range resp.ReturnObj.ProdInstList {
		item := CtyunRocketmqInstancesModel{
			ID:            types.StringValue(r.ProdInstId),
			MachineSpec:   types.StringValue(r.MachineSpec),
			ProdType:      types.StringValue(r.ProdType),
			BillMode:      types.StringValue(map[string]string{"1": "包年包月", "2": "按需计费"}[r.BillMode]),
			InstanceName:  types.StringValue(r.ProdInstName),
			RunningState:  types.StringValue(r.RunningState),
			State:         types.Int32Value(r.State),
			StatusDesc:    types.StringValue(r.RunningState),
			DiskSpace:     types.StringValue(r.DiskSpace),
			DiskType:      types.StringValue(r.DiskType),
			NodeSize:      types.Int32Value(r.NodeSize),
			ClusterType:   types.Int32Value(r.ClusterType),
			Version:       types.StringValue(r.Version),
			EngineType:    types.StringValue(r.EngineType),
			VpcId:         types.StringValue(r.VpcId),
			NetName:       types.StringValue(r.NetName),
			Subnet:        types.StringValue(r.Subnet),
			SecurityGroup: types.StringValue(r.SecurityGroup),
			Vip:           types.StringValue(r.Vip),
		}

		if r.CrtTime != "" && r.CrtTime != "null" {
			item.CrtTime = types.StringValue(utils.ConvertToUTCZ(time.RFC3339, r.CrtTime))
		}
		if r.ModTime != "" && r.ModTime != "null" {
			item.ModTime = types.StringValue(utils.ConvertToUTCZ(time.RFC3339, r.ModTime))
		}

		config.Instances = append(config.Instances, item)
	}
	return
}
