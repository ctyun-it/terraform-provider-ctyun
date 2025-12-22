package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamTotpEffectiveChangeApi
/* MFA绑定或解绑 */
type CtiamTotpEffectiveChangeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamTotpEffectiveChangeApi(client *core.CtyunClient) *CtiamTotpEffectiveChangeApi {
	return &CtiamTotpEffectiveChangeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/user/totpEffectiveChange",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamTotpEffectiveChangeApi) Do(ctx context.Context, credential core.Credential, req *CtiamTotpEffectiveChangeRequest) (*CtiamTotpEffectiveChangeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamTotpEffectiveChangeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamTotpEffectiveChangeRequest struct{}

type CtiamTotpEffectiveChangeResponse struct{}
