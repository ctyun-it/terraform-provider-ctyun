package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreatePolicyApi
/* 创建自定义策略 */
type CtiamCreatePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreatePolicyApi(client *core.CtyunClient) *CtiamCreatePolicyApi {
	return &CtiamCreatePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/policy/createPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreatePolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreatePolicyRequest) (*CtiamCreatePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreatePolicyRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamCreatePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreatePolicyRequest struct {
	PolicyName    string `json:"policyName"`    /*  策略名称  */
	PolicyContent string `json:"policyContent"` /*  策略内容(JSON String)
	{
	"Version": "1.1",
	"Statement": [{
	"Resource": ["ctrn1:ctiam:hab-j:202d18a69ecf4f2b91780f50ae5d973d:11"],
	"Action": ["ecs:cloudServers:put", "ecs1:cloudServerNics:binding", "ecs:ServersGroups:list"],
	"Effect": "allow"
	}]
	}  */
	PolicyRange       int32   `json:"policyRange"`                 /*  策略范围（1：资源池 、2：全局级）  */
	PolicyDescription *string `json:"policyDescription,omitempty"` /*  策略描述  */
}

type CtiamCreatePolicyResponse struct {
	StatusCode *string                             `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                             `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                             `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCreatePolicyReturnObjResponse `json:"returnObj"`  /*  返回对象  */
}

type CtiamCreatePolicyReturnObjResponse struct {
	PolicyId *string `json:"policyId"` /*  策略id  */
}
