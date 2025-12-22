package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamDeleteDelegateRoleApi
/* 删除委托角色 */
type CtiamDeleteDelegateRoleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamDeleteDelegateRoleApi(client *core.CtyunClient) *CtiamDeleteDelegateRoleApi {
	return &CtiamDeleteDelegateRoleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/delegate/deleteDelegateRole",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamDeleteDelegateRoleApi) Do(ctx context.Context, credential core.Credential, req *CtiamDeleteDelegateRoleRequest) (*CtiamDeleteDelegateRoleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamDeleteDelegateRoleRequest
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
	var resp CtiamDeleteDelegateRoleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamDeleteDelegateRoleRequest struct {
	Id   *int32  `json:"id,omitempty"`   /*  ID （id和委托名称name 必填其中之一）  */
	Name *string `json:"name,omitempty"` /*  委托名称（id和委托名称name 必填其中之一）  */
}

type CtiamDeleteDelegateRoleResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
