package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQuotaTopicUpdateApi
/* 更新Topic流控配置。
 */type CtgkafkaQuotaTopicUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQuotaTopicUpdateApi(client *core.CtyunClient) *CtgkafkaQuotaTopicUpdateApi {
	return &CtgkafkaQuotaTopicUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/quota/topic/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaQuotaTopicUpdateApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQuotaTopicUpdateRequest) (*CtgkafkaQuotaTopicUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaQuotaTopicUpdateRequest
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
	var resp CtgkafkaQuotaTopicUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQuotaTopicUpdateRequest struct {
	RegionId         string `json:"regionId,omitempty"`         /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例ID，实例需是Kafka引擎类型且非单机版。  */
	Topic            string `json:"topic,omitempty"`            /*  主题名称。  */
	ProducerByteRate int32  `json:"producerByteRate,omitempty"` /*  生产者流控上限，单位字节，取值范围[1048576, 1073741824]，与consumerByteRate至少一个不为空。  */
	ConsumerByteRate int32  `json:"consumerByteRate,omitempty"` /*  消费者流控上限，单位字节，取值范围[1048576, 1073741824]，与producerByteRate至少一个不为空。  */
}

type CtgkafkaQuotaTopicUpdateResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string `json:"message,omitempty"`    /*  描述状态。  */
}
