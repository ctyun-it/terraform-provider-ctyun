package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamMfaValidationApi
/* 校验虚拟MFA */
type CtiamMfaValidationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamMfaValidationApi(client *core.CtyunClient) *CtiamMfaValidationApi {
	return &CtiamMfaValidationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/user/mfaValidation",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamMfaValidationApi) Do(ctx context.Context, credential core.Credential, req *CtiamMfaValidationRequest) (*CtiamMfaValidationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtiamMfaValidationResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamMfaValidationRequest struct{}

type CtiamMfaValidationResponse struct{}
