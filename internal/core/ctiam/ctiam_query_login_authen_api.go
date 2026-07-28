package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryLoginAuthenApi
type CtiamQueryLoginAuthenApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryLoginAuthenApi(client *core.CtyunClient) *CtiamQueryLoginAuthenApi {
	return &CtiamQueryLoginAuthenApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/queryLoginAuthen",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryLoginAuthenApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryLoginAuthenRequest) (*CtiamQueryLoginAuthenResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryLoginAuthenRequest
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
	var resp CtiamQueryLoginAuthenResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryLoginAuthenRequest struct {
	UserId string `json:"userId"` /*  用户ID  */
}

type CtiamQueryLoginAuthenResponse struct {
	StatusCode *string                                 `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                 `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                 `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryLoginAuthenReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQueryLoginAuthenReturnObjResponse struct {
	AuthenCode *string `json:"authenCode"` /*  是否重置密码 （10001：下次登录必须重置密码，其他：无需重置）  */
}
