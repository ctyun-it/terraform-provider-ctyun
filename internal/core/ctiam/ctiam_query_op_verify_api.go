package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryOpVerifyApi
type CtiamQueryOpVerifyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryOpVerifyApi(client *core.CtyunClient) *CtiamQueryOpVerifyApi {
	return &CtiamQueryOpVerifyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/security/queryOpVerify",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryOpVerifyApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryOpVerifyRequest) (*CtiamQueryOpVerifyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamQueryOpVerifyResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryOpVerifyRequest struct{}

type CtiamQueryOpVerifyResponse struct {
	StatusCode *string                              `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                              `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                              `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamQueryOpVerifyReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamQueryOpVerifyReturnObjResponse struct {
	OpVerify *string `json:"opVerify"` /*  开启敏感操作保护。0:关闭敏感操作、1:开启 用户级敏感操作、2:开启账号级敏感操作  */
}
