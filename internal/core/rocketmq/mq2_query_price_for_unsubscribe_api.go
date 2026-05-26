package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2QueryPriceForUnsubscribeApi
/* 退订查价 */
type Mq2QueryPriceForUnsubscribeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2QueryPriceForUnsubscribeApi(client *core.CtyunClient) *Mq2QueryPriceForUnsubscribeApi {
	return &Mq2QueryPriceForUnsubscribeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instance/queryPriceForUnsubscribe",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2QueryPriceForUnsubscribeApi) Do(ctx context.Context, credential core.Credential, req *Mq2QueryPriceForUnsubscribeRequest) (*Mq2QueryPriceForUnsubscribeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2QueryPriceForUnsubscribeRequest
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
	var resp Mq2QueryPriceForUnsubscribeResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2QueryPriceForUnsubscribeRequest struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
}

type Mq2QueryPriceForUnsubscribeResponse struct {
	StatusCode *string                                       `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                                       `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2QueryPriceForUnsubscribeReturnObjResponse `json:"returnObj"`  /*  响应对象  */
	Error      *string                                       `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2QueryPriceForUnsubscribeReturnObjResponse struct {
	PaidPrice   *string `json:"paidPrice"`   /*  支付金额  */
	RefundPrice *string `json:"refundPrice"` /*  退款金额  */
}
