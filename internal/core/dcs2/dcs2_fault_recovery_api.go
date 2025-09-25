package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2FaultRecoveryApi
/* 故障恢复
 */type Dcs2FaultRecoveryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2FaultRecoveryApi(client *core.CtyunClient) *Dcs2FaultRecoveryApi {
	return &Dcs2FaultRecoveryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/inject/faultRecovery",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2FaultRecoveryApi) Do(ctx context.Context, credential core.Credential, req *Dcs2FaultRecoveryRequest) (*Dcs2FaultRecoveryResponse, error) {
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
	var resp Dcs2FaultRecoveryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2FaultRecoveryRequest struct {
	RegionId     string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例ID  */
	ExperimentId string `json:"experimentId,omitempty"` /*  演练id，故障注入与故障恢复一起构成完整的一次故障演练。使用该id查询故障执行详情与执行撤销故障  */
}

type Dcs2FaultRecoveryResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2FaultRecoveryReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2FaultRecoveryReturnObjResponse struct {
	TaskId       string `json:"taskId,omitempty"`       /*  撤销故障任务id，使用该任务id查询撤销故障执行详情  */
	ExperimentId string `json:"experimentId,omitempty"` /*  演练id，故障注入与故障恢复一起构成完整的一次故障演练。使用该id查询故障执行详情与执行撤销故障  */
}
