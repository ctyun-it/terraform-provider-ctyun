package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqInstQueryV3Api
/* 查询租户实例 v3
 */type RocketmqInstQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqInstQueryV3Api(client *core.CtyunClient) *RocketmqInstQueryV3Api {
	return &RocketmqInstQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *RocketmqInstQueryV3Api) Do(ctx context.Context, credential core.Credential, req *RocketmqInstQueryV3Request) (*RocketmqInstQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp RocketmqInstQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqInstQueryV3Request struct {
	RegionId string `json:"regionId,omitempty"` /*  资源池 ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
}

type RocketmqInstQueryV3Response struct {
	StatusCode int32                            `json:"statusCode"` // 接口系统层面状态码。成功："800"，失败："900"
	Message    string                           `json:"message"`    // 接口调用状态描述。成功时为"success"，失败时为具体失败信息
	ReturnObj  *RocketmqInstancesQueryReturnObj `json:"returnObj"`  // 核心返回对象。成功时包含 MQ 实例列表数据，失败时为空对象
	Error      string                           `json:"error"`      // 错误码。仅失败时返回，描述具体错误信息
}

type RocketmqInstancesQueryReturnObj struct {
	ProdInstList []*ProdInstInfo `json:"prodInstList"` // MQ 实例列表数据。仅接口调用成功时返回
}

type ProdInstInfo struct {
	ProdInstId        string             `json:"prodInstId"`        // 实例 id
	BillMode          string             `json:"billMode"`          // 计费模式 1：包周期；2：按需
	ProdInstName      string             `json:"prodInstName"`      // MQ 实例名称
	RunningState      string             `json:"runningState"`      // 实例运行状态
	State             int32              `json:"state"`             // 实例状态编码 1.运行中 2.已过期 3.已注销 4.已退订 5.扩容中 6.开通中 7.已取消 8.缩容中 9.重启中 10.网络变更中 11.运维恢复 12.运维停止 13.异常中 15.已欠费 -1.变更 101.开通失败
	ProdType          string             `json:"prodType"`          // 产品类型
	MachineSpec       string             `json:"machineSpec"`       // 机器规格名
	TopicLimit        string             `json:"topicLimit"`        // Topic 数量限制
	DiskSpace         string             `json:"diskSpace"`         // 磁盘空间大小
	NetName           string             `json:"netName"`           // 网络名称（VPC 名称）
	Subnet            string             `json:"subnet"`            // 子网名称
	SecurityGroup     string             `json:"securityGroup"`     // 安全组 ID
	ModTime           string             `json:"modTime"`           // 最后修改时间（UTC 格式）
	DeployType        string             `json:"deployType"`        // 部署类型编码 1-虚机部署
	DiskType          string             `json:"diskType"`          // 磁盘类型
	DiskIsEncrypt     bool               `json:"diskIsEncrypt"`     // 磁盘是否加密
	FileReservedTime  int32              `json:"fileReservedTime"`  // 文件保留时间（单位：小时）
	ClusterProperties *ClusterProperties `json:"clusterProperties"` // 集群配置属性
	EngineType        string             `json:"engineType"`        // 引擎类型
	Labels            []*LabelInfo       `json:"labels"`            // 标签列表
	Vip               string             `json:"vip"`               // 实例 VIP 地址
	ClusterType       int32              `json:"clusterType"`       // 集群类型编码 1-单机版 2-集群版
	Version           string             `json:"version"`           // 版本号
	NodeSize          int32              `json:"nodeSize"`          // broker 节点数量
	OuterProjectName  string             `json:"outerProjectName"`  // 企业项目名称
	VpcId             string             `json:"vpcId"`             // VPC ID
	EnableIpv6        bool               `json:"enableIpv6"`        // 是否启用 IPv6
	Resources         []string           `json:"resources"`         // 资源列表
	ExpandAccess      []string           `json:"expandAccess"`      // 可扩展权限列表
	ShrinkAccess      []string           `json:"shrinkAccess"`      // 可缩容权限列表
	CrtTime           string             `json:"crtTime"`           // 创建时间（UTC 格式）
	ExpTime           string             `json:"expTime"`           // 过期时间。按需返回，可能为 null
}

type ClusterProperties struct {
	EnableTpsLimit    string `json:"enableTpsLimit"`    // 是否启用 TPS 限制。0：未启用，1：已启用
	DeleteWhen        string `json:"deleteWhen"`        // 清理时间配置
	TpsLimit          string `json:"tpsLimit"`          // 集群级 TPS 限制值
	MaxQueuesPerTopic string `json:"maxQueuesPerTopic"` // 每个 Topic 最大队列数
}

type LabelInfo struct {
	LabelId string `json:"labelId"` // 标签唯一标识 ID
	Key     string `json:"key"`     // 标签键
	Value   string `json:"value"`   // 标签值
}
