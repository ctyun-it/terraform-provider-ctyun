package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamPauseAkApi
/* 禁用ak */
type CtiamPauseAkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamPauseAkApi(client *core.CtyunClient) *CtiamPauseAkApi {
	return &CtiamPauseAkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/credential/pauseAk",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamPauseAkApi) Do(ctx context.Context, credential core.Credential, req *CtiamPauseAkRequest) (*CtiamPauseAkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamPauseAkRequest
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
	var resp CtiamPauseAkResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamPauseAkRequest struct {
	UserId string `json:"userId"` /*  用户id  */
	Ak     string `json:"ak"`     /*  ak  */
}

type CtiamPauseAkResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息  */
	Error      *string `json:"error"`      /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
