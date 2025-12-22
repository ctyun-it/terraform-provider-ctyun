package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpVhostCreateV3Api
/* 创建虚拟机
 */type AmqpVhostCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpVhostCreateV3Api(client *core.CtyunClient) *AmqpVhostCreateV3Api {
	return &AmqpVhostCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/vhost/create",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpVhostCreateV3Api) Do(ctx context.Context, credential core.Credential, req *AmqpVhostCreateV3Request) (*AmqpVhostCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpVhostCreateV3Request
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
	var resp AmqpVhostCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpVhostCreateV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  资源池id  */
	Name       string `json:"name,omitempty"`       /*  虚拟机名称  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID  */
}

type AmqpVhostCreateV3Response struct {
	Message    string `json:"message"`    /*  描述状态  */
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Error      string `json:"error"`      /*  错误码，描述错误信息，只有失败才显示  */
}
