package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2DeleteTopicApi
/* 删除主题 */
type Mq2DeleteTopicApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2DeleteTopicApi(client *core.CtyunClient) *Mq2DeleteTopicApi {
	return &Mq2DeleteTopicApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/topic/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2DeleteTopicApi) Do(ctx context.Context, credential core.Credential, req *Mq2DeleteTopicRequest) (*Mq2DeleteTopicResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2DeleteTopicRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Mq2DeleteTopicResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2DeleteTopicRequest struct {
	ProdInstId string `json:"prodInstId"` /*  实例ID  */
	TopicName  string `json:"topicName"`  /*  topic名字  */
}

type Mq2DeleteTopicResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
