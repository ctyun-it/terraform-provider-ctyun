package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaChaosQueryReqMsgApi
/* 查询故障演练请求详情。
 */type CtgkafkaChaosQueryReqMsgApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaChaosQueryReqMsgApi(client *core.CtyunClient) *CtgkafkaChaosQueryReqMsgApi {
	return &CtgkafkaChaosQueryReqMsgApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/chaos/queryReqMsg",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaChaosQueryReqMsgApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaChaosQueryReqMsgRequest) (*CtgkafkaChaosQueryReqMsgResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("experimentId", req.ExperimentId)
	ctReq.AddParam("taskId", req.TaskId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaChaosQueryReqMsgResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaChaosQueryReqMsgRequest struct {
	RegionId     string `json:"regionId,omitempty"`     /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例ID。  */
	ExperimentId string `json:"experimentId,omitempty"` /*  故障演练ID。  */
	TaskId       string `json:"taskId,omitempty"`       /*  任务ID。  */
}

type CtgkafkaChaosQueryReqMsgResponse struct {
	StatusCode int32                                      `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                     `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaChaosQueryReqMsgReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaChaosQueryReqMsgReturnObjResponse struct {
	NodeKillType int32    `json:"nodeKillType,omitempty"` /*  节点宕机模式。  */
	Nodes        []string `json:"nodes"`                  /*  故障演练节点对象列表。  */
	Duration     int32    `json:"duration,omitempty"`     /*  故障演练持续时间。  */
	CpuPercent   int32    `json:"cpuPercent,omitempty"`   /*  演练CPU负载率。  */
}
