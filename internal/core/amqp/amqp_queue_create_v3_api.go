package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueueCreateV3Api
/* 创建队列v3
 */type AmqpQueueCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueueCreateV3Api(client *core.CtyunClient) *AmqpQueueCreateV3Api {
	return &AmqpQueueCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/queue/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpQueueCreateV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpQueueCreateV3Request) (*AmqpQueueCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpQueueCreateV3Request
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
	var resp AmqpQueueCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueueCreateV3Request struct {
	RegionId              string  `json:"regionId,omitempty"`                  /*  资源池id  */
	ProdInstId            string  `json:"prodInstId,omitempty"`                /*  实例ID  */
	Vhost                 string  `json:"vhost,omitempty"`                     /*  vhost名称  */
	Name                  string  `json:"name,omitempty"`                      /*  队列名称  */
	Durable               *bool   `json:"durable"`                             /*  是否持久化  */
	Auto_delete           *bool   `json:"auto_delete"`                         /*  是否自动删除  */
	Node                  *string `json:"node,omitempty"`                      // 队列所在节点，可选，默认为实例随机节点
	XExpires              *int64  `json:"x-expires,omitempty"`                 // 队列过期时间，单位ms，过期后自动删除，可选
	XDeadLetterExchange   *string `json:"x-dead-letter-exchange,omitempty"`    // 死信交换器名称，可选
	XDeadLetterRoutingKey *string `json:"x-dead-letter-routing-key,omitempty"` // 死信路由键，可选
	XMessageTTL           *int64  `json:"x-message-ttl,omitempty"`             // 消息过期时间，单位ms，可选
	XMaxLength            *int64  `json:"x-max-length,omitempty"`              // 队列最大消息长度，可选
	XMaxLengthBytes       *int64  `json:"x-max-length-bytes,omitempty"`        // 队列消息的总字节数上限，可选
	XOverflow             *string `json:"x-overflow,omitempty"`                // 队列消息处理策略，可选，值为drop-head或reject-publish
	XMaxPriority          *int64  `json:"x-max-priority,omitempty"`            // 队列最大优先级，范围0~255，可选
	XQueueMode            *string `json:"x-queue-mode,omitempty"`              // 队列模式，可选，值为default或lazy
}

type AmqpQueueCreateV3Response struct {
	ReturnObj  *AmqpQueueCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                              `json:"message"`    /*  描述状态  */
	StatusCode string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                              `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpQueueCreateV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
