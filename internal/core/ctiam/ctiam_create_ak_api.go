package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamCreateAkApi
/* 创建密钥 */
type CtiamCreateAkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamCreateAkApi(client *core.CtyunClient) *CtiamCreateAkApi {
	return &CtiamCreateAkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/credential/createAk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamCreateAkApi) Do(ctx context.Context, credential core.Credential, req *CtiamCreateAkRequest) (*CtiamCreateAkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamCreateAkRequest
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
	var resp CtiamCreateAkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamCreateAkRequest struct {
	UserId string `json:"userId"` /*  用户id  */
}

type CtiamCreateAkResponse struct {
	StatusCode *string                         `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string                         `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string                         `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
	ReturnObj  *CtiamCreateAkReturnObjResponse `json:"returnObj"`  /*  返回参数  */
}

type CtiamCreateAkReturnObjResponse struct {
	AppId  *string `json:"appId"`  /*  AK  */
	AppKey *string `json:"appKey"` /*  SK  */
}
