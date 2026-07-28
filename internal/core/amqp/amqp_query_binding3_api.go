package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueryBinding3Api
/* 查队列绑定
 */type AmqpQueryBinding3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueryBinding3Api(client *core.CtyunClient) *AmqpQueryBinding3Api {
	return &AmqpQueryBinding3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/binding/queue/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpQueryBinding3Api) Do(ctx context.Context, credential core.Credential, req *AmqpQueryBinding3Request) (*AmqpQueryBinding3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	ctReq.AddParam("queue", req.Queue)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpQueryBinding3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueryBinding3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
	Vhost      string `json:"vhost,omitempty"`      /*  虚拟主机  */
	Queue      string `json:"queue,omitempty"`      /*  队列名称  */
}

type AmqpQueryBinding3Response struct {
	ReturnObj  *AmqpQueryBinding3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                              `json:"message"`    /*  描述状态  */
	StatusCode string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                              `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpQueryBinding3ReturnObjResponse struct {
	Data []*AmqpQueryBinding3ReturnObjDataResponse `json:"data"` /*  绑定数据  */
}

type AmqpQueryBinding3ReturnObjDataResponse struct {
	Source           string                                           `json:"source"`           /*  交换器  */
	Vhost            string                                           `json:"vhost"`            /*  虚拟机  */
	Destination      string                                           `json:"destination"`      /*  队列  */
	Destination_type string                                           `json:"destination_type"` /*  绑定目标类型，只支持queue  */
	Routing_key      string                                           `json:"routing_key"`      /*  路由键  */
	Arguments        *AmqpQueryBinding3ReturnObjDataArgumentsResponse `json:"arguments"`        /*  绑定参数  */
	Properties_key   string                                           `json:"properties_key"`   /*  properties键  */
}

type AmqpQueryBinding3ReturnObjDataArgumentsResponse struct{}
