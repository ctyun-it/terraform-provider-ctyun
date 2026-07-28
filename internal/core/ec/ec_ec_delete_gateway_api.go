package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcDeleteGatewayApi
/* 删除云网关 */
type EcEcDeleteGatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcDeleteGatewayApi(client *core.CtyunClient) *EcEcDeleteGatewayApi {
	return &EcEcDeleteGatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cloud-gateway/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcDeleteGatewayApi) Do(ctx context.Context, credential core.Credential, req *EcEcDeleteGatewayRequest) (*EcEcDeleteGatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcDeleteGatewayRequest
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
	var resp EcEcDeleteGatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcDeleteGatewayRequest struct {
	CgwID string `json:"cgwID"` /*  云网关实例ID  */
}

type EcEcDeleteGatewayResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	ErrorCode   *string                             `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     *string                             `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                             `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	TraceID     *string                             `json:"traceID"`     /*  链路追踪ID  */
	ReturnObj   *EcEcDeleteGatewayReturnObjResponse `json:"returnObj"`   /*  返回参数  */
}

type EcEcDeleteGatewayReturnObjResponse struct{}
