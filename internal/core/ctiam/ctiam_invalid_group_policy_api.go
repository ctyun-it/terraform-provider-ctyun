package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamInvalidGroupPolicyApi
/* 用户组取消权限 */
type CtiamInvalidGroupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamInvalidGroupPolicyApi(client *core.CtyunClient) *CtiamInvalidGroupPolicyApi {
	return &CtiamInvalidGroupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/invalidGroupPolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamInvalidGroupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamInvalidGroupPolicyRequest) (*CtiamInvalidGroupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamInvalidGroupPolicyRequest
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
	var resp CtiamInvalidGroupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamInvalidGroupPolicyRequest struct {
	PrivilegeId string `json:"privilegeId"` /*  权限id  */
}

type CtiamInvalidGroupPolicyResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
