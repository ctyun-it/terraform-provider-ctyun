package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ListV3Api
/* 实例列表 */
type Mq2ListV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ListV3Api(client *core.CtyunClient) *Mq2ListV3Api {
	return &Mq2ListV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *Mq2ListV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ListV3Request) (*Mq2ListV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2ListV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ListV3Request struct {
	RegionId string `json:"regionId"` /*  资源池编码  */
}

type Mq2ListV3Response struct {
	StatusCode *string                     `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                     `json:"message"`    /*  接口调用状态描述。成功时为"success"，失败时为具体失败信息  */
	ReturnObj  *Mq2ListV3ReturnObjResponse `json:"returnObj"`  /*  核心返回对象。成功时包含MQ实例列表数据，失败时为空对象  */
	Error      *string                     `json:"error"`      /*  错误码。仅失败时返回，描述具体错误信息  */
}

type Mq2ListV3ReturnObjResponse struct {
	ProdInstList []*Mq2ListV3ReturnObjProdInstListResponse `json:"prodInstList"` /*  MQ实例列表数据。仅接口调用成功时返回  */
}

type Mq2ListV3ReturnObjProdInstListResponse struct {
	ProdInstId        *string                                                  `json:"prodInstId"`        /*  实例id  */
	BillMode          *string                                                  `json:"billMode"`          /*  计费模式 1：包周期；2：按需  */
	ProdInstName      *string                                                  `json:"prodInstName"`      /*  MQ实例名称  */
	RunningState      *string                                                  `json:"runningState"`      /*  实例运行状态  */
	State             *int32                                                   `json:"state"`             /*  实例状态编码 1.运行中 2.已过期 3.已注销 4.已退订 5.扩容中 6.开通中 7.已取消 8.缩容中 9.重启中 10.网络变更中 11.运维恢复 12.运维停止 13.异常中 15.已欠费 -1.变更 101.开通失败  */
	ProdType          *string                                                  `json:"prodType"`          /*  产品类型  */
	MachineSpec       *string                                                  `json:"machineSpec"`       /*  机器规格名  */
	TopicLimit        *string                                                  `json:"topicLimit"`        /*  Topic数量限制  */
	DiskSpace         *string                                                  `json:"diskSpace"`         /*  磁盘空间大小  */
	NetName           *string                                                  `json:"netName"`           /*  网络名称（VPC名称）  */
	Subnet            *string                                                  `json:"subnet"`            /*  子网名称  */
	SecurityGroup     *string                                                  `json:"securityGroup"`     /*  安全组ID  */
	ModTime           *string                                                  `json:"modTime"`           /*  最后修改时间（UTC格式）  */
	DeployType        *string                                                  `json:"deployType"`        /*  部署类型编码  1-虚机部署  */
	DiskType          *string                                                  `json:"diskType"`          /*  磁盘类型  */
	DiskIsEncrypt     *bool                                                    `json:"diskIsEncrypt"`     /*  磁盘是否加密  */
	FileReservedTime  *int32                                                   `json:"fileReservedTime"`  /*  文件保留时间（单位：小时）  */
	ClusterProperties *Mq2ListV3ReturnObjProdInstListClusterPropertiesResponse `json:"clusterProperties"` /*  集群配置属性  */
	EngineType        *string                                                  `json:"engineType"`        /*  引擎类型  */
	Labels            []*Mq2ListV3ReturnObjProdInstListLabelsResponse          `json:"labels"`            /*  标签列表  */
	Vip               *string                                                  `json:"vip"`               /*  实例VIP地址  */
	ClusterType       *int32                                                   `json:"clusterType"`       /*  集群类型编码 1-单机版 2-集群版  */
	Version           *string                                                  `json:"version"`           /*  版本号  */
	NodeSize          *int32                                                   `json:"nodeSize"`          /*  broker节点数量  */
	OuterProjectName  *string                                                  `json:"outerProjectName"`  /*  企业项目名称  */
	VpcId             *string                                                  `json:"vpcId"`             /*  VPC ID  */
	EnableIpv6        *bool                                                    `json:"enableIpv6"`        /*  是否启用IPv6  */
	Resources         []*string                                                `json:"resources"`         /*  资源列表  */
	ExpandAccess      []*string                                                `json:"expandAccess"`      /*  可扩展权限列表 MQ2_INST_VM_EXPAND-规格扩容；MQ2_INST_DISK_EXPAND-磁盘扩容；MQ2_INST_NODE_EXPAND-节点扩容  */
	ShrinkAccess      []*string                                                `json:"shrinkAccess"`      /*  可缩容权限列表 MQ2_INST_VM_SHRINK-规格缩容  当前只支持规格缩容  */
	CrtTime           *string                                                  `json:"crtTime"`           /*  创建时间（UTC格式）  */
}

type Mq2ListV3ReturnObjProdInstListClusterPropertiesResponse struct {
	EnableTpsLimit    *string `json:"enableTpsLimit"`    /*  是否启用TPS限制。0：未启用，1：已启用  */
	DeleteWhen        *string `json:"deleteWhen"`        /*  清理时间配置  */
	TpsLimit          *string `json:"tpsLimit"`          /*  集群级TPS限制值  */
	MaxQueuesPerTopic *string `json:"maxQueuesPerTopic"` /*  每个Topic最大队列数  */
}

type Mq2ListV3ReturnObjProdInstListLabelsResponse struct {
	LabelId *string `json:"labelId"` /*  标签唯一标识ID  */
	Key     *string `json:"key"`     /*  标签键  */
	Value   *string `json:"value"`   /*  标签值  */
}
