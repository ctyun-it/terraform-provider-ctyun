package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamAttachDelegateRolePolicyApi
/* 委托角色授权 */
type CtiamAttachDelegateRolePolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamAttachDelegateRolePolicyApi(client *core.CtyunClient) *CtiamAttachDelegateRolePolicyApi {
	return &CtiamAttachDelegateRolePolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/perm/attachDelegateRolePolicy",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamAttachDelegateRolePolicyApi) Do(ctx context.Context, credential core.Credential, req *CtiamAttachDelegateRolePolicyRequest) (*CtiamAttachDelegateRolePolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamAttachDelegateRolePolicyRequest
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
	var resp CtiamAttachDelegateRolePolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamAttachDelegateRolePolicyRequest struct {
	AssumeUserId  string   `json:"assumeUserId"`  /*  委托用户id  */
	RangeType     string   `json:"rangeType"`     /*  授权类型  */
	PolicyIds     []string `json:"policyIds"`     /*  策略id列表  */
	RegionIds     []string `json:"regionIds"`     /*  资源池列表  */
	PrincipalType string   `json:"principalType"` /*  授权实体类型（1：代表用户授权，2：代表委托授权）  */
}

type CtiamAttachDelegateRolePolicyResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
