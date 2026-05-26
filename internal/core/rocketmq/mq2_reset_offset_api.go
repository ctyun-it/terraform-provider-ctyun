package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// Mq2ResetOffsetApi
/* 重置订阅组消费进度到指定时间戳 */
type Mq2ResetOffsetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2ResetOffsetApi(client *core.CtyunClient) *Mq2ResetOffsetApi {
	return &Mq2ResetOffsetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/consumer/resetOffset",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2ResetOffsetApi) Do(ctx context.Context, credential core.Credential, req *Mq2ResetOffsetRequest) (*Mq2ResetOffsetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("resetTime", strconv.FormatInt(int64(req.ResetTime), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2ResetOffsetResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2ResetOffsetRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	GroupName  string `json:"groupName"`  /*  订阅组名字  */
	TopicName  string `json:"topicName"`  /*  Topic名字  */
	ResetTime  int64  `json:"resetTime"`  /*  要重置的毫秒时间戳  */
}

type Mq2ResetOffsetResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
