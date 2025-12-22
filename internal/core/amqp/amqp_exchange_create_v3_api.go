package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpExchangeCreateV3Api
/* 创建交换器v3
 */type AmqpExchangeCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpExchangeCreateV3Api(client *core.CtyunClient) *AmqpExchangeCreateV3Api {
	return &AmqpExchangeCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/exchange/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpExchangeCreateV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpExchangeCreateV3Request) (*AmqpExchangeCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpExchangeCreateV3Request
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
	var resp AmqpExchangeCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpExchangeCreateV3Request struct {
	RegionId          string  `json:"regionId,omitempty"`   /*  资源池id  */
	ProdInstId        string  `json:"prodInstId,omitempty"` /*  实例ID  */
	Vhost             string  `json:"vhost,omitempty"`      /*  vhost名称(vhost需提前创建，否则接口放回成功但不会真正创建交换机)   */
	Name              string  `json:"name,omitempty"`       /*  交换器名称  */
	Auto_delete       *bool   `json:"auto_delete"`          /*  是否自动删除  */
	RawType           string  `json:"type,omitempty"`       /*  交换器类型  */
	AlternateExchange *string `json:"alternate-exchange,omitempty"`
	XDelayedType      *string `json:"x-delayed-type,omitempty"`
	Durable           *bool   `json:"durable,omitempty"`
	Internal          *bool   `json:"internal,omitempty"`
}

type AmqpExchangeCreateV3Response struct {
	ReturnObj  *AmqpExchangeCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象  */
	Message    string                                 `json:"message"`    /*  描述状态  */
	StatusCode string                                 `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string                                 `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}

type AmqpExchangeCreateV3ReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据  */
}
