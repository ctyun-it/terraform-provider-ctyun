package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSecurityOpVerifyApi
type CtiamSecurityOpVerifyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSecurityOpVerifyApi(client *core.CtyunClient) *CtiamSecurityOpVerifyApi {
	return &CtiamSecurityOpVerifyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/security/securityOpVerify",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSecurityOpVerifyApi) Do(ctx context.Context, credential core.Credential, req *CtiamSecurityOpVerifyRequest) (*CtiamSecurityOpVerifyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSecurityOpVerifyRequest
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
	var resp CtiamSecurityOpVerifyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSecurityOpVerifyRequest struct {
	OpVerify string `json:"opVerify"` /*  0:关闭敏感操作、1: 开启用户级敏感操作、2:开启账号级敏感操作  */
}

type CtiamSecurityOpVerifyResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
