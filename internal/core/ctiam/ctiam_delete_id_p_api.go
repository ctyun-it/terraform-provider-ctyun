package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamDeleteIdPApi
/* 删除身份供应商 */
type CtiamDeleteIdPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamDeleteIdPApi(client *core.CtyunClient) *CtiamDeleteIdPApi {
	return &CtiamDeleteIdPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/identityProvider/deleteIdP",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamDeleteIdPApi) Do(ctx context.Context, credential core.Credential, req *CtiamDeleteIdPRequest) (*CtiamDeleteIdPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamDeleteIdPRequest
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
	var resp CtiamDeleteIdPResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamDeleteIdPRequest struct {
	Id int32 `json:"id"` /*  id  */
}

type CtiamDeleteIdPResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
