package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbQueryCreatePgelbPriceApi
/* 保障型负载均衡创建询价
 */type CtelbQueryCreatePgelbPriceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbQueryCreatePgelbPriceApi(client *core.CtyunClient) *CtelbQueryCreatePgelbPriceApi {
	return &CtelbQueryCreatePgelbPriceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/query-create-price",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbQueryCreatePgelbPriceApi) Do(ctx context.Context, credential core.Credential, req *CtelbQueryCreatePgelbPriceRequest) (*CtelbQueryCreatePgelbPriceResponse, error) {
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
	var resp CtelbQueryCreatePgelbPriceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbQueryCreatePgelbPriceRequest struct {
	ClientToken      string `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID         string `json:"regionID,omitempty"`         /*  区域ID  */
	ProjectID        string `json:"projectID,omitempty"`        /*  企业项目 ID，默认为'0'  */
	VpcID            string `json:"vpcID,omitempty"`            /*  vpc的ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网的ID  */
	Name             string `json:"name,omitempty"`             /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	EipID            string `json:"eipID,omitempty"`            /*  弹性公网IP的ID。当resourceType=external为必填  */
	SlaName          string `json:"slaName,omitempty"`          /*  lb的规格名称, 支持:elb.s2.small，elb.s3.small，elb.s4.small，elb.s5.small，elb.s2.large，elb.s3.large，elb.s4.large，elb.s5.large  */
	ResourceType     string `json:"resourceType,omitempty"`     /*  资源类型。internal：内网负载均衡，external：公网负载均衡  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  负载均衡的私有IP地址，不指定则自动分配  */
	CycleType        string `json:"cycleType,omitempty"`        /*  订购类型：month（包月） / year（包年）  */
	CycleCount       int32  `json:"cycleCount,omitempty"`       /*  订购时长, 当 cycleType = month, 支持续订 1 - 11 个月; 当 cycleType = year, 支持续订 1 - 3 年  */
}

type CtelbQueryCreatePgelbPriceResponse struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbQueryCreatePgelbPriceReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	Error       string                                       `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbQueryCreatePgelbPriceReturnObjResponse struct {
	TotalPrice     float64                                                      `json:"totalPrice"`     /*  总价格  */
	DiscountPrice  float64                                                      `json:"discountPrice"`  /*  折后价格，云主机相关产品有  */
	FinalPrice     float64                                                      `json:"finalPrice"`     /*  最终价格  */
	SubOrderPrices []*CtelbQueryCreatePgelbPriceReturnObjSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格信息  */
}

type CtelbQueryCreatePgelbPriceReturnObjSubOrderPricesResponse struct {
	ServiceTag      string                                                                      `json:"serviceTag,omitempty"` /*  服务类型  */
	TotalPrice      float64                                                                     `json:"totalPrice"`           /*  子订单总价格  */
	FinalPrice      float64                                                                     `json:"finalPrice"`           /*  最终价格  */
	OrderItemPrices []*CtelbQueryCreatePgelbPriceReturnObjSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"`      /*  item价格信息  */
}

type CtelbQueryCreatePgelbPriceReturnObjSubOrderPricesOrderItemPricesResponse struct {
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型  */
	TotalPrice   string `json:"totalPrice,omitempty"`   /*  总价格  */
	FinalPrice   string `json:"finalPrice,omitempty"`   /*  最终价格  */
}
