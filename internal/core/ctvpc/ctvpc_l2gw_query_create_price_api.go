package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcL2gwQueryCreatePriceApi
/* 创建询价。
 */type CtvpcL2gwQueryCreatePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcL2gwQueryCreatePriceApi(client *core.CtyunClient) *CtvpcL2gwQueryCreatePriceApi {
	return &CtvpcL2gwQueryCreatePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/l2gw/query-create-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcL2gwQueryCreatePriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcL2gwQueryCreatePriceRequest) (*CtvpcL2gwQueryCreatePriceResponse, error) {
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
	var resp CtvpcL2gwQueryCreatePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcL2gwQueryCreatePriceRequest struct {
	RegionID   string `json:"regionID,omitempty"`  /*  资源池 ID  */
	CycleType  string `json:"cycleType,omitempty"` /*  订购类型：month（包月） / year（包年） / on_demand（按需）  */
	CycleCount int32  `json:"cycleCount"`          /*  订购时长，包年和包月时必填，包月1~11，包年1~3  */
	Spec       string `json:"spec,omitempty"`      /*  规格 STANDARD：标准版  ENHANCED：增强版  */
}

type CtvpcL2gwQueryCreatePriceResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcL2gwQueryCreatePriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
}

type CtvpcL2gwQueryCreatePriceReturnObjResponse struct {
	IsSucceed        *bool                                                       `json:"isSucceed"`        /*  是否调用成功  */
	TotalPrice       float64                                                     `json:"totalPrice"`       /*  总价格  */
	DiscountPrice    float64                                                     `json:"discountPrice"`    /*  折后价格，云主机相关产品有  */
	FinalPrice       float64                                                     `json:"finalPrice"`       /*  最终价格  */
	UsedDiscounts    []*string                                                   `json:"usedDiscounts"`    /*  使用的折扣  */
	SubOrderPrices   []*CtvpcL2gwQueryCreatePriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"`   /*  子订单价格信息  */
	VerifyStatusCode int32                                                       `json:"verifyStatusCode"` /*  返回值状态码  */
}

type CtvpcL2gwQueryCreatePriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                    `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                    `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                    `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcL2gwQueryCreatePriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcL2gwQueryCreatePriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId,omitempty"`       /*  订单项ID  */
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
