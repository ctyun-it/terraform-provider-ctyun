package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamInvalidDelegateRolePolicyApi
/* 解除角色授权 */
type CtiamInvalidDelegateRolePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamInvalidDelegateRolePolicyApi(client *core.CtyunClient) *CtiamInvalidDelegateRolePolicyApi {
	return &CtiamInvalidDelegateRolePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/invalidDelegateRolePolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamInvalidDelegateRolePolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamInvalidDelegateRolePolicyRequest) (*CtiamInvalidDelegateRolePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamInvalidDelegateRolePolicyRequest
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
	var resp CtiamInvalidDelegateRolePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamInvalidDelegateRolePolicyRequest struct {
	PermissionId string `json:"permissionId"` /*  权限id  */
	AssumeUserId string `json:"assumeUserId"` /*  委托用户id  */
}

type CtiamInvalidDelegateRolePolicyResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
