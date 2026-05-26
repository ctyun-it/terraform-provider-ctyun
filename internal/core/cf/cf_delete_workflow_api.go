package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfDeleteWorkflowApi
/* 删除云工作流 */
type CfDeleteWorkflowApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfDeleteWorkflowApi(client *core.CtyunClient) *CfDeleteWorkflowApi {
	return &CfDeleteWorkflowApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodDelete,
			UrlPath:      "/openapi/v1/workflow",
			ContentType:  "application/json",
		},
	}
}

func (a *CfDeleteWorkflowApi) Do(ctx context.Context, credential core.Credential, req *CfDeleteWorkflowRequest) (*CfDeleteWorkflowResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("workflowId", req.WorkflowId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfDeleteWorkflowResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfDeleteWorkflowRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池ID，您可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>接口获取所需资源池ID  */
	WorkflowId string `json:"workflowId"` /*  工作流标识  */
}

type CfDeleteWorkflowResponse struct {
	StatusCode *int32  `json:"statusCode"` /*  状态码，0表示成功，非0表示不成功  */
	Error      *string `json:"error"`      /*  错误码  */
	Message    *string `json:"message"`    /*  结果描述  */
}
