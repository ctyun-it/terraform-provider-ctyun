package ec

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EcEcQueryPacketNewPriceApi
/* 支持按需包年/包月 询价云间高速带宽包 */
type EcEcQueryPacketNewPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEcEcQueryPacketNewPriceApi(client *core.CtyunClient) *EcEcQueryPacketNewPriceApi {
	return &EcEcQueryPacketNewPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ec/packet/query-price-new",
			ContentType:  "application/json",
		},
	}
}

func (a *EcEcQueryPacketNewPriceApi) Do(ctx context.Context, credential core.Credential, req *EcEcQueryPacketNewPriceRequest) (*EcEcQueryPacketNewPriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EcEcQueryPacketNewPriceRequest
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
	var resp EcEcQueryPacketNewPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type EcEcQueryPacketNewPriceRequest struct {
	RegionID   string `json:"regionID"`   /*  资源池ID  */
	EcID       string `json:"ecID"`       /*  云间高速ID  */
	Bandwidth  int32  `json:"bandwidth"`  /*  带宽，单位MB  */
	CycleType  string `json:"cycleType"`  /*  包周期类型,<br/>取值如下：<br/>'YEAR': 包年<br/>'MONTH':包月  */
	CycleCount int32  `json:"cycleCount"` /*  包周期数。onDemand为False时必须指定。周期最大长度不能超过36个月  */
	OnDemand   bool   `json:"onDemand"`   /*  布尔类型，是否按需下单。默认为false, 当前不支持按需下单  */
}

type EcEcQueryPacketNewPriceResponse struct {
	StatusCode  *int32                                    `json:"statusCode"`  /*  返回状态码<br/>取值范围:<br/>800:成功<br/>900:失败   */
	Message     *string                                   `json:"message"`     /*   失败时的错误描述，一般为英文描述  */
	Description *string                                   `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ErrorCode   *string                                   `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	ReturnObj   *EcEcQueryPacketNewPriceReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	Error       *string                                   `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type EcEcQueryPacketNewPriceReturnObjResponse struct {
	TotalPrice     *float32                                                  `json:"totalPrice"`     /*  所有订单总共价格，单位：元。  */
	FinalPrice     *float32                                                  `json:"finalPrice"`     /*  所有订单最终价格，单位：元。  */
	IsSucceed      *bool                                                     `json:"isSucceed"`      /*  是否成功  */
	SubOrderPrices []*EcEcQueryPacketNewPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格列表  */
}

type EcEcQueryPacketNewPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                  `json:"serviceTag"`      /*  产品标签  */
	TotalPrice      *float32                                                                 `json:"totalPrice"`      /*  订单总共价格，单位：元。  */
	FinalPrice      *float32                                                                 `json:"finalPrice"`      /*  订单最终价格，单位：元。  */
	OrderItemPrices []*EcEcQueryPacketNewPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  订单资源价格列表  */
}

type EcEcQueryPacketNewPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId     *string `json:"itemId"`     /*  订单资源ID  */
	ServiceTag *string `json:"serviceTag"` /*  本参数表示订单资源类型  */
	TotalPrice *string `json:"totalPrice"` /*  订单资源总共价格，单位：元。  */
	FinalPrice *string `json:"finalPrice"` /*  订单资源最终价格，单位：元。  */
}
