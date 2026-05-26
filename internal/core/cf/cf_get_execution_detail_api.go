package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetExecutionDetailApi
/* 获取云工作流执行 */
type CfGetExecutionDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetExecutionDetailApi(client *core.CtyunClient) *CfGetExecutionDetailApi {
	return &CfGetExecutionDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/execution",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetExecutionDetailApi) Do(ctx context.Context, credential core.Credential, req *CfGetExecutionDetailRequest) (*CfGetExecutionDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("executionId", req.ExecutionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetExecutionDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetExecutionDetailRequest struct {
	RegionId    string `json:"regionId"`    /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	ExecutionId string `json:"executionId"` /*  执行id标识  */
}

type CfGetExecutionDetailResponse struct {
	StatusCode *int32                                 `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                                `json:"error"`      /*  错误码  */
	Message    *string                                `json:"message"`    /*  结果描述  */
	ReturnObj  *CfGetExecutionDetailReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfGetExecutionDetailReturnObjResponse struct {
	ExecutionId          *string `json:"executionId"`          /*  执行标识  */
	ExecutionName        *string `json:"executionName"`        /*  执行名称  */
	WorkflowId           *string `json:"workflowId"`           /*  工作流标识  */
	WorkflowName         *string `json:"workflowName"`         /*  工作流名称  */
	Status               *string `json:"status"`               /*  执行状态  */
	Input                *string `json:"input"`                /*  执行输入  */
	Output               *string `json:"output"`               /*  执行输出  */
	ExecutionWorkflowDsl *string `json:"executionWorkflowDsl"` /*  执行dsl  */
	StartedTime          *int32  `json:"startedTime"`          /*  执行开始事件  */
	StoppedTime          *int32  `json:"stoppedTime"`          /*  执行结束时间  */
}
