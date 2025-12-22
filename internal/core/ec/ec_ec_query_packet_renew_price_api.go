package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcQueryPacketRenewPriceApi
/* 支持云间高速带宽包包周期计费的续订 */
type EcEcQueryPacketRenewPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcQueryPacketRenewPriceApi(client *core.CtyunClient) *EcEcQueryPacketRenewPriceApi {
	return &EcEcQueryPacketRenewPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/query-price-renew",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcQueryPacketRenewPriceApi) Do(ctx context.Context, credential core.Credential, req *EcEcQueryPacketRenewPriceRequest) (*EcEcQueryPacketRenewPriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcQueryPacketRenewPriceRequest
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
	var resp EcEcQueryPacketRenewPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcQueryPacketRenewPriceRequest struct {
	EcID       string `json:"ecID"`       /*  云间高速ID  */
	RegionID   string `json:"regionID"`   /*  资源池ID, 例:100054c0416811e9a6690242ac110002  */
	ResourceID string `json:"resourceID"` /*  云间高速带宽包资源ID  */
	CycleType  string `json:"cycleType"`  /*  包周期类型，<br/>取值如下：<br/>'YEAR': 包年<br/>'MONTH':包月  */
	CycleCount int32  `json:"cycleCount"` /*  包周期数，周期最大长度不能超过36个月  */
}

type EcEcQueryPacketRenewPriceResponse struct {
	StatusCode  *int32                                      `json:"statusCode"`  /*  返回状态码，取值如下<br/>800:成功<br/>900:失败  */
	Message     *string                                     `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                     `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcQueryPacketRenewPriceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                     `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
}

type EcEcQueryPacketRenewPriceReturnObjResponse struct {
	TotalPrice     *float64                                                  `json:"totalPrice"`     /*  所有订单总共价格，单位：元  */
	FinalPrice     *float64                                                  `json:"finalPrice"`     /*  所有订单最终价格，单位：元  */
	IsSucceed      *bool                                                     `json:"isSucceed"`      /*  是否成功  */
	SubOrderPrices *EcEcQueryPacketRenewPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格列表  */
}

type EcEcQueryPacketRenewPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                  `json:"serviceTag"`      /*  产品标签  */
	TotalPrice      *float64                                                                 `json:"totalPrice"`      /*  订单总共价格，单位：元  */
	FinalPrice      *float64                                                                 `json:"finalPrice"`      /*  订单最终价格，单位：元  */
	OrderItemPrices *EcEcQueryPacketRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  订单资源价格列表  */
}

type EcEcQueryPacketRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string  `json:"itemId"`       /*  订单资源ID  */
	ResourceType *string  `json:"resourceType"` /*  本参数表示订单资源类型。  */
	TotalPrice   *float64 `json:"totalPrice"`   /*  订单资源总共价格，单位：元  */
	FinalPrice   *float64 `json:"finalPrice"`   /*  订单资源最终价格，单位：元  */
}
