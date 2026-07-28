package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamInvalidUserApi
/* 注销用户 */
type CtiamInvalidUserApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamInvalidUserApi(client *core.CtyunClient) *CtiamInvalidUserApi {
	return &CtiamInvalidUserApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/invalidUser",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamInvalidUserApi) Do(ctx context.Context, credential core.Credential, req *CtiamInvalidUserRequest) (*CtiamInvalidUserResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamInvalidUserRequest
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
	var resp CtiamInvalidUserResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamInvalidUserRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamInvalidUserResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
