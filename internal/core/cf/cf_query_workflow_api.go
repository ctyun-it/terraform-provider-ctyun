package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfQueryWorkflowApi
/* 分页列出云工作流 */
type CfQueryWorkflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfQueryWorkflowApi(client *core.CtyunClient) *CfQueryWorkflowApi {
	return &CfQueryWorkflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/workflow/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CfQueryWorkflowApi) Do(ctx context.Context, credential core.Credential, req *CfQueryWorkflowRequest) (*CfQueryWorkflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.KeyWord != nil && *req.KeyWord != "" {
		ctReq.AddParam("keyWord", *req.KeyWord)
	}
	if req.PageNum != nil && *req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(*req.PageNum), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfQueryWorkflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfQueryWorkflowRequest struct {
	RegionId string  `json:"regionId"`           /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	KeyWord  *string `json:"keyWord,omitempty"`  /*  关键字  */
	PageNum  *int32  `json:"pageNum,omitempty"`  /*  页码  */
	PageSize *int32  `json:"pageSize,omitempty"` /*  每页大小  */
}

type CfQueryWorkflowResponse struct {
	StatusCode *int32                            `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                           `json:"error"`      /*  错误码  */
	Message    *string                           `json:"message"`    /*  结果描述  */
	ReturnObj  *CfQueryWorkflowReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfQueryWorkflowReturnObjResponse struct {
	PageNum  *int32                                   `json:"pageNum"`  /*  页码  */
	PageSize *int32                                   `json:"pageSize"` /*  每页大小  */
	Total    *int32                                   `json:"total"`    /*  工作流总数量  */
	Flows    []*CfQueryWorkflowReturnObjFlowsResponse `json:"flows"`    /*  工作流列表  */
}

type CfQueryWorkflowReturnObjFlowsResponse struct {
	WorkflowId       *string `json:"workflowId"`       /*  工作流标识  */
	Name             *string `json:"name"`             /*  工作流名称  */
	Region           *string `json:"region"`           /*  资源池id  */
	Role             *string `json:"role"`             /*  执行角色  */
	ExecutionMode    *string `json:"executionMode"`    /*  工作流类型  */
	ExecutionTimeout *int32  `json:"executionTimeout"` /*  执行最长时间(s)  */
	Description      *string `json:"description"`      /*  描述  */
	Definition       *string `json:"definition"`       /*  dsl定义  */
	CreatedTime      *int32  `json:"createdTime"`      /*  创建时间  */
	UpdatedTime      *int32  `json:"updatedTime"`      /*  更新时间  */
}
