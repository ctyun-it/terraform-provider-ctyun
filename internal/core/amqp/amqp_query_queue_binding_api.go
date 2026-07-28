package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryQueueBindingApi
/* 查队列绑定
 */type AmqpQueryQueueBindingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryQueueBindingApi(client *core.CtyunClient) *AmqpQueryQueueBindingApi {
	return &AmqpQueryQueueBindingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/binding/queue/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpQueryQueueBindingApi) Do(ctx context.Context, credential core.Credential, req *AmqpQueryQueueBindingRequest) (*AmqpQueryQueueBindingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	ctReq.AddParam("queue", req.Queue)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpQueryQueueBindingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueryQueueBindingRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Queue      string `json:"queue,omitempty"`      /*  队列名称  */
}

type AmqpQueryQueueBindingResponse struct {
	StatusCode string                                    `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string                                    `json:"message"`    /*  描述状态  */
	ReturnObj  []*AmqpQueryQueueBindingReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}

type AmqpQueryQueueBindingReturnObjResponse struct{}
