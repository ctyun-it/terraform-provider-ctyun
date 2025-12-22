package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpBindingCreateV3Api
/* 绑定
 */type AmqpBindingCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpBindingCreateV3Api(client *core.CtyunClient) *AmqpBindingCreateV3Api {
	return &AmqpBindingCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/binding/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpBindingCreateV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpBindingCreateV3Request) (*AmqpBindingCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpBindingCreateV3Request
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
	var resp AmqpBindingCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpBindingCreateV3Request struct {
	RegionId         string `json:"regionId,omitempty"`         /*  资源池id  */
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例id  */
	Source           string `json:"source,omitempty"`           /*  绑定来源交换器名称（交换机需存在)  */
	Destination_type string `json:"destination_type,omitempty"` /*  绑定目标类型，只支持q  */
	Destination      string `json:"destination,omitempty"`      /*  队列名称 （需存在）  */
	Routing_key      string `json:"routing_key,omitempty"`      /*  路由键  */
	Vhost            string `json:"vhost,omitempty"`            /*  vhost名称  */
	Arguments        string `json:"arguments,omitempty"`        /*  绑定参数  */
}

type AmqpBindingCreateV3Response struct {
	ReturnObj  *AmqpBindingCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                                `json:"message"`    /*  描述状态  */
	StatusCode string                                `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                                `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpBindingCreateV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
