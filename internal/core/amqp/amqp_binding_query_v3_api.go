package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpBindingQueryV3Api
/* 绑定查询
 */type AmqpBindingQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpBindingQueryV3Api(client *core.CtyunClient) *AmqpBindingQueryV3Api {
	return &AmqpBindingQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/binding/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpBindingQueryV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpBindingQueryV3Request) (*AmqpBindingQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpBindingQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpBindingQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
}

type AmqpBindingQueryV3Response struct {
	ReturnObj  *AmqpBindingQueryV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                               `json:"message"`    /*  描述状态  */
	StatusCode string                               `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                               `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpBindingQueryV3ReturnObjResponse struct {
	Source           string `json:"source"`           /*  绑定来源，交换器名称  */
	Vhost            string `json:"vhost"`            /*  虚拟机  */
	Destination      string `json:"destination"`      /*  绑定目标，队列名称  */
	Destination_type string `json:"destination_type"` /*  绑定目标类型，只有一种类型queue  */
	Routing_key      string `json:"routing_key"`      /*  路由键  */
	Arguments        string `json:"arguments"`        /*  绑定参数  */
	Properties_key   string `json:"properties_key"`   /*  properties键  */
}
