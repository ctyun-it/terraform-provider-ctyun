package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2PushApi
/* 向指定的消费者推送消息 */
type Mq2PushApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2PushApi(client *core.CtyunClient) *Mq2PushApi {
	return &Mq2PushApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/message/push",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2PushApi) Do(ctx context.Context, credential core.Credential, req *Mq2PushRequest) (*Mq2PushResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("msgId", req.MsgId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("clientId", req.ClientId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2PushResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2PushRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	MsgId      string `json:"msgId"`      /*  消息ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	ClientId   string `json:"clientId"`   /*  客户端ID  */
}

type Mq2PushResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
