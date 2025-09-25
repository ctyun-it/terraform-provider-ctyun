package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwQueryRenewPriceApi
/* 续订询价。
 */type CtvpcL2gwQueryRenewPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwQueryRenewPriceApi(client *core.CtyunClient) *CtvpcL2gwQueryRenewPriceApi {
	return &CtvpcL2gwQueryRenewPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw/query-renew-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwQueryRenewPriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwQueryRenewPriceRequest) (*CtvpcL2gwQueryRenewPriceResponse, error) {
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
	var resp CtvpcL2gwQueryRenewPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwQueryRenewPriceRequest struct {
	RegionID   string `json:"regionID,omitempty"`  /*  资源池 ID  */
	CycleType  string `json:"cycleType,omitempty"` /*  订购类型：month（包月） / year（包年）  */
	CycleCount int32  `json:"cycleCount"`          /*  订购时长，包月1~11，包年1~3  */
	L2gwID     string `json:"l2gwID,omitempty"`    /*  l2gwID  */
}

type CtvpcL2gwQueryRenewPriceResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcL2gwQueryRenewPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcL2gwQueryRenewPriceReturnObjResponse struct {
	TotalPrice     float64                                                    `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                    `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                    `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcL2gwQueryRenewPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcL2gwQueryRenewPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                   `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                   `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                   `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcL2gwQueryRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcL2gwQueryRenewPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
