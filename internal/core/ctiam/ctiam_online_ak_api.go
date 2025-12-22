package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamOnlineAkApi
/* 启用ak */
type CtiamOnlineAkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamOnlineAkApi(client *core.CtyunClient) *CtiamOnlineAkApi {
	return &CtiamOnlineAkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/credential/onlineAk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamOnlineAkApi) Do(ctx context.Context, credential core.Credential, req *CtiamOnlineAkRequest) (*CtiamOnlineAkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamOnlineAkRequest
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
	var resp CtiamOnlineAkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamOnlineAkRequest struct {
	UserId string `json:"userId"` /*  用户id  */
	Ak     string `json:"ak"`     /*  ak  */
}

type CtiamOnlineAkResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
