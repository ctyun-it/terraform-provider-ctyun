package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpCreateVhostApi
/* 创建虚拟机
 */type AmqpCreateVhostApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpCreateVhostApi(client *core.CtyunClient) *AmqpCreateVhostApi {
	return &AmqpCreateVhostApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/vhost/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpCreateVhostApi) Do(ctx context.Context, credential core.Credential, req *AmqpCreateVhostRequest) (*AmqpCreateVhostResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*AmqpCreateVhostRequest
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
	var resp AmqpCreateVhostResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpCreateVhostRequest struct {
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Name       string `json:"name,omitempty"`       /*  虚拟机名称  */
}

type AmqpCreateVhostResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
