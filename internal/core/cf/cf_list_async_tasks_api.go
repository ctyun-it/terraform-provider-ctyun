package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListAsyncTasksApi
/* 获取异步任务列表 */
type CfListAsyncTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListAsyncTasksApi(client *core.CtyunClient) *CfListAsyncTasksApi {
	return &CfListAsyncTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/async-tasks",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListAsyncTasksApi) Do(ctx context.Context, credential core.Credential, req *CfListAsyncTasksRequest) (*CfListAsyncTasksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.Prefix != nil && *req.Prefix != "" {
		ctReq.AddParam("prefix", *req.Prefix)
	}
	ctReq.AddParam("qualifier", req.Qualifier)
	if req.Status != nil && *req.Status != "" {
		ctReq.AddParam("status", *req.Status)
	}
	ctReq.AddParam("startedTimeBegin", strconv.FormatInt(int64(req.StartedTimeBegin), 10))
	ctReq.AddParam("startedTimeEnd", strconv.FormatInt(int64(req.StartedTimeEnd), 10))
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.OrderByTime != nil && *req.OrderByTime != "" {
		ctReq.AddParam("orderByTime", *req.OrderByTime)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListAsyncTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListAsyncTasksRequest struct {
	FunctionName     string  `json:"functionName"`          /*  函数名称，函数必须已存在  */
	RegionId         string  `json:"regionId"`              /*  资源池id，标识不同的地区，如：华东1、西南1  */
	Prefix           *string `json:"prefix,omitempty"`      /*  异步任务 ID 前缀，指定后会返回符合前缀的异步任务列表。  */
	Qualifier        string  `json:"qualifier"`             /*  版本或别名，版本包括1,2,...这样的普通版本，和特殊版本LATEST。  */
	Status           *string `json:"status,omitempty"`      /*  异步任务执行状态。<br>Enqueued：异步消息已入队，等待处理。<br>Succeeded：调用执行成功。<br>Failed：调用执行失败。<br>Running：调用执行中。<br>Stopped：调用执行终止。<br>Stopping：执行停止中。<br>Invalid：执行因函数被删除等原因处于无效状态（任务未被执行）。<br>Expired：为任务配置了最长排队等待的期限。该任务因为超期被丢弃（任务未被执行）。<br>Retrying：异步调用因执行错误重试中。  */
	StartedTimeBegin int32   `json:"startedTimeBegin"`      /*  异步任务启动时间段的开始时间，单位毫秒  */
	StartedTimeEnd   int32   `json:"startedTimeEnd"`        /*  异步任务启动时间段的结束时间，单位毫秒  */
	PageIndex        *int32  `json:"pageIndex,omitempty"`   /*  页码，正整数，默认为1  */
	PageSize         *int32  `json:"pageSize,omitempty"`    /*  页大小，正整数，取值范围[1,100]  */
	OrderByTime      *string `json:"orderByTime,omitempty"` /*  返回异步任务列表的排序方式，默认降序<br>asc 表示升序<br>desc 表示降序  */
}

type CfListAsyncTasksResponse struct {
	StatusCode *int32                             `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                            `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败。  */
	Message    *string                            `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListAsyncTasksReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListAsyncTasksReturnObjResponse struct {
	Data       []*CfListAsyncTasksReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListAsyncTasksReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListAsyncTasksReturnObjDataResponse struct {
	TaskPayload         *string `json:"taskPayload"`         /*  入参内容  */
	Method              *string `json:"method"`              /*  方法，GET/POST等HTTP方法  */
	TaskId              *string `json:"taskId"`              /*  任务ID  */
	Status              *string `json:"status"`              /*  异步任务执行状态。<br>Enqueued：异步消息已入队，等待处理。<br>Succeeded：调用执行成功。<br>Failed：调用执行失败。<br>Running：调用执行中。<br>Stopped：调用执行终止。<br>Stopping：执行停止中。<br>Invalid：执行因函数被删除等原因处于无效状态（任务未被执行）。<br>Expired：为任务配置了最长排队等待的期限。该任务因为超期被丢弃（任务未被执行）。<br>Retrying：异步调用因执行错误重试中。  */
	StartTime           *int32  `json:"startTime"`           /*  异步任务开始时间，单位为毫秒  */
	EndTime             *int32  `json:"endTime"`             /*  异步任务结束时间，单位为毫秒  */
	DurationMs          *int32  `json:"durationMs"`          /*  异步任务的执行时长，单位为毫秒  */
	AlreadyRetriedTimes *int32  `json:"alreadyRetriedTimes"` /*  已重试次数  */
}

type CfListAsyncTasksReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}
