package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfGetExecutionHistoryApi
/* 根据筛选条件列出执行历史 */
type CfGetExecutionHistoryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetExecutionHistoryApi(client *core.CtyunClient) *CfGetExecutionHistoryApi {
	return &CfGetExecutionHistoryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/execution/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetExecutionHistoryApi) Do(ctx context.Context, credential core.Credential, req *CfGetExecutionHistoryRequest) (*CfGetExecutionHistoryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("workflowId", req.WorkflowId)
	if req.Status != nil && *req.Status != "" {
		ctReq.AddParam("status", *req.Status)
	}
	if req.KeyWord != nil && *req.KeyWord != "" {
		ctReq.AddParam("keyWord", *req.KeyWord)
	}
	ctReq.AddParam("startedTimeBegin", strconv.FormatInt(int64(req.StartedTimeBegin), 10))
	ctReq.AddParam("startedTimeEnd", strconv.FormatInt(int64(req.StartedTimeEnd), 10))
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.ContinuationToken != nil && *req.ContinuationToken != "" {
		ctReq.AddParam("continuationToken", *req.ContinuationToken)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetExecutionHistoryResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetExecutionHistoryRequest struct {
	RegionId          string  `json:"regionId"`                    /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	WorkflowId        string  `json:"workflowId"`                  /*  工作流标识  */
	Status            *string `json:"status,omitempty"`            /*  执行状态，可选值： <br> - started <br> - completed <br> - faulted <br> - canceled  */
	KeyWord           *string `json:"keyWord,omitempty"`           /*  关键字  */
	StartedTimeBegin  int32   `json:"startedTimeBegin"`            /*  时间戳范围起始，按照执行开始时间计算（微秒级Unix时间戳）  */
	StartedTimeEnd    int32   `json:"startedTimeEnd"`              /*  时间戳范围终止，按照执行开始时间计算（微秒级Unix时间戳）  */
	PageSize          *int32  `json:"pageSize,omitempty"`          /*  每页大小  */
	ContinuationToken *string `json:"continuationToken,omitempty"` /*  首次查询留空，后续查询从上一次查询响应中获取  */
}

type CfGetExecutionHistoryResponse struct {
	StatusCode *int32                                  `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                                 `json:"error"`      /*  错误码  */
	Message    *string                                 `json:"message"`    /*  结果描述  */
	ReturnObj  *CfGetExecutionHistoryReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfGetExecutionHistoryReturnObjResponse struct {
	ContinuationToken *string                                                 `json:"continuationToken"` /*  继续查询token，没有下一页时为空  */
	PageSize          *int32                                                  `json:"pageSize"`          /*  每页大小  */
	ExecutionLists    []*CfGetExecutionHistoryReturnObjExecutionListsResponse `json:"executionLists"`    /*  执行列表  */
}

type CfGetExecutionHistoryReturnObjExecutionListsResponse struct {
	ExecutionId   *string `json:"executionId"`   /*  执行标识  */
	ExecutionName *string `json:"executionName"` /*  执行名称  */
	WorkflowId    *string `json:"workflowId"`    /*  工作流标识  */
	WorkflowName  *string `json:"workflowName"`  /*  工作流名称  */
	Status        *string `json:"status"`        /*  执行状态 <br> - started <br> - completed <br> - faulted <br> - canceled  */
	StartedTime   *int32  `json:"startedTime"`   /*  开始时间  */
	StoppedTime   *int32  `json:"stoppedTime"`   /*  结束时间  */
}
