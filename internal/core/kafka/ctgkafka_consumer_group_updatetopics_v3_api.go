package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaConsumerGroupUpdatetopicsV3Api
/* 更新消费组订阅的主题。
 */type CtgkafkaConsumerGroupUpdatetopicsV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaConsumerGroupUpdatetopicsV3Api(client *core.CtyunClient) *CtgkafkaConsumerGroupUpdatetopicsV3Api {
	return &CtgkafkaConsumerGroupUpdatetopicsV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/consumerGroup/updateTopics",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaConsumerGroupUpdatetopicsV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaConsumerGroupUpdatetopicsV3Request) (*CtgkafkaConsumerGroupUpdatetopicsV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaConsumerGroupUpdatetopicsV3Request
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
	var resp CtgkafkaConsumerGroupUpdatetopicsV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaConsumerGroupUpdatetopicsV3Request struct {
	RegionId         string   `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId       string   `json:"prodInstId,omitempty"` /*  实例ID。  */
	GroupName        string   `json:"groupName,omitempty"`  /*  消费组名称。  */
	AddTopicNames    []string `json:"addTopicNames"`        /*  需要订阅的主题，和deleteTopicNames不能同时为空。  */
	DeleteTopicNames []string `json:"deleteTopicNames"`     /*  取消订阅的主题，和addTopicNames不能同时为空。  */
}

type CtgkafkaConsumerGroupUpdatetopicsV3Response struct {
	StatusCode string                                                `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                                `json:"message,omitempty"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaConsumerGroupUpdatetopicsV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                                `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaConsumerGroupUpdatetopicsV3ReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
