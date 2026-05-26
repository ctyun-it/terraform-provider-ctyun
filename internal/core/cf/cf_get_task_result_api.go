package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetTaskResultApi
/* 列出执行任务（task）状态 */
type CfGetTaskResultApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetTaskResultApi(client *core.CtyunClient) *CfGetTaskResultApi {
	return &CfGetTaskResultApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/taskevent",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetTaskResultApi) Do(ctx context.Context, credential core.Credential, req *CfGetTaskResultRequest) (*CfGetTaskResultResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("executionId", req.ExecutionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetTaskResultResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetTaskResultRequest struct {
	RegionId    string `json:"regionId"`    /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	ExecutionId string `json:"executionId"` /*  执行标识  */
}

type CfGetTaskResultResponse struct {
	StatusCode *int32                            `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                           `json:"error"`      /*  错误码  */
	Message    *string                           `json:"message"`    /*  结果描述  */
	ReturnObj  *CfGetTaskResultReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfGetTaskResultReturnObjResponse struct {
	Events []*CfGetTaskResultReturnObjEventsResponse `json:"events"` /*  任务列表  */
}

type CfGetTaskResultReturnObjEventsResponse struct {
	TaskName    *string `json:"taskName"`    /*  任务名称  */
	TaskType    *string `json:"taskType"`    /*  任务类型  */
	EventId     *string `json:"eventId"`     /*  事件标识  */
	EventType   *string `json:"eventType"`   /*  事件类型。取值说明如下： <br> - started <br> - submitted <br> - failed <br> - succeeded  */
	Input       *string `json:"input"`       /*  任务输入  */
	Output      *string `json:"output"`      /*  任务输出  */
	ErrorType   *string `json:"errorType"`   /*  错误类型  */
	ErrorMsg    *string `json:"errorMsg"`    /*  错误消息  */
	StartedTime *int32  `json:"startedTime"` /*  开始时间  */
	StoppedTime *int32  `json:"stoppedTime"` /*  结束时间  */
	Attempt     *int32  `json:"attempt"`     /*  第几次尝试（重试策略下可能大于1）  */
	Parent      *string `json:"parent"`      /*  父任务  */
}
