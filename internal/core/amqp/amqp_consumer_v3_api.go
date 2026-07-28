package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpConsumerV3Api
/* 查队列消费者v3
 */type AmqpConsumerV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpConsumerV3Api(client *core.CtyunClient) *AmqpConsumerV3Api {
	return &AmqpConsumerV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/queue/consumer",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *AmqpConsumerV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpConsumerV3Request) (*AmqpConsumerV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("vhost", req.Vhost)
	ctReq.AddParam("name", req.Name)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp AmqpConsumerV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpConsumerV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  vhost名称  */
	Name       string `json:"name,omitempty"`       /*  队列名称  */
}

type AmqpConsumerV3Response struct {
	ReturnObj  *AmqpConsumerV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                           `json:"message"`    /*  描述状态  */
	StatusCode string                           `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                           `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpConsumerV3ReturnObjResponse struct {
	Data []*AmqpConsumerV3ReturnObjDataResponse `json:"data"` /*  返回数据  */
}

type AmqpConsumerV3ReturnObjDataResponse struct {
	Ack_required    *bool  `json:"ack_required"`    /*  是否需要确认  */
	Active          *bool  `json:"active"`          /*  消费者是否活动  */
	Activity_status string `json:"activity_status"` /*  消费者活动状态  */
	Consumer_tag    string `json:"consumer_tag"`    /*  消费者tag  */
	Exclusive       *bool  `json:"exclusive"`       /*  是否专享  */
	Prefetch_count  int32  `json:"prefetch_count"`  /*  预拉取数  */
	Arguments       string `json:"arguments"`       /*  消费者参数  */
}
