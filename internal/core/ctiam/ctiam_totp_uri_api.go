package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamTotpUriApi
/* 获取虚拟MFA二维码 */
type CtiamTotpUriApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamTotpUriApi(client *core.CtyunClient) *CtiamTotpUriApi {
	return &CtiamTotpUriApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/user/totpUri",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamTotpUriApi) Do(ctx context.Context, credential core.Credential, req *CtiamTotpUriRequest) (*CtiamTotpUriResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamTotpUriResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamTotpUriRequest struct{}

type CtiamTotpUriResponse struct{}
