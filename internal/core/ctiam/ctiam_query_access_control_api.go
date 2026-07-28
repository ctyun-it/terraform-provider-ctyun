package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamQueryAccessControlApi
type CtiamQueryAccessControlApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamQueryAccessControlApi(client *core.CtyunClient) *CtiamQueryAccessControlApi {
	return &CtiamQueryAccessControlApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/queryAccessControl",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamQueryAccessControlApi) Do(ctx context.Context, credential core.Credential, req *CtiamQueryAccessControlRequest) (*CtiamQueryAccessControlResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamQueryAccessControlRequest
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
	var resp CtiamQueryAccessControlResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamQueryAccessControlRequest struct {
	UserId string `json:"userId"` /*  用户ID  */
}

type CtiamQueryAccessControlResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *string `json:"returnObj"`  /*  返回参数  */
}
