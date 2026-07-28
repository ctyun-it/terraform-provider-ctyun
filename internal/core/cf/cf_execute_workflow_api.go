package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfExecuteWorkflowApi
/* 异步执行云工作流（仅适用于标准工作流） */
type CfExecuteWorkflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfExecuteWorkflowApi(client *core.CtyunClient) *CfExecuteWorkflowApi {
	return &CfExecuteWorkflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/openapi/v1/execution/start",
			ContentType:  "application/json",
		},
	}
}

func (a *CfExecuteWorkflowApi) Do(ctx context.Context, credential core.Credential, req *CfExecuteWorkflowRequest) (*CfExecuteWorkflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CfExecuteWorkflowRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfExecuteWorkflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfExecuteWorkflowRequest struct {
	RegionId      string `json:"regionId"`      /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	ExecutionName string `json:"executionName"` /*  执行名称  */
	WorkflowId    string `json:"workflowId"`    /*  工作流id  */
	Input         string `json:"input"`         /*  执行输入  */
}

type CfExecuteWorkflowResponse struct {
	StatusCode *int32                              `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                             `json:"error"`      /*  错误码  */
	Message    *string                             `json:"message"`    /*  结果描述  */
	ReturnObj  *CfExecuteWorkflowReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfExecuteWorkflowReturnObjResponse struct {
	ExecutionId    *string `json:"executionId"`    /*  执行id  */
	ExecutionName  *string `json:"executionName"`  /*  执行名称  */
	Status         *string `json:"status"`         /*  执行状态。取值说明如下：<br> - started <br> - completed <br> - faulted <br> - canceled  */
	FlowDefinition *string `json:"flowDefinition"` /*  工作流定义  */
	Output         *string `json:"output"`         /*  执行输出  */
	StartedTime    *int32  `json:"startedTime"`    /*  开始事件  */
	StoppedTime    *int32  `json:"stoppedTime"`    /*  结束时间  */
	ErrorCode      *string `json:"errorCode"`      /*  执行错误码  */
	ErrorMessage   *int32  `json:"errorMessage"`   /*  执行错误消息  */
}
