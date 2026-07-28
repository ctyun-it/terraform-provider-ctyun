package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpDeleteBindingApi
/* 删除绑定
 */type AmqpDeleteBindingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpDeleteBindingApi(client *core.CtyunClient) *AmqpDeleteBindingApi {
	return &AmqpDeleteBindingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/binding/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpDeleteBindingApi) Do(ctx context.Context, credential core.Credential, req *AmqpDeleteBindingRequest) (*AmqpDeleteBindingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddHeader("customInfo", req.CustomInfo)
	_, err := ctReq.WriteJson(struct {
		*AmqpDeleteBindingRequest
		RegionId   interface{} `json:"regionId,omitempty"`
		CustomInfo interface{} `json:"customInfo,omitempty"`
	}{
		req, nil, nil,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpDeleteBindingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpDeleteBindingRequest struct {
	RegionId         string `json:"regionId,omitempty"`         /*  资源池ID  */
	CustomInfo       string `json:"customInfo,omitempty"`       /*  请参照公共参数说明及customInfo对象参数说明进行传参（accountId、userId必传）  */
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例ID  */
	Destination_type string `json:"destination_type,omitempty"` /*  绑定目标类型,交换机:e，队列q  */
	Source           string `json:"source,omitempty"`           /*  绑定来源交换器名称  */
	Destination      string `json:"destination,omitempty"`      /*  绑定目标名称  */
	Vhost            string `json:"vhost,omitempty"`            /*  vhost名称  */
	Properties_key   string `json:"properties_key,omitempty"`   /*  对应创建绑定接口routing_key字段  */
}

type AmqpDeleteBindingResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    string `json:"message"`    /*  描述状态  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
