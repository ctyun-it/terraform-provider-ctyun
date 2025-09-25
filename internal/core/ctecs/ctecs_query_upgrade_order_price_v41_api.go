package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryUpgradeOrderPriceV41Api
/* 支持云主机、云硬盘、弹性公网IP、NAT网关、共享带宽、性能保障型负载均衡、云主机备份存储库和云硬盘备份存储库产品的包年/包月或按量订单变配时的询价功能
 */type CtecsQueryUpgradeOrderPriceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryUpgradeOrderPriceV41Api(client *core.CtyunClient) *CtecsQueryUpgradeOrderPriceV41Api {
	return &CtecsQueryUpgradeOrderPriceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/order/upgrade-query-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryUpgradeOrderPriceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryUpgradeOrderPriceV41Request) (*CtecsQueryUpgradeOrderPriceV41Response, error) {
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
	var resp CtecsQueryUpgradeOrderPriceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryUpgradeOrderPriceV41Request struct {
	RegionID        string `json:"regionID,omitempty"`        /*  资源池ID  */
	ResourceUUID    string `json:"resourceUUID,omitempty"`    /*  资源uuid  */
	ResourceType    string `json:"resourceType,omitempty"`    /*  资源类型  */
	FlavorName      string `json:"flavorName,omitempty"`      /*  云主机规格，当resourceType为VM时必填  */
	Bandwidth       int32  `json:"bandwidth,omitempty"`       /*  带宽大小，范围[1,2000]，需大于当前带宽，当resourceType为IP时必填  */
	DiskSize        int32  `json:"diskSize,omitempty"`        /*  磁盘大小，范围[10,2000]，需大于当前大小，当resourceType为EBS时必填  */
	NatType         string `json:"natType,omitempty"`         /*  nat规格，当resourceType为NAT时必填  */
	IpPoolBandwidth int32  `json:"ipPoolBandwidth,omitempty"` /*  共享带宽大小，范围[5,2000]，需大于当前带宽，当resourceType为IP_POOL时必填  */
	ElbType         string `json:"elbType,omitempty"`         /*  性能保障型负载均衡类型(支持standardI/standardII/enhancedI/enhancedII/higherI)，当resourceType为PGELB时必填  */
	CbrValue        int32  `json:"cbrValue,omitempty"`        /*  存储库大小，100-1024000GB，当resourceType为CBR_VM或CBR_VBS时必填  */
}

type CtecsQueryUpgradeOrderPriceV41Response struct {
	StatusCode  int32                                            `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                           `json:"errorCode,omitempty"`   /*  具体错误码标志  */
	Message     string                                           `json:"message,omitempty"`     /*  失败时的错误信息  */
	Description string                                           `json:"description,omitempty"` /*  失败时的错误描述  */
	ReturnObj   *CtecsQueryUpgradeOrderPriceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据，参见returnObj对象结构  */
	Error       string                                           `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码。请求成功时不返回该字段  */
}

type CtecsQueryUpgradeOrderPriceV41ReturnObjResponse struct {
	TotalPrice     float32                                                          `json:"totalPrice"`     /*  总价格，单位CNY  */
	DiscountPrice  float32                                                          `json:"discountPrice"`  /*  折后价格，单位CNY  */
	FinalPrice     float32                                                          `json:"finalPrice"`     /*  最终价格，单位CNY  */
	SubOrderPrices []*CtecsQueryUpgradeOrderPriceV41ReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtecsQueryUpgradeOrderPriceV41ReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                          `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float32                                                                         `json:"totalPrice"`           /*  总价格，单位CNY  */
	FinalPrice      float32                                                                         `json:"finalPrice"`           /*  最终价格，单位CNY  */
	OrderItemPrices []*CtecsQueryUpgradeOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  资源价格信息  */
}

type CtecsQueryUpgradeOrderPriceV41ReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType string  `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   float32 `json:"totalPrice"`             /*  总价格，单位CNY  */
	FinalPrice   float32 `json:"finalPrice"`             /*  最终价格，单位CNY  */
}
