package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamRestPasswordApi
type CtiamRestPasswordApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamRestPasswordApi(client *core.CtyunClient) *CtiamRestPasswordApi {
	return &CtiamRestPasswordApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/restPassword",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamRestPasswordApi) Do(ctx context.Context, credential core.Credential, req *CtiamRestPasswordRequest) (*CtiamRestPasswordResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamRestPasswordRequest
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
	var resp CtiamRestPasswordResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamRestPasswordRequest struct {
	UserId      string `json:"userId"`      /*  用户id  */
	OldPassword string `json:"oldPassword"` /*  旧密码  */
	NewPassword string `json:"newPassword"` /*  新密码（密码必须包含大小写，密码长度必须在8-26位之间，密码必须包含数字，密码不能包含指定特殊字符：&,(),=,|,'"<,>;,s，密码必须包含特殊字符）  */
}

type CtiamRestPasswordResponse struct {
	StatusCode *string `json:"statusCode"` /*  状态码  */
	Message    *string `json:"message"`    /*  对于调用失败的补充信息说明  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
