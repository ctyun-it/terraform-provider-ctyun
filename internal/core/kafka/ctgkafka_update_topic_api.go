package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUpdateTopicApi
/* 修改主题配置。
 */type CtgkafkaUpdateTopicApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUpdateTopicApi(client *core.CtyunClient) *CtgkafkaUpdateTopicApi {
	return &CtgkafkaUpdateTopicApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/updateTopic",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUpdateTopicApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUpdateTopicRequest) (*CtgkafkaUpdateTopicResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUpdateTopicRequest
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
	var resp CtgkafkaUpdateTopicResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUpdateTopicRequest struct {
	RegionId          string `json:"regionId,omitempty"`          /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId        string `json:"prodInstId,omitempty"`        /*  实例ID。  */
	TopicName         string `json:"topicName,omitempty"`         /*  主题名称。  */
	PartitionCapacity int32  `json:"partitionCapacity,omitempty"` /*  分区容量限制，单位GB，取值-1或范围[1, 100]。-1表示无限制。不传入则不修改。  */
	RetentionTime     int32  `json:"retentionTime,omitempty"`     /*  消息保留时长，单位毫秒，取值-1或范围[3600000, 315360000000]，单位毫秒，-1表示永久保留。不传入则不修改。  */
	MaxMessage        int32  `json:"maxMessage,omitempty"`        /*  最大消息大小，单位字节，取值范围[1, 10485760]。不传入则不修改。  */
	TimestampType     string `json:"timestampType,omitempty"`     /*  消息时间戳类型，不传入则不修改。<br><li>CreateTime<br><li>LogAppendTime  */
}

type CtgkafkaUpdateTopicResponse struct {
	StatusCode int32                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaUpdateTopicReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUpdateTopicReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
