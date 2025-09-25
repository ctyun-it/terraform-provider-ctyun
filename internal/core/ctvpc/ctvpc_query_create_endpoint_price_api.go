package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcQueryCreateEndpointPriceApi
/* 创建询价
 */type CtvpcQueryCreateEndpointPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryCreateEndpointPriceApi(client *core.CtyunClient) *CtvpcQueryCreateEndpointPriceApi {
	return &CtvpcQueryCreateEndpointPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/vpce/query-create-endpoint-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryCreateEndpointPriceApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryCreateEndpointPriceRequest) (*CtvpcQueryCreateEndpointPriceResponse, error) {
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
	var resp CtvpcQueryCreateEndpointPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryCreateEndpointPriceRequest struct {
	ClientToken       string    `json:"clientToken,omitempty"`       /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID          string    `json:"regionID,omitempty"`          /*  资源池ID  */
	CycleType         string    `json:"cycleType,omitempty"`         /*  收费类型：只能填写 on_demand  */
	EndpointServiceID string    `json:"endpointServiceID,omitempty"` /*  终端节点关联的终端节点服务  */
	EndpointName      string    `json:"endpointName,omitempty"`      /*  终端节点名称，只能由数字，字母，-组成不能以数字和-开头，最大长度28  */
	SubnetID          string    `json:"subnetID,omitempty"`          /*  子网id  */
	VpcID             string    `json:"vpcID,omitempty"`             /*  vpc-xxxx  */
	IP                *string   `json:"IP,omitempty"`                /*  vpc address  */
	WhitelistFlag     int32     `json:"whitelistFlag"`               /*  白名单开关 1.开启 0.关闭，默认1  */
	Whitelist         []*string `json:"whitelist"`                   /*  白名单  */
	Description       *string   `json:"description,omitempty"`       /*  描述,内容限制：1、长度限制100 2、支持汉字，大小写字母，数字 3、支持英文特殊字符：~!@#$%^&*()_-+=<>?:'{}  */
}

type CtvpcQueryCreateEndpointPriceResponse struct {
	StatusCode  int32                                           `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                         `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                         `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                         `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcQueryCreateEndpointPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       *string                                         `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcQueryCreateEndpointPriceReturnObjResponse struct {
	TotalPrice     float64                                                         `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                         `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                         `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtvpcQueryCreateEndpointPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtvpcQueryCreateEndpointPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      *string                                                                        `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                        `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                        `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtvpcQueryCreateEndpointPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtvpcQueryCreateEndpointPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   *string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   *string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
