package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqInstQueryDetailV3Api
/* 查询实例详情 V3
 */type RocketmqInstQueryDetailV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqInstQueryDetailV3Api(client *core.CtyunClient) *RocketmqInstQueryDetailV3Api {
	return &RocketmqInstQueryDetailV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/instance/baseInfo",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *RocketmqInstQueryDetailV3Api) Do(ctx context.Context, credential core.Credential, req *RocketmqInstQueryDetailV3Request) (*RocketmqInstQueryDetailV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp RocketmqInstQueryDetailV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqInstQueryDetailV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池 ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例 ID。  */
}

type RocketmqInstQueryDetailV3Response struct {
	StatusCode int32                                  `json:"statusCode"` // 接口系统层面状态码。成功："800"，失败："900"
	Message    string                                 `json:"message"`    // 接口调用状态描述。成功时为"success"，失败时为具体失败信息
	ReturnObj  *RocketmqInstancesQueryDetailReturnObj `json:"returnObj"`  // 核心返回对象
	Error      string                                 `json:"error"`      // 错误码。仅失败时返回，描述具体错误信息
}

type RocketmqInstancesQueryDetailReturnObj struct {
	InstanceBaseInfo *InstanceBaseInfo `json:"instanceBaseInfo"` // MQ 实例基础信息详情。仅接口调用成功时返回
}

type InstanceBaseInfo struct {
	ProdInstId    string `json:"prodInstId"`    // 实例 id
	BillMode      string `json:"billMode"`      // 计费模式 1：包周期；2：按需
	ProdInstName  string `json:"prodInstName"`  // MQ 实例名称
	RunningState  string `json:"runningState"`  // 实例运行状态
	ProdType      string `json:"prodType"`      // 产品类型
	MachineSpec   string `json:"machineSpec"`   // 机器规格
	TpsLimit      string `json:"tpsLimit"`      // TPS 限制值
	TopicLimit    string `json:"topicLimit"`    // Topic 数量限制。按需返回，可能为 null
	DiskSpace     string `json:"diskSpace"`     // 磁盘空间大小
	NetName       string `json:"netName"`       // 网络名称（VPC 名称）
	Subnet        string `json:"subnet"`        // 子网名称
	SecurityGroup string `json:"securityGroup"` // 安全组 ID
	ExpTime       string `json:"expTime"`       // 过期时间。按需返回，可能为 null
	NameServer    string `json:"nameServer"`    // NameServer 地址列表，多个地址以分号分隔
	CrtTime       string `json:"crtTime"`       // 创建时间（UTC 格式）
}
