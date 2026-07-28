package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcCgwBillRefundApi
/* 云企业路由器按需订单退订 */
type EcEcCgwBillRefundApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcCgwBillRefundApi(client *core.CtyunClient) *EcEcCgwBillRefundApi {
	return &EcEcCgwBillRefundApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/cgw-bill/refund",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcCgwBillRefundApi) Do(ctx context.Context, credential core.Credential, req *EcEcCgwBillRefundRequest) (*EcEcCgwBillRefundResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcCgwBillRefundRequest
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
	var resp EcEcCgwBillRefundResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcCgwBillRefundRequest struct {
	RegionID    string  `json:"regionID"`              /* 资源池ID */
	ClientToken *string `json:"clientToken,omitempty"` /* 客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一 */
	EcID        string  `json:"ecID"`                  /* 云间高速ID */
	ResourceID  string  `json:"resourceID"`            /* 单项资源的变配、续订、退订等需要该资源项的ID */
}

type EcEcCgwBillRefundResponse struct {
	StatusCode  *int32                              `json:"statusCode"`  /* 返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败 */
	ErrorCode   *string                             `json:"errorCode"`   /* 业务细分码，为product.module.code三段式码 */
	Message     *string                             `json:"message"`     /* 失败时的错误描述，一般为英文描述 */
	Description *string                             `json:"description"` /* 失败时的错误描述，一般为中文描述 */
	TraceID     *string                             `json:"traceID"`     /* 链路追踪ID */
	ReturnObj   *EcEcCgwBillRefundReturnObjResponse `json:"returnObj"`   /* 返回参数 */
}

type EcEcCgwBillRefundReturnObjResponse struct {
	MasterOrderNO *string `json:"masterOrderNO"` /* 订单号 */
	MasterOrderID *string `json:"masterOrderID"` /* 订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态 */
	RegionID      *string `json:"regionID"`      /* 资源所属资源池ID */
}
