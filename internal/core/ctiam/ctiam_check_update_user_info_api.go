package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCheckUpdateUserInfoApi
type CtiamCheckUpdateUserInfoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCheckUpdateUserInfoApi(client *core.CtyunClient) *CtiamCheckUpdateUserInfoApi {
	return &CtiamCheckUpdateUserInfoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/checkUpdateUserInfo",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCheckUpdateUserInfoApi) Do(ctx context.Context, credential core.Credential, req *CtiamCheckUpdateUserInfoRequest) (*CtiamCheckUpdateUserInfoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCheckUpdateUserInfoRequest
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
	var resp CtiamCheckUpdateUserInfoResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCheckUpdateUserInfoRequest struct {
	UserId      string  `json:"userId"`                /*  用户id  */
	LoginEmail  *string `json:"loginEmail,omitempty"`  /*  邮箱（邮箱与手机号必填其中之一）  */
	MobilePhone *string `json:"mobilePhone,omitempty"` /*  手机号（邮箱与手机号必填其中之一）  */
}

type CtiamCheckUpdateUserInfoResponse struct {
	StatusCode *string                                    `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                    `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                    `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCheckUpdateUserInfoReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamCheckUpdateUserInfoReturnObjResponse struct {
	CheckLoginEmailRes  *bool `json:"checkLoginEmailRes"`  /*  邮箱是否可以使用：true：是，false：否  */
	CheckMobilePhoneRes *bool `json:"checkMobilePhoneRes"` /*  手机号是否可以使用：true：是，false：否  */
}
