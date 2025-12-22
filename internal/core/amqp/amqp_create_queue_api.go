package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpCreateQueueApi
/* 创建队列
 */type AmqpCreateQueueApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpCreateQueueApi(client *core.CtyunClient) *AmqpCreateQueueApi {
	return &AmqpCreateQueueApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/queue/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpCreateQueueApi) Do(ctx context.Context, credential core.Credential, req *AmqpCreateQueueRequest) (*AmqpCreateQueueResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpCreateQueueRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpCreateQueueResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpCreateQueueRequest struct {
	ProdInstId            string `json:"prodInstId,omitempty"`                /*  实例ID  */
	Vhost                 string `json:"vhost,omitempty"`                     /*  vhost名称  */
	Name                  string `json:"name,omitempty"`                      /*  队列名称  */
	Durable               *bool  `json:"durable"`                             /*  是否持久化  */
	Auto_delete           *bool  `json:"auto_delete"`                         /*  是否自动删除  */
	XMessageTtl           int32  `json:"x-message-ttl,omitempty"`             /*  消息过期时间（单位ms）  */
	XExpires              int32  `json:"x-expires,omitempty"`                 /*  队列过期时间，过期后队列自动删除 （单位ms）  */
	XMaxLength            int32  `json:"x-max-length,omitempty"`              /*  队列能保存的最大消息数  */
	XDeadLetterExchange   string `json:"x-dead-letter-exchange,omitempty"`    /*  死信交换器名称，（死信交换机必须提前创建，否则接口返回成功但不会真正创建队列）  */
	XDeadLetterRoutingKey string `json:"x-dead-letter-routing-key,omitempty"` /*  死信路由键  */
	XMaxPriority          int32  `json:"x-max-priority,omitempty"`            /*  队列最大优先级：要开启消息的优先级，必须设置消息所在队列的优先级，0\~ 255  */
}

type AmqpCreateQueueResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
