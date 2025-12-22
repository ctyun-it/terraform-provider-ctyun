package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpChangeInstanceNameApi
/* 更改实例名称
 */type AmqpChangeInstanceNameApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpChangeInstanceNameApi(client *core.CtyunClient) *AmqpChangeInstanceNameApi {
	return &AmqpChangeInstanceNameApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/instances/instanceName",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpChangeInstanceNameApi) Do(ctx context.Context, credential core.Credential, req *AmqpChangeInstanceNameRequest) (*AmqpChangeInstanceNameResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpChangeInstanceNameRequest
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
	var resp AmqpChangeInstanceNameResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpChangeInstanceNameRequest struct {
	InstanceName string `json:"instanceName,omitempty"` /*  新的实例名称  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例Id  */
}

type AmqpChangeInstanceNameResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
