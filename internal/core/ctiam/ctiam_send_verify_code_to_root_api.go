package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSendVerifyCodeToRootApi
/* 重置密码发送短信验证码 */
type CtiamSendVerifyCodeToRootApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSendVerifyCodeToRootApi(client *core.CtyunClient) *CtiamSendVerifyCodeToRootApi {
	return &CtiamSendVerifyCodeToRootApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/sendVerifyCodeToRoot",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSendVerifyCodeToRootApi) Do(ctx context.Context, credential core.Credential, req *CtiamSendVerifyCodeToRootRequest) (*CtiamSendVerifyCodeToRootResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSendVerifyCodeToRootRequest
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
	var resp CtiamSendVerifyCodeToRootResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSendVerifyCodeToRootRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamSendVerifyCodeToRootResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
