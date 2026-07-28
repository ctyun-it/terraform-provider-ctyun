package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSetLoginAuthenApi
type CtiamSetLoginAuthenApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSetLoginAuthenApi(client *core.CtyunClient) *CtiamSetLoginAuthenApi {
	return &CtiamSetLoginAuthenApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/setLoginAuthen",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSetLoginAuthenApi) Do(ctx context.Context, credential core.Credential, req *CtiamSetLoginAuthenRequest) (*CtiamSetLoginAuthenResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSetLoginAuthenRequest
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
	var resp CtiamSetLoginAuthenResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSetLoginAuthenRequest struct {
	UserId string `json:"userId"` /*  用户ID  */
	Enable bool   `json:"enable"` /*  true:用户在下次登录时必须重置密码, false:无需重置  */
}

type CtiamSetLoginAuthenResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
