package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcOrderPacketRenewApi
/* 支持云间高速带宽包包周期计费的续订 */
type EcEcOrderPacketRenewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcOrderPacketRenewApi(client *core.CtyunClient) *EcEcOrderPacketRenewApi {
	return &EcEcOrderPacketRenewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/renew",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcOrderPacketRenewApi) Do(ctx context.Context, credential core.Credential, req *EcEcOrderPacketRenewRequest) (*EcEcOrderPacketRenewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcOrderPacketRenewRequest
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
	var resp EcEcOrderPacketRenewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcOrderPacketRenewRequest struct {
	ClientToken *string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一  */
	EcID        string  `json:"ecID"`                  /*  云间高速ID  */
	RegionID    string  `json:"regionID"`              /*  资源池ID, 例:100054c0416811e9a6690242ac110002  */
	ResourceID  string  `json:"resourceID"`            /*  云间高速带宽包资源ID  */
	CycleType   string  `json:"cycleType"`             /*  包周期类型<br/>取值如下<br/>YEAR:包年 ;<br/>MONTH:包月  */
	CycleCount  int32   `json:"cycleCount"`            /*  包周期数，周期最大长度不能超过36个月  */
}

type EcEcOrderPacketRenewResponse struct {
	StatusCode  *int32                                 `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败  */
	Message     *string                                `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcOrderPacketRenewReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
}

type EcEcOrderPacketRenewReturnObjResponse struct {
	MasterOrderID *string `json:"masterOrderID"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO *string `json:"masterOrderNO"` /*  订单号  */
	RegionID      *string `json:"regionID"`      /*  资源池ID  */
}
