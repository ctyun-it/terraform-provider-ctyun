package ctiam

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtiamInvalidEpGroupPloyDecreaseApi
type CtiamInvalidEpGroupPloyDecreaseApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtiamInvalidEpGroupPloyDecreaseApi(client *core.CtyunClient) *CtiamInvalidEpGroupPloyDecreaseApi {
	return &CtiamInvalidEpGroupPloyDecreaseApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v1/project/invalidEpGroupPloyDecrease",
			ContentType:  "application/json",
		},
	}
}

func (a *CtiamInvalidEpGroupPloyDecreaseApi) Do(ctx context.Context, credential core.Credential, req *CtiamInvalidEpGroupPloyDecreaseRequest) (*CtiamInvalidEpGroupPloyDecreaseResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*CtiamInvalidEpGroupPloyDecreaseRequest
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
	var resp CtiamInvalidEpGroupPloyDecreaseResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtiamInvalidEpGroupPloyDecreaseRequest struct {
	PrivilegeId string `json:"privilegeId"` /*  授权id  */
}

type CtiamInvalidEpGroupPloyDecreaseResponse struct {
	StatusCode *string `json:"statusCode"` /*  兼容性返回码，800代表成功，CTIAM_XXX 为失败码  */
	Message    *string `json:"message"`    /*  返回信息，请求失败时会回传错误信息
	 */
	Error *string `json:"error"` /*  异常码，请求失败时会返回(CTIAM_XXXX)  */
}
