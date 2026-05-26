package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Mq2TopicCreateV3Api
/* 创建主题 */
type Mq2TopicCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewMq2TopicCreateV3Api(client *core.CtyunClient) *Mq2TopicCreateV3Api {
	return &Mq2TopicCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/create",
			ContentType:  "application/json",
		},
	}
}

func (a *Mq2TopicCreateV3Api) Do(ctx context.Context, credential core.Credential, req *Mq2TopicCreateV3Request) (*Mq2TopicCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*Mq2TopicCreateV3Request
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
	var resp Mq2TopicCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type Mq2TopicCreateV3Request struct {
	RegionId             string   `json:"regionId"`             /*  资源池编码  */
	ProdInstId           string   `json:"prodInstId"`           /*  实例id  */
	BrokerNameList       []string `json:"brokerNameList"`       /*  关联的Broker节点名称列表，需填写实际部署的Broker节点名称，不可为空数组  */
	WriteQueueNums       int32    `json:"writeQueueNums"`       /*  主题的写队列数量，用于控制消息写入的并发能力，需为正整数  */
	ReadQueueNums        int32    `json:"readQueueNums"`        /*  主题的读队列数量，需与写队列数量保持一致，确保消息消费的负载均衡，需为正整数  */
	Order                bool     `json:"order"`                /*  标识主题是否为顺序消息队列，true表示顺序队列，false表示普通队列  */
	Perm                 int32    `json:"perm"`                 /*  主题的权限控制值，固定可选值为2（只读）、4（只写）、6（读写），默认推荐6  */
	AllowdConsumerGroups []string `json:"allowdConsumerGroups"` /*  允许订阅该主题的消费者组列表，空数组表示不限制消费者组订阅  */
	TopicName            string   `json:"topicName"`            /*  主题名称，需符合MQ主题命名规范（字母、数字、下划线组合，长度1-64字符），不可重复  */
	Remark               string   `json:"remark"`               /*  主题备注  */
}

type Mq2TopicCreateV3Response struct {
	StatusCode int32                              `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    *string                            `json:"message"`    /*  描述状态。  */
	ReturnObj  *Mq2TopicCreateV3ReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      *string                            `json:"error"`      /*  错误码，描述错误信息。  */
}

type Mq2TopicCreateV3ReturnObjResponse struct{}
