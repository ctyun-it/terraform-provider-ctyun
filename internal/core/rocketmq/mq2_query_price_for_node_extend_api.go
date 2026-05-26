package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryPriceForNodeExtendApi
/* 节点扩容查价 */
type Mq2QueryPriceForNodeExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryPriceForNodeExtendApi(client *core.CtyunClient) *Mq2QueryPriceForNodeExtendApi {
	return &Mq2QueryPriceForNodeExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPriceForNodeExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2QueryPriceForNodeExtendApi) Do(ctx context.Context, credential core.Credential, req *Mq2QueryPriceForNodeExtendRequest) (*Mq2QueryPriceForNodeExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2QueryPriceForNodeExtendRequest
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
	var resp Mq2QueryPriceForNodeExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryPriceForNodeExtendRequest struct {
	RegionId      string `json:"regionId"`      /*  资源池编码  */
	ProdInstId    string `json:"prodInstId"`    /*  实例ID  */
	ExtendNodeNum int32  `json:"extendNodeNum"` /*  扩容后的节点数  */
}

type Mq2QueryPriceForNodeExtendResponse struct {
	StatusCode *string                                      `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                      `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryPriceForNodeExtendReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                                      `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2QueryPriceForNodeExtendReturnObjResponse struct {
	Data *Mq2QueryPriceForNodeExtendReturnObjDataResponse `json:"data"` /*  价格详情数据  */
}

type Mq2QueryPriceForNodeExtendReturnObjDataResponse struct {
	TotalPrice     *string                                                          `json:"totalPrice"`     /*  总价格  */
	SubOrderPrices []*Mq2QueryPriceForNodeExtendReturnObjDataSubOrderPricesResponse `json:"subOrderPrices"` /*  子订单价格列表  */
	FinalPrice     *string                                                          `json:"finalPrice"`     /*  最终价格  */
}

type Mq2QueryPriceForNodeExtendReturnObjDataSubOrderPricesResponse struct {
	TotalPrice      *string                                                                         `json:"totalPrice"`      /*  子订单总价格  */
	ServiceTag      *string                                                                         `json:"serviceTag"`      /*  服务标签  */
	FinalPrice      *string                                                                         `json:"finalPrice"`      /*  子订单最终价格  */
	OrderItemPrices []*Mq2QueryPriceForNodeExtendReturnObjDataSubOrderPricesOrderItemPricesResponse `json:"orderItemPrices"` /*  订单项价格列表  */
}

type Mq2QueryPriceForNodeExtendReturnObjDataSubOrderPricesOrderItemPricesResponse struct {
	ItemId       *string `json:"itemId"`       /*  项目id  */
	TotalPrice   *string `json:"totalPrice"`   /*  项目总价格  */
	FinalPrice   *string `json:"finalPrice"`   /*  项目最终价格  */
	ResourceType *string `json:"resourceType"` /*  资源类型（PAAS_DMQ为rocketmq实例，PAAS_DMQ_EBS为云硬盘）  */
}
