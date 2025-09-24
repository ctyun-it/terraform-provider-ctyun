package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcIpv6BandwidthQueryRenewPriceApi
/* 续订询价。
 */type CtvpcIpv6BandwidthQueryRenewPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcIpv6BandwidthQueryRenewPriceApi(client *core.CtyunClient) *CtvpcIpv6BandwidthQueryRenewPriceApi {
	return &CtvpcIpv6BandwidthQueryRenewPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ipv6_bandwidth/query-renew-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcIpv6BandwidthQueryRenewPriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcIpv6BandwidthQueryRenewPriceRequest) (*CtvpcIpv6BandwidthQueryRenewPriceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcIpv6BandwidthQueryRenewPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcIpv6BandwidthQueryRenewPriceRequest struct {
	RegionID    string `json:"regionID,omitempty"`    /*  共享带宽的区域id。  */
	BandwidthID string `json:"bandwidthID,omitempty"` /*  共享带宽 ID  */
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订购类型：month / year  */
	CycleCount  int32  `json:"cycleCount"`            /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
}

type CtvpcIpv6BandwidthQueryRenewPriceResponse struct {
	StatusCode  int32                                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcIpv6BandwidthQueryRenewPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcIpv6BandwidthQueryRenewPriceReturnObjResponse struct {
	TotalPrice     float64                                                             `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                             `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                             `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcIpv6BandwidthQueryRenewPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcIpv6BandwidthQueryRenewPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                            `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                            `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                            `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcIpv6BandwidthQueryRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcIpv6BandwidthQueryRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
