package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryPriceForRenewApi
/* 续订查价 */
type Mq2QueryPriceForRenewApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryPriceForRenewApi(client *core.CtyunClient) *Mq2QueryPriceForRenewApi {
	return &Mq2QueryPriceForRenewApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPriceForRenew",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2QueryPriceForRenewApi) Do(ctx context.Context, credential core.Credential, req *Mq2QueryPriceForRenewRequest) (*Mq2QueryPriceForRenewResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2QueryPriceForRenewRequest
		RegionId interface{} `json:"regionId,omitempty"`
	}{
		req, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2QueryPriceForRenewResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryPriceForRenewRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	CycleCnt   int32  `json:"cycleCnt"`   /*  续费周期，取值范围：值需大于零，不超过384个月  */
}

type Mq2QueryPriceForRenewResponse struct {
	StatusCode *string                                 `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                 `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryPriceForRenewReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Error      *string                                 `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2QueryPriceForRenewReturnObjResponse struct {
	Data []*Mq2QueryPriceForRenewReturnObjDataResponse `json:"data"` /*  价格详情数据列表  */
}

type Mq2QueryPriceForRenewReturnObjDataResponse struct {
	TotalPrice     *string                                                     `json:"totalPrice"`     /*  总价格  */
	SubOrderPrices []*Mq2QueryPriceForRenewReturnObjDataSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格列表  */
	FinalPrice     *string                                                     `json:"finalPrice"`     /*  最终价格  */
}

type Mq2QueryPriceForRenewReturnObjDataSubOrderPricesResponse struct {
	TotalPrice      *string                                                                    `json:"totalPrice"`      /*  子订单总价格  */
	ServiceTag      *string                                                                    `json:"serviceTag"`      /*  服务标签  */
	FinalPrice      *string                                                                    `json:"finalPrice"`      /*  子订单最终价格  */
	OrderItemPrices []*Mq2QueryPriceForRenewReturnObjDataSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  订单项价格列表  */
}

type Mq2QueryPriceForRenewReturnObjDataSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId"`       /*  项目ID  */
	TotalPrice   *string `json:"totalPrice"`   /*  项目总价格  */
	FinalPrice   *string `json:"finalPrice"`   /*  项目最终价格  */
	ResourceType *string `json:"resourceType"` /*  资源类型  */
}
