package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetDelegateRoleApi
/* 根据id查询委托角色详情 */
type CtiamGetDelegateRoleApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetDelegateRoleApi(client *core.CtyunClient) *CtiamGetDelegateRoleApi {
	return &CtiamGetDelegateRoleApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/delegate/getDelegateRole",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetDelegateRoleApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetDelegateRoleRequest) (*CtiamGetDelegateRoleResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("id", req.Id)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetDelegateRoleResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetDelegateRoleRequest struct {
	Id string `json:"id"` /*  委托id  */
}

type CtiamGetDelegateRoleResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
