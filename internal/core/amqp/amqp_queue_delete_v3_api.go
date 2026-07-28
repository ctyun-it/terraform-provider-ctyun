package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpQueueDeleteV3Api
/* 删除队列v3
 */type AmqpQueueDeleteV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpQueueDeleteV3Api(client *core.CtyunClient) *AmqpQueueDeleteV3Api {
	return &AmqpQueueDeleteV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/queue/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpQueueDeleteV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpQueueDeleteV3Request) (*AmqpQueueDeleteV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpQueueDeleteV3Request
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
	var resp AmqpQueueDeleteV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpQueueDeleteV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost      string `json:"vhost,omitempty"`      /*  虚拟机名称  */
	Name       string `json:"name,omitempty"`       /*  队列名称  */
}

type AmqpQueueDeleteV3Response struct {
	ReturnObj  *AmqpQueueDeleteV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                              `json:"message"`    /*  描述状态  */
	StatusCode string                              `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                              `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpQueueDeleteV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
