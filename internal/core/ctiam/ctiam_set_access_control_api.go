package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamSetAccessControlApi
type CtiamSetAccessControlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamSetAccessControlApi(client *core.CtyunClient) *CtiamSetAccessControlApi {
	return &CtiamSetAccessControlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/setAccessControl",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamSetAccessControlApi) Do(ctx context.Context, credential core.Credential, req *CtiamSetAccessControlRequest) (*CtiamSetAccessControlResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamSetAccessControlRequest
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
	var resp CtiamSetAccessControlResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamSetAccessControlRequest struct {
	UserId        string `json:"userId"`        /*  用户ID  */
	ConsoleAccess string `json:"consoleAccess"` /*  控制台访问(PUB_100_02_0001: 禁止;PUB_100_02_0002: 关闭禁止)  */
}

type CtiamSetAccessControlResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
