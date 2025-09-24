package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUnsubscribeInstV3Api
/* 退订实例。
 */type CtgkafkaUnsubscribeInstV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUnsubscribeInstV3Api(client *core.CtyunClient) *CtgkafkaUnsubscribeInstV3Api {
	return &CtgkafkaUnsubscribeInstV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/unsubscribeInst",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUnsubscribeInstV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUnsubscribeInstV3Request) (*CtgkafkaUnsubscribeInstV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUnsubscribeInstV3Request
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
	var resp CtgkafkaUnsubscribeInstV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUnsubscribeInstV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaUnsubscribeInstV3Response struct {
	StatusCode string                                      `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                      `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaUnsubscribeInstV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUnsubscribeInstV3ReturnObjResponse struct {
	ErrorMessage               string                                                                  `json:"errorMessage,omitempty"`     /*  错误信息。  */
	BatchOrderPlacementResults []*CtgkafkaUnsubscribeInstV3ReturnObjBatchOrderPlacementResultsResponse `json:"batchOrderPlacementResults"` /*  退订返回信息  */
}

type CtgkafkaUnsubscribeInstV3ReturnObjBatchOrderPlacementResultsResponse struct {
	ErrorMessage      string                                                                                   `json:"errorMessage,omitempty"` /*  错误信息。  */
	Submitted         *bool                                                                                    `json:"submitted"`              /*  是否成功提交。  */
	OrderPlacedEvents []*CtgkafkaUnsubscribeInstV3ReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse `json:"orderPlacedEvents"`      /*  退订事件信息。  */
}

type CtgkafkaUnsubscribeInstV3ReturnObjBatchOrderPlacementResultsOrderPlacedEventsResponse struct {
	ErrorMessage string  `json:"errorMessage,omitempty"` /*  错误信息。  */
	Submitted    *bool   `json:"submitted"`              /*  是否成功提交。  */
	NewOrderId   string  `json:"newOrderId,omitempty"`   /*  订单ID。  */
	NewOrderNo   string  `json:"newOrderNo,omitempty"`   /*  订单编号。  */
	TotalPrice   float64 `json:"totalPrice"`             /*  价格。  */
}
