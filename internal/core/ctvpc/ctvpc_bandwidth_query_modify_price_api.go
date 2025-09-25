package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcBandwidthQueryModifyPriceApi
/* 变配询价。
 */type CtvpcBandwidthQueryModifyPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcBandwidthQueryModifyPriceApi(client *core.CtyunClient) *CtvpcBandwidthQueryModifyPriceApi {
	return &CtvpcBandwidthQueryModifyPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/bandwidth/query-modify-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcBandwidthQueryModifyPriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcBandwidthQueryModifyPriceRequest) (*CtvpcBandwidthQueryModifyPriceResponse, error) {
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
	var resp CtvpcBandwidthQueryModifyPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcBandwidthQueryModifyPriceRequest struct {
	ClientToken string  `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  `json:"regionID,omitempty"`    /*  共享带宽的区域 id。  */
	BandwidthID string  `json:"bandwidthID,omitempty"` /*  共享带宽 id。  */
	ProjectID   *string `json:"projectID,omitempty"`   /*  企业项目 ID，默认为'0'  */
	Bandwidth   int32   `json:"bandwidth"`             /*  共享带宽的带宽峰值。  */
}

type CtvpcBandwidthQueryModifyPriceResponse struct {
	StatusCode  int32                                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcBandwidthQueryModifyPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                            `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcBandwidthQueryModifyPriceReturnObjResponse struct {
	TotalPrice     float64                                                          `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                          `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                          `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcBandwidthQueryModifyPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcBandwidthQueryModifyPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                         `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                         `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                         `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcBandwidthQueryModifyPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcBandwidthQueryModifyPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
