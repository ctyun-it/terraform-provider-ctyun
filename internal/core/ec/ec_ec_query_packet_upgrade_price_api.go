package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcQueryPacketUpgradePriceApi
/* 支持云间高速带宽包变配询价，目前只支持升配，即增加带宽 */
type EcEcQueryPacketUpgradePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcQueryPacketUpgradePriceApi(client *core.CtyunClient) *EcEcQueryPacketUpgradePriceApi {
	return &EcEcQueryPacketUpgradePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/query-price-upgrade",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcQueryPacketUpgradePriceApi) Do(ctx context.Context, credential core.Credential, req *EcEcQueryPacketUpgradePriceRequest) (*EcEcQueryPacketUpgradePriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcQueryPacketUpgradePriceRequest
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
	var resp EcEcQueryPacketUpgradePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcQueryPacketUpgradePriceRequest struct {
	EcID       string `json:"ecID"`       /*  云间高速ID  */
	RegionID   string `json:"regionID"`   /*  资源池ID, 例:100054c0416811e9a6690242ac110002  */
	Bandwidth  int32  `json:"bandwidth"`  /*  带宽，单位MB  */
	ResourceID string `json:"resourceID"` /*  云间高速带宽包资源ID  */
}

type EcEcQueryPacketUpgradePriceResponse struct {
	StatusCode  *int32                                        `json:"statusCode"`  /*  返回状态码，取值如下<br/>800:成功<br/>900:失败  */
	Message     *string                                       `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description *string                                       `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *EcEcQueryPacketUpgradePriceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   *string                                       `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
}

type EcEcQueryPacketUpgradePriceReturnObjResponse struct {
	TotalPrice     *float64                                                      `json:"totalPrice"`     /*  所有订单总共价格，单位：元  */
	FinalPrice     *float64                                                      `json:"finalPrice"`     /*  所有订单最终价格，单位：元  */
	IsSucceed      *bool                                                         `json:"isSucceed"`      /*  是否成功  */
	SubOrderPrices []*EcEcQueryPacketUpgradePriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格列表  */
	DiscountPrice  *int32                                                        `json:"discountPrice"`  /*  所有订单折扣价，单位：元  */
}

type EcEcQueryPacketUpgradePriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                      `json:"serviceTag"`      /*  产品标签  */
	TotalPrice      *float64                                                                     `json:"totalPrice"`      /*  订单总共价格，单位：元  */
	FinalPrice      *float64                                                                     `json:"finalPrice"`      /*  订单最终价格，单位：元  */
	OrderItemPrices []*EcEcQueryPacketUpgradePriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  订单资源价格列表  */
	CycleType       *int32                                                                       `json:"cycleType"`       /*  订购周期，按月订购  */
}

type EcEcQueryPacketUpgradePriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string  `json:"itemId"`       /*  订单资源ID  */
	ResourceType *string  `json:"resourceType"` /*  本参数表示订单资源类型。  */
	TotalPrice   *float64 `json:"totalPrice"`   /*  订单资源总共价格，单位：元  */
	FinalPrice   *float64 `json:"finalPrice"`   /*  订单资源最终价格，单位：元  */
}
