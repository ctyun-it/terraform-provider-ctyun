package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2ConsumerResetoffsetV3Api
/* 重置订阅组消费进度到指定时间戳 */
type Mq2ConsumerResetoffsetV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ConsumerResetoffsetV3Api(client *core.CtyunClient) *Mq2ConsumerResetoffsetV3Api {
	return &Mq2ConsumerResetoffsetV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/consumer/resetOffset",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2ConsumerResetoffsetV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2ConsumerResetoffsetV3Request) (*Mq2ConsumerResetoffsetV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2ConsumerResetoffsetV3Request
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
	var resp Mq2ConsumerResetoffsetV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ConsumerResetoffsetV3Request struct {
	RegionId   string `json:"regionId"`   /*  资源池编码  */
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
	ResetTime  int64  `json:"resetTime"`  /*  要重置的毫秒时间戳  */
}

type Mq2ConsumerResetoffsetV3Response struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象  */
	Error      *string `json:"error"`      /*  错误码，只有非成功才有这个字段，方便快速定位问题  */
}
