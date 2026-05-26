package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2MessagePushV3Api
/* 向指定的消费者推送消息 */
type Mq2MessagePushV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2MessagePushV3Api(client *core.CtyunClient) *Mq2MessagePushV3Api {
	return &Mq2MessagePushV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/message/push",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2MessagePushV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2MessagePushV3Request) (*Mq2MessagePushV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2MessagePushV3Request
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
	var resp Mq2MessagePushV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2MessagePushV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	MsgId      string `json:"msgId"`      /*  消息的offsetMsgId  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	ClientId   string `json:"clientId"`   /*  客户端ID  */
}

type Mq2MessagePushV3Response struct {
	StatusCode *string                            `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string                            `json:"message"`    /*  描述状态  */
	ReturnObj  *Mq2MessagePushV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例'里面的注释  */
	Error      *string                            `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}

type Mq2MessagePushV3ReturnObjResponse struct{}
