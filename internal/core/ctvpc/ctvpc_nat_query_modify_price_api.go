package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcNatQueryModifyPriceApi
/* 变配询价。
 */type CtvpcNatQueryModifyPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNatQueryModifyPriceApi(client *core.CtyunClient) *CtvpcNatQueryModifyPriceApi {
	return &CtvpcNatQueryModifyPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/nat/query-modify-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNatQueryModifyPriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNatQueryModifyPriceRequest) (*CtvpcNatQueryModifyPriceResponse, error) {
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
	var resp CtvpcNatQueryModifyPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNatQueryModifyPriceRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  NAT 网关所在的地域 ID。  */
	NatGatewayID string `json:"natGatewayID,omitempty"` /*  NAT 网关 ID。  */
	Spec         int32  `json:"spec"`                   /*  规格(可传值：1-SMALL,2-MEDIUM,3-LARGE,4-XLARGE)  */
	ClientToken  string `json:"clientToken,omitempty"`  /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
}

type CtvpcNatQueryModifyPriceResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNatQueryModifyPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNatQueryModifyPriceReturnObjResponse struct {
	TotalPrice     float64                                                    `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                    `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                    `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcNatQueryModifyPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcNatQueryModifyPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                   `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                   `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                   `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcNatQueryModifyPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcNatQueryModifyPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
