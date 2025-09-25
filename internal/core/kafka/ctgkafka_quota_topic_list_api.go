package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQuotaTopicListApi
/* 查询Topic流控配置。
 */type CtgkafkaQuotaTopicListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQuotaTopicListApi(client *core.CtyunClient) *CtgkafkaQuotaTopicListApi {
	return &CtgkafkaQuotaTopicListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/quota/topic/list",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaQuotaTopicListApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQuotaTopicListRequest) (*CtgkafkaQuotaTopicListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaQuotaTopicListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQuotaTopicListRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaQuotaTopicListResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                   `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaQuotaTopicListReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaQuotaTopicListReturnObjResponse struct {
	Topic            string `json:"topic,omitempty"`            /*  主题名称。  */
	ProducerByteRate int64  `json:"producerByteRate,omitempty"` /*  生产者流控上限，单位字节。  */
	ConsumerByteRate int64  `json:"consumerByteRate,omitempty"` /*  消费者流控上限，单位字节。  */
}
