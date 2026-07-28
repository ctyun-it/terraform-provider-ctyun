package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamGetTokenApi
type CtiamGetTokenApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamGetTokenApi(client *core.CtyunClient) *CtiamGetTokenApi {
	return &CtiamGetTokenApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/credential/getToken",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamGetTokenApi) Do(ctx context.Context, credential core.Credential, req *CtiamGetTokenRequest) (*CtiamGetTokenResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamGetTokenResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamGetTokenRequest struct{}

type CtiamGetTokenResponse struct {
	StatusCode *string                         `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                         `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                         `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamGetTokenReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamGetTokenReturnObjResponse struct {
	Token      *string `json:"token"`      /*  token  */
	AccountId  *string `json:"accountId"`  /*  账号ID  */
	UserId     *string `json:"userId"`     /*  用户id  */
	RootUserId *string `json:"rootUserId"` /*  主用户id  */
	ExpireTime *string `json:"expireTime"` /*  过期时间  */
}
