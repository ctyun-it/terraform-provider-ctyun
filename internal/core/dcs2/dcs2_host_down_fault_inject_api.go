package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2HostDownFaultInjectApi
/* 主机宕机故障注入
 */type Dcs2HostDownFaultInjectApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2HostDownFaultInjectApi(client *core.CtyunClient) *Dcs2HostDownFaultInjectApi {
	return &Dcs2HostDownFaultInjectApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/inject/hostDownFaultInject",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2HostDownFaultInjectApi) Do(ctx context.Context, credential core.Credential, req *Dcs2HostDownFaultInjectRequest) (*Dcs2HostDownFaultInjectResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2HostDownFaultInjectResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2HostDownFaultInjectRequest struct {
	RegionId   string                                 /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string                                 `json:"prodInstId,omitempty"` /*  实例ID  */
	Nodes      []*Dcs2HostDownFaultInjectNodesRequest `json:"nodes"`                /*  目标故障节点列表  */
	Duration   int32                                  `json:"duration,omitempty"`   /*  故障持续时间(秒), 最小值为60  */
	ActionCode string                                 `json:"actionCode,omitempty"` /*  故障动作类型 <li>node-shutdown：主机宕机  */
}

type Dcs2HostDownFaultInjectNodesRequest struct {
	Ip    string `json:"ip,omitempty"`    /*  管理IP  */
	VpcIp string `json:"vpcIp,omitempty"` /*  业务IP  */
}

type Dcs2HostDownFaultInjectResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                    `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2HostDownFaultInjectReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                                    `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                    `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2HostDownFaultInjectReturnObjResponse struct {
	TaskId       string `json:"taskId,omitempty"`       /*  撤销故障任务id，使用该任务id查询撤销故障执行详情  */
	ExperimentId string `json:"experimentId,omitempty"` /*  演练id，故障注入与故障恢复一起构成完整的一次故障演练。使用该id查询故障执行详情与执行撤销故障  */
}
