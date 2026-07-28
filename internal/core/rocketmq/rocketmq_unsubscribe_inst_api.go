package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// RocketmqUnsubscribeInstApi
/* 退订实例。
 */type RocketmqUnsubscribeInstApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewRocketmqUnsubscribeInstApi(client *core.CtyunClient) *RocketmqUnsubscribeInstApi {
	return &RocketmqUnsubscribeInstApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/unsubscribeInst",
			ContentType:  "application/json",
		},
	}
}

func (a *RocketmqUnsubscribeInstApi) Do(ctx context.Context, credential core.Credential, req *RocketmqUnsubscribeInstRequest) (*RocketmqUnsubscribeInstResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*RocketmqUnsubscribeInstRequest
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
	var resp RocketmqUnsubscribeInstResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type RocketmqUnsubscribeInstRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池 ID。您通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例 ID。  */
}

type RocketmqUnsubscribeInstResponse struct {
	StatusCode int32                                  `json:"statusCode"` // 接口系统层面状态码。成功：800，失败：900
	Message    string                                 `json:"message"`    // 描述状态
	ReturnObj  *RocketmqInstancesUnsubscribeReturnObj `json:"returnObj"`  // 返回对象
	Error      string                                 `json:"error"`      // 错误码，只有非成功才有这个字段，方便快速定位问题
}

type RocketmqInstancesUnsubscribeReturnObj struct {
	Data *UnsubscribeData `json:"data"` // 响应数据对象
}

type UnsubscribeData struct {
	ErrorMessage               string                      `json:"errorMessage"`               // 错误信息描述
	BatchOrderPlacementResults []BatchOrderPlacementResult `json:"batchOrderPlacementResults"` // 批量下单结果列表
}

type BatchOrderPlacementResult struct {
	ErrorMessage      string             `json:"errorMessage"`      // 错误信息描述
	Submitted         bool               `json:"submitted"`         // 是否提交成功
	OrderPlacedEvents []OrderPlacedEvent `json:"orderPlacedEvents"` // 下单事件列表
}

type OrderPlacedEvent struct {
	ErrorMessage string `json:"errorMessage"` // 错误信息描述
	Submitted    bool   `json:"submitted"`    // 是否提交成功
	NewOrderId   string `json:"newOrderId"`   // 订单 id
	NewOrderNo   string `json:"newOrderNo"`   // 订单编号
	TotalPrice   string `json:"totalPrice"`   // 订单价格
}
