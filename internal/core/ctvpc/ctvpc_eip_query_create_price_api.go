package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcEipQueryCreatePriceApi
/* 创建询价。
 */type CtvpcEipQueryCreatePriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcEipQueryCreatePriceApi(client *core.CtyunClient) *CtvpcEipQueryCreatePriceApi {
	return &CtvpcEipQueryCreatePriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/query-create-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcEipQueryCreatePriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcEipQueryCreatePriceRequest) (*CtvpcEipQueryCreatePriceResponse, error) {
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
	var resp CtvpcEipQueryCreatePriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcEipQueryCreatePriceRequest struct {
	ClientToken       string  `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string  `json:"regionID,omitempty"`          /*  资源池 ID  */
	ProjectID         *string `json:"projectID,omitempty"`         /*  不填默认为默认企业项目，如果需要指定企业项目，则需要填写  */
	CycleType         string  `json:"cycleType,omitempty"`         /*  订购类型：month（包月） / year（包年） / on_demand（按需）  */
	CycleCount        int32   `json:"cycleCount"`                  /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年, 当 cycleType = on_demand 时，可以不传  */
	Name              string  `json:"name,omitempty"`              /*  弹性 IP 名称  */
	Bandwidth         int32   `json:"bandwidth"`                   /*  弹性 IP 的带宽峰值，默认为 1 Mbps,1-300  */
	BandwidthID       *string `json:"bandwidthID,omitempty"`       /*  当 cycleType 为 on_demand 时，可以使用 bandwidthID，将弹性 IP 加入到共享带宽中  */
	DemandBillingType *string `json:"demandBillingType,omitempty"` /*  按需计费类型，当 cycleType 为 on_demand 时必填，支持 bandwidth（按带宽）/ upflowc（按流量）  */
}

type CtvpcEipQueryCreatePriceResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcEipQueryCreatePriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                      `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcEipQueryCreatePriceReturnObjResponse struct {
	TotalPrice     float64                                                    `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                    `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                    `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcEipQueryCreatePriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcEipQueryCreatePriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                   `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                   `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                   `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcEipQueryCreatePriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcEipQueryCreatePriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
