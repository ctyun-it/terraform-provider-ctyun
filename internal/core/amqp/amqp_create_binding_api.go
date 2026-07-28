package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpCreateBindingApi
/* 绑定
 */type AmqpCreateBindingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpCreateBindingApi(client *core.CtyunClient) *AmqpCreateBindingApi {
	return &AmqpCreateBindingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/binding/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpCreateBindingApi) Do(ctx context.Context, credential core.Credential, req *AmqpCreateBindingRequest) (*AmqpCreateBindingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpCreateBindingRequest
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
	var resp AmqpCreateBindingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpCreateBindingRequest struct {
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例ID  */
	Destination_type string `json:"destination_type,omitempty"` /*  绑定目标类型,交换机:e，队列q  */
	Source           string `json:"source,omitempty"`           /*  绑定来源交换器名称（交换机需存在）  */
	Destination      string `json:"destination,omitempty"`      /*  绑定目标名称取值如下： destination_type=e 交换机名称 destination_type=q 队列名称  */
	Vhost            string `json:"vhost,omitempty"`            /*  vhost名称  */
	Routing_key      string `json:"routing_key,omitempty"`      /*  路由键  */
	Arguments        string `json:"arguments,omitempty"`        /*  绑定参数。该参数仅用于headers类型交换机，此时必须包括键x-match，可选值有all、any。  */
}

type AmqpCreateBindingResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900   */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
