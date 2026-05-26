package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfUpdateWorkflowApi
/* 更新云工作流 */
type CfUpdateWorkflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfUpdateWorkflowApi(client *core.CtyunClient) *CfUpdateWorkflowApi {
	return &CfUpdateWorkflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/openapi/v1/workflow",
			ContentType:  "application/json",
		},
	}
}

func (a *CfUpdateWorkflowApi) Do(ctx context.Context, credential core.Credential, req *CfUpdateWorkflowRequest) (*CfUpdateWorkflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.WorkflowId != nil && *req.WorkflowId != "" {
		ctReq.AddParam("workflowId", *req.WorkflowId)
	}
	_, err := ctReq.WriteJson(struct {
		*CfUpdateWorkflowRequest
		RegionId   interface{} `json:"regionId,omitempty"`
		WorkflowId interface{} `json:"workflowId,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfUpdateWorkflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfUpdateWorkflowRequest struct {
	RegionId         string  `json:"regionId"`              /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	WorkflowId       *string `json:"workflowId,omitempty"`  /*  工作流id标识  */
	Description      *string `json:"description,omitempty"` /*  工作流描述  */
	Definition       string  `json:"definition"`            /*  工作流定义，参考<a href="https://www.ctyun.cn/document/10006234/11056484">流程定义</a>  */
	Role             *string `json:"role,omitempty"`        /*  <a href="https://www.ctyun.cn/document/10006234/11056614">执行角色</a>CTRN，可在<a href="https://iam.ctyun.cn/entrustList">IAM委托页面</a>查看委托详情获取  */
	ExecutionTimeout int32   `json:"executionTimeout"`      /*  执行最长时间(s)，必须大于0 <br> - 快速工作流默认值300，最大值300 <br> - 标准工作流默认值31536000，最大值31536000（365天）  */
}

type CfUpdateWorkflowResponse struct {
	StatusCode *int32                             `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string                            `json:"error"`      /*  错误码  */
	Message    *string                            `json:"message"`    /*  结果描述  */
	ReturnObj  *CfUpdateWorkflowReturnObjResponse `json:"returnObj"`  /*  结果数据  */
}

type CfUpdateWorkflowReturnObjResponse struct {
	WorkflowId  *string `json:"workflowId"`  /*  工作流标识  */
	UpdatedTime *int32  `json:"updatedTime"` /*  更新时间  */
}
