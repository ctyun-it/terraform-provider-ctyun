package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryRenewOrderPriceV41Api
/* 支持云主机、云硬盘、弹性公网IP、NAT网关、共享带宽、物理机、性能保障型负载均衡、云主机备份存储库和云硬盘备份存储库产品的包年/包月订单的续订询价功能
 */type CtecsQueryRenewOrderPriceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryRenewOrderPriceV41Api(client *core.CtyunClient) *CtecsQueryRenewOrderPriceV41Api {
	return &CtecsQueryRenewOrderPriceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/renew-query-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryRenewOrderPriceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryRenewOrderPriceV41Request) (*CtecsQueryRenewOrderPriceV41Response, error) {
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
	var resp CtecsQueryRenewOrderPriceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryRenewOrderPriceV41Request struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型  */
	ResourceUUID string `json:"resourceUUID,omitempty"` /*  资源uuid  */
	CycleType    string `json:"cycleType,omitempty"`    /*  订购周期类型，可选值：MONTH 月，YEAR 年  */
	CycleCount   int32  `json:"cycleCount,omitempty"`   /*  订购周期大小，订购周期类型为MONTH时范围[1,36]，订购周期类型为YEAR时范围[1,3]  */
}

type CtecsQueryRenewOrderPriceV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  具体错误码标志  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误信息  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述  */
	ReturnObj   *CtecsQueryRenewOrderPriceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据，参见returnObj对象结构  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryRenewOrderPriceV41ReturnObjResponse struct {
	TotalPrice     float32                                                        `json:"totalPrice"`     /*  总价格，单位CNY  */
	FinalPrice     float32                                                        `json:"finalPrice"`     /*  最终价格，单位CNY  */
	SubOrderPrices []*CtecsQueryRenewOrderPriceV41ReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtecsQueryRenewOrderPriceV41ReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                        `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float32                                                                       `json:"totalPrice"`           /*  总价格，单位CNY  */
	FinalPrice      float32                                                                       `json:"finalPrice"`           /*  最终价格，单位CNY  */
	OrderItemPrices []*CtecsQueryRenewOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  资源价格信息  */
}

type CtecsQueryRenewOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType string  `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   float32 `json:"totalPrice"`             /*  总价格，单位CNY  */
	FinalPrice   float32 `json:"finalPrice"`             /*  最终价格，单位CNY  */
}
