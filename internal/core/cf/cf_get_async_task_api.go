package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetAsyncTaskApi
/* 获取异步任务 */
type CfGetAsyncTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetAsyncTaskApi(client *core.CtyunClient) *CfGetAsyncTaskApi {
	return &CfGetAsyncTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/async-tasks/{taskId}",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetAsyncTaskApi) Do(ctx context.Context, credential core.Credential, req *CfGetAsyncTaskRequest) (*CfGetAsyncTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder = builder.ReplaceUrl("taskId", req.TaskId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("qualifier", req.Qualifier)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetAsyncTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetAsyncTaskRequest struct {
	FunctionName string `json:"functionName"` /*  函数名称，函数必须存在  */
	TaskId       string `json:"taskId"`       /*  异步任务 ID，任务ID必须存在  */
	RegionId     string `json:"regionId"`     /*  资源池id，标识不同的地区，如：华东1、西南1  */
	Qualifier    string `json:"qualifier"`    /*  版本或别名，版本包括1,2,...这样的普通版本，和特殊版本LATEST  */
}

type CfGetAsyncTaskResponse struct {
	StatusCode *int32                           `json:"statusCode"` /*  状态码。0表示成功，非0表示不成功  */
	Error      *string                          `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                          `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfGetAsyncTaskReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfGetAsyncTaskReturnObjResponse struct {
	StartTime           *int32  `json:"startTime"`           /*  异步任务开始时间，单位为毫秒  */
	EndTime             *int32  `json:"endTime"`             /*  异步任务结束时间，单位为毫秒  */
	DurationMs          *int32  `json:"durationMs"`          /*  异步任务的执行时长，单位为毫秒  */
	AlreadyRetriedTimes *int32  `json:"alreadyRetriedTimes"` /*  已重试次数  */
	Method              *string `json:"method"`              /*  方法，GET/POST等HTTP方法  */
	TaskId              *string `json:"taskId"`              /*  任务 ID  */
	TaskPayload         *string `json:"taskPayload"`         /*  任务入参  */
	Status              *string `json:"status"`              /*  任务状态。<br>Enqueued：异步消息已入队，等待处理。  <br>Succeeded：调用执行成功。  <br>Failed：调用执行失败。<br>Running：调用执行中。  <br>Stopped：调用执行终止。  <br>Stopping：执行停止中。  <br>Invalid：执行因函数被删除等原因处于无效状态（任务未被执行）。<br>Expired：为任务配置了最长排队等待的期限。该任务因为超期被丢弃（任务未被执行）。  <br>Retrying：异步调用因执行错误重试中。  */
}
