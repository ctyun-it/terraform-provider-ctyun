package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamResetPwdByRootCodeApi
/* 重置用户密码 */
type CtiamResetPwdByRootCodeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamResetPwdByRootCodeApi(client *core.CtyunClient) *CtiamResetPwdByRootCodeApi {
	return &CtiamResetPwdByRootCodeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/resetPwdByRootCode",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamResetPwdByRootCodeApi) Do(ctx context.Context, credential core.Credential, req *CtiamResetPwdByRootCodeRequest) (*CtiamResetPwdByRootCodeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamResetPwdByRootCodeRequest
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
	var resp CtiamResetPwdByRootCodeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamResetPwdByRootCodeRequest struct {
	UserId   string `json:"userId"`   /*  用户id  */
	Password string `json:"password"` /*  加密后密码  */
	Code     string `json:"code"`     /*  验证码  */
}

type CtiamResetPwdByRootCodeResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
