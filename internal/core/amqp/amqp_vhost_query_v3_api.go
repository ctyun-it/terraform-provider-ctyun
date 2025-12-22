package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpVhostQueryV3Api
/* 查询虚拟机
 */type AmqpVhostQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpVhostQueryV3Api(client *core.CtyunClient) *AmqpVhostQueryV3Api {
	return &AmqpVhostQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/vhost/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpVhostQueryV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpVhostQueryV3Request) (*AmqpVhostQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpVhostQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpVhostQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例id  */
}

type AmqpVhostQueryV3Response struct {
	ReturnObj  *AmqpVhostQueryV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                             `json:"message"`    /*  描述状态  */
	StatusCode string                             `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                             `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpVhostQueryV3ReturnObjResponse struct {
	Data *AmqpVhostQueryV3ReturnObjDataResponse `json:"data"` /*  返回数据  */
}

type AmqpVhostQueryV3ReturnObjDataResponse struct {
	Vhosts       []string           `json:"vhosts"` /*  虚拟机列表  */
	VhostsDetail []AmqpVhostsDetail `json:"vhostsDetail"`
}

type AmqpVhostsDetail struct {
	Name                   string      `json:"name"`
	Messages               interface{} `json:"messages"`
	MessagesReady          interface{} `json:"messagesReady"`
	MessagesUnacknowledged interface{} `json:"messagesUnacknowledged"`
	PublishBytesRate       float64     `json:"publishBytesRate"`
	DeliverBytesRate       float64     `json:"deliverBytesRate"`
	PublishMessagesRate    float64     `json:"publishMessagesRate"`
	DeliverMessagesRate    float64     `json:"deliverMessagesRate"`
}
