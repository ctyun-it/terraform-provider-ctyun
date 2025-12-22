package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamTotpEffectiveApi
/* 查询虚拟MFA是否绑定 */
type CtiamTotpEffectiveApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamTotpEffectiveApi(client *core.CtyunClient) *CtiamTotpEffectiveApi {
	return &CtiamTotpEffectiveApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/user/totpEffective",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamTotpEffectiveApi) Do(ctx context.Context, credential core.Credential, req *CtiamTotpEffectiveRequest) (*CtiamTotpEffectiveResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamTotpEffectiveResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamTotpEffectiveRequest struct{}

type CtiamTotpEffectiveResponse struct{}
