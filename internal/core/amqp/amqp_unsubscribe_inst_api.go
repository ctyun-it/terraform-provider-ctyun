package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpUnsubscribeInstApi
/* 退订实例。
 */type AmqpUnsubscribeInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpUnsubscribeInstApi(client *core.CtyunClient) *AmqpUnsubscribeInstApi {
	return &AmqpUnsubscribeInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/unsubscribeInst",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpUnsubscribeInstApi) Do(ctx context.Context, credential core.Credential, req *AmqpUnsubscribeInstRequest) (*AmqpUnsubscribeInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpUnsubscribeInstRequest
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
	var resp AmqpUnsubscribeInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpUnsubscribeInstRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type AmqpUnsubscribeInstResponse struct {
	StatusCode string                                `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                                `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpUnsubscribeInstReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpUnsubscribeInstReturnObjResponse struct {
	ErrorMessage               string                                                            `json:"errorMessage"`               /*  错误信息。  */
	BatchOrderPlacementResults []*AmqpUnsubscribeInstReturnObjBatchOrderPlacementResultsResponse `json:"batchOrderPlacementResults"` /*  退订返回信息  */
}

type AmqpUnsubscribeInstReturnObjBatchOrderPlacementResultsResponse struct {
	ErrorMessage      string                                                                             `json:"errorMessage"`      /*  错误信息。  */
	Submitted         *bool                                                                              `json:"submitted"`         /*  是否成功提交。  */
	OrderPlacedEvents []*AmqpUnsubscribeInstReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse `json:"orderPlacedEvents"` /*  退订事件信息。  */
}

type AmqpUnsubscribeInstReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse struct {
	ErrorMessage string  `json:"errorMessage"` /*  错误信息。  */
	Submitted    *bool   `json:"submitted"`    /*  是否成功提交。  */
	NewOrderId   string  `json:"newOrderId"`   /*  订单ID。  */
	NewOrderNo   string  `json:"newOrderNo"`   /*  订单编号。  */
	TotalPrice   float64 `json:"totalPrice"`   /*  价格。  */
}
