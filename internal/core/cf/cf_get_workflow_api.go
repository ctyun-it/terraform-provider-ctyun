package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfGetWorkflowApi
/* 使用工作流id或工作流名称获取云工作流 */
type CfGetWorkflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfGetWorkflowApi(client *core.CtyunClient) *CfGetWorkflowApi {
	return &CfGetWorkflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/workflow",
			ContentType:  "application/json",
		},
	}
}

func (a *CfGetWorkflowApi) Do(ctx context.Context, credential core.Credential, req *CfGetWorkflowRequest) (*CfGetWorkflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionIdHeader)
	if req.WorkflowId != nil && *req.WorkflowId != "" {
		ctReq.AddParam("workflowId", *req.WorkflowId)
	}
	if req.WorkflowName != nil && *req.WorkflowName != "" {
		ctReq.AddParam("workflowName", *req.WorkflowName)
	}
	if req.RegionIdParam != nil && *req.RegionIdParam != "" {
		ctReq.AddParam("regionId", *req.RegionIdParam)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfGetWorkflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfGetWorkflowRequest struct {
	RegionIdHeader string  `json:"regionId"`               /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	WorkflowId     *string `json:"workflowId,omitempty"`   /*  工作流id标识，与workflowName互斥  */
	WorkflowName   *string `json:"workflowName,omitempty"` /*  工作流名称，与workflowId互斥  */
	RegionIdParam  *string `json:"regionId,omitempty"`     /*  资源池id，通过工作流名称获取时必填  */
}

type CfGetWorkflowResponse struct {
	StatusCode *int32                          `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                         `json:"error"`      /*  错误码  */
	Message    *string                         `json:"message"`    /*  结果描述  */
	ReturnObj  *CfGetWorkflowReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfGetWorkflowReturnObjResponse struct {
	WorkflowId       *string `json:"workflowId"`       /*  工作流标识  */
	Name             *string `json:"name"`             /*  工作流名称  */
	Region           *string `json:"region"`           /*  资源池id  */
	Role             *string `json:"role"`             /*  执行角色  */
	ExecutionMode    *string `json:"executionMode"`    /*  工作流类型  */
	ExecutionTimeout *int32  `json:"executionTimeout"` /*  执行最长时间(s)  */
	Description      *string `json:"description"`      /*  工作流描述  */
	Definition       *string `json:"definition"`       /*  dsl定义  */
	CreatedTime      *int32  `json:"createdTime"`      /*  创建时间  */
	UpdatedTime      *int32  `json:"updatedTime"`      /*  更新时间  */
}
