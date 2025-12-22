package sdwan

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SdwanSdwanOrderEdgeRefundApi
/* 支持退订一个包周期计费/按需的SDWAN智能网关。退订后，将退还对应部分SDWAN智能网关费用。 */
type SdwanSdwanOrderEdgeRefundApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSdwanSdwanOrderEdgeRefundApi(client *core.CtyunClient) *SdwanSdwanOrderEdgeRefundApi {
	return &SdwanSdwanOrderEdgeRefundApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sdwan/refund",
			ContentType:  "application/json",
		},
	}
}

func (a *SdwanSdwanOrderEdgeRefundApi) Do(ctx context.Context, credential core.Credential, req *SdwanSdwanOrderEdgeRefundRequest) (*SdwanSdwanOrderEdgeRefundResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SdwanSdwanOrderEdgeRefundRequest
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
	var resp SdwanSdwanOrderEdgeRefundResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type SdwanSdwanOrderEdgeRefundRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一。参考[订单幂等性说明](#33-订单幂等操作)。  */
	ResourceID  string  `json:"resourceID"`            /*  智能网关资源ID  */
}

type SdwanSdwanOrderEdgeRefundResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *SdwanSdwanOrderEdgeRefundReturnObjResponse `json:"returnObj"`   /*  返回参数  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码.参考[结果码](#61-通用结果码)  */
	Details     *string                                     `json:"details"`     /*  错误明细。一般情况下，会对订单侧(bss)的SDWAN智能网关订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode。  */
	Error       *string                                     `json:"error"`       /*  业务细分码，为product.module.code三段式码.参考[结果码](#61-通用结果码)  */
}

type SdwanSdwanOrderEdgeRefundReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态。参考[结果码](#61-通用结果码)  */
	MasterOrderNO *string `json:"masterOrderNO"` /*  订单号  */
	RegionID      *string `json:"regionID"`      /*  对于SDWAN产品为空。  */
}
