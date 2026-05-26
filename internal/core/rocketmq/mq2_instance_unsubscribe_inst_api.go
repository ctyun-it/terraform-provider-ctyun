package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2InstanceUnsubscribeInstApi
/* 退订 */
type Mq2InstanceUnsubscribeInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2InstanceUnsubscribeInstApi(client *core.CtyunClient) *Mq2InstanceUnsubscribeInstApi {
	return &Mq2InstanceUnsubscribeInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/unsubscribeInst",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2InstanceUnsubscribeInstApi) Do(ctx context.Context, credential core.Credential, req *Mq2InstanceUnsubscribeInstRequest) (*Mq2InstanceUnsubscribeInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2InstanceUnsubscribeInstRequest
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
	var resp Mq2InstanceUnsubscribeInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2InstanceUnsubscribeInstRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
}

type Mq2InstanceUnsubscribeInstResponse struct {
	StatusCode *string                                      `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                      `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2InstanceUnsubscribeInstReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                                      `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2InstanceUnsubscribeInstReturnObjResponse struct {
	Data *Mq2InstanceUnsubscribeInstReturnObjDataResponse `json:"data"` /*  响应数据对象  */
}

type Mq2InstanceUnsubscribeInstReturnObjDataResponse struct {
	ErrorMessage               *string                                                                      `json:"errorMessage"`               /*  错误信息描述  */
	BatchOrderPlacementResults []*Mq2InstanceUnsubscribeInstReturnObjDataBatchOrderPlacementResultsResponse `json:"batchOrderPlacementResults"` /*  批量下单结果列表  */
}

type Mq2InstanceUnsubscribeInstReturnObjDataBatchOrderPlacementResultsResponse struct {
	ErrorMessage      *string                                                                                       `json:"errorMessage"`      /*  错误信息描述  */
	Submitted         *bool                                                                                         `json:"submitted"`         /*  是否提交成功  */
	OrderPlacedEvents []*Mq2InstanceUnsubscribeInstReturnObjDataBatchOrderPlacementResultsOrderPlacedEventsResponse `json:"orderPlacedEvents"` /*  下单事件列表  */
}

type Mq2InstanceUnsubscribeInstReturnObjDataBatchOrderPlacementResultsOrderPlacedEventsResponse struct {
	ErrorMessage *string  `json:"errorMessage"` /*  错误信息描述  */
	Submitted    *string  `json:"submitted"`    /*  是否提交成功  */
	NewOrderId   *string  `json:"newOrderId"`   /*  订单id  */
	NewOrderNo   *string  `json:"newOrderNo"`   /*  订单编号  */
	TotalPrice   *float64 `json:"totalPrice"`   /*  订单价格  */
}
