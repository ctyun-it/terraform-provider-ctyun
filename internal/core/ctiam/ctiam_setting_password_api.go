package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSettingPasswordApi
type CtiamSettingPasswordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSettingPasswordApi(client *core.CtyunClient) *CtiamSettingPasswordApi {
	return &CtiamSettingPasswordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/settingPassword",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSettingPasswordApi) Do(ctx context.Context, credential core.Credential, req *CtiamSettingPasswordRequest) (*CtiamSettingPasswordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSettingPasswordRequest
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
	var resp CtiamSettingPasswordResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSettingPasswordRequest struct {
	UserId   string  `json:"userId"`             /*  用户id  */
	RawType  string  `json:"type"`               /*  类型：1:保留当前密码, 2:自动生成密码, 3:自定义密码  */
	Password *string `json:"password,omitempty"` /*  自定义密码(类型为3时必填)  */
}

type CtiamSettingPasswordResponse struct {
	StatusCode *string                                `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                                `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                                `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamSettingPasswordReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamSettingPasswordReturnObjResponse struct {
	Password *string `json:"password"` /*  密码  */
}
