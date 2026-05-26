package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetTaskEventTraceApi
/* 列出执行任务（task）的所有事件 */
type CfGetTaskEventTraceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetTaskEventTraceApi(client *core.CtyunClient) *CfGetTaskEventTraceApi {
	return &CfGetTaskEventTraceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/taskevent/trace",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetTaskEventTraceApi) Do(ctx context.Context, credential core.Credential, req *CfGetTaskEventTraceRequest) (*CfGetTaskEventTraceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("executionId", req.ExecutionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetTaskEventTraceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetTaskEventTraceRequest struct {
	RegionId    string `json:"regionId"`    /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	ExecutionId string `json:"executionId"` /*  执行标识  */
}

type CfGetTaskEventTraceResponse struct {
	StatusCode *int32                                `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                               `json:"error"`      /*  错误码  */
	Message    *string                               `json:"message"`    /*  结果描述  */
	ReturnObj  *CfGetTaskEventTraceReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfGetTaskEventTraceReturnObjResponse struct {
	Records []*CfGetTaskEventTraceReturnObjRecordsResponse `json:"records"` /*  任务列表  */
}

type CfGetTaskEventTraceReturnObjRecordsResponse struct {
	TaskName  *string `json:"taskName"`  /*  任务名称  */
	TaskType  *string `json:"taskType"`  /*  任务类型  */
	EventType *string `json:"eventType"` /*  事件类型。取值说明如下： <br> - started <br> - submitted <br> - failed <br> - succeeded  */
	EventTime *int32  `json:"eventTime"` /*  事件发生时间  */
	Input     *string `json:"input"`     /*  事件输入  */
	Output    *string `json:"output"`    /*  事件输出  */
}
