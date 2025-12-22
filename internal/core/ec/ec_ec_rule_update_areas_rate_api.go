package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcRuleUpdateAreasRateApi
/* 跨区域互通带宽升降配 */
type EcEcRuleUpdateAreasRateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcRuleUpdateAreasRateApi(client *core.CtyunClient) *EcEcRuleUpdateAreasRateApi {
	return &EcEcRuleUpdateAreasRateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cross-region/update",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcRuleUpdateAreasRateApi) Do(ctx context.Context, credential core.Credential, req *EcEcRuleUpdateAreasRateRequest) (*EcEcRuleUpdateAreasRateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EcEcRuleUpdateAreasRateResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcRuleUpdateAreasRateRequest struct{}

type EcEcRuleUpdateAreasRateResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                                   `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                                   `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcRuleUpdateAreasRateReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcRuleUpdateAreasRateReturnObjResponse struct{}
