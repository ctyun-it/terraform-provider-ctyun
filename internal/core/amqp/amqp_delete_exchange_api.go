package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpDeleteExchangeApi
/* 删除交换器
 */type AmqpDeleteExchangeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpDeleteExchangeApi(client *core.CtyunClient) *AmqpDeleteExchangeApi {
	return &AmqpDeleteExchangeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/exchange/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpDeleteExchangeApi) Do(ctx context.Context, credential core.Credential, req *AmqpDeleteExchangeRequest) (*AmqpDeleteExchangeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpDeleteExchangeRequest
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
	var resp AmqpDeleteExchangeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpDeleteExchangeRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Name       string `json:"name,omitempty"`       /*  交换机名称  */
}

type AmqpDeleteExchangeResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
