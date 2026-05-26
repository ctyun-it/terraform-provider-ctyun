package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2CreateTopicApi
/* 创建主题 */
type Mq2CreateTopicApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2CreateTopicApi(client *core.CtyunClient) *Mq2CreateTopicApi {
	return &Mq2CreateTopicApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/topic/create",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2CreateTopicApi) Do(ctx context.Context, credential core.Credential, req *Mq2CreateTopicRequest) (*Mq2CreateTopicResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*Mq2CreateTopicRequest
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
	var resp Mq2CreateTopicResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2CreateTopicRequest struct {
	ProdInstId           string    `json:"prodInstId"`                     /*  实例ID  */
	BrokerNameList       []string  `json:"brokerNameList"`                 /*  需要创建topic的Broker名字列表  */
	AllowdConsumerGroups []*string `json:"allowdConsumerGroups,omitempty"` /*  1  */
	TopicName            string    `json:"topicName"`                      /*  topic名字  */
	WriteQueueNums       int32     `json:"writeQueueNums"`                 /*  topic的属性，可写queue数量，默认是16  */
	ReadQueueNums        int32     `json:"readQueueNums"`                  /*  topic的属性，可读queue数量，默认是16  */
	Perm                 int32     `json:"perm"`                           /*  权限：2可写，4可读，6可读写。 默认可读写  */
	Order                bool      `json:"order"`                          /*  是否有序topic，默认是false  */
	Remark               *string   `json:"remark,omitempty"`               /*  主题备注信息  */
}

type Mq2CreateTopicResponse struct {
	StatusCode *string `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900  */
	Message    *string `json:"message"`    /*  描述状态  */
	ReturnObj  *string `json:"returnObj"`  /*  返回对象。此参数所包含的参数请见“响应示例"里面的注释  */
}
