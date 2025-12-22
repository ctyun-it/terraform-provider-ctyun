package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpCreateExchangeApi
/* 创建交换器
 */type AmqpCreateExchangeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpCreateExchangeApi(client *core.CtyunClient) *AmqpCreateExchangeApi {
	return &AmqpCreateExchangeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/exchange/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpCreateExchangeApi) Do(ctx context.Context, credential core.Credential, req *AmqpCreateExchangeRequest) (*AmqpCreateExchangeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpCreateExchangeRequest
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
	var resp AmqpCreateExchangeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpCreateExchangeRequest struct {
	ProdInstId        string `json:"prodInstId,omitempty"`         /*  实例ID  */
	Vhost             string `json:"vhost,omitempty"`              /*  vhost名称(vhost需提前创建，否则接口放回成功但不会真正创建交换机)  */
	Name              string `json:"name,omitempty"`               /*  队列名称  */
	RawType           string `json:"type,omitempty"`               /*  交换器类型  */
	Durable           bool   `json:"durable"`                      /*  是否持久化  */
	Auto_delete       bool   `json:"auto_delete"`                  /*  是否自动删除  */
	Internal          bool   `json:"internal"`                     /*  是否为内部Exchange。取值： false：否 true：是  */
	AlternateExchange string `json:"alternate-exchange,omitempty"` /*  备用交换器  */
}

type AmqpCreateExchangeResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
