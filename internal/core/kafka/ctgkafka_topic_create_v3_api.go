package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaTopicCreateV3Api
/* 创建主题。
 */type CtgkafkaTopicCreateV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTopicCreateV3Api(client *core.CtyunClient) *CtgkafkaTopicCreateV3Api {
	return &CtgkafkaTopicCreateV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaTopicCreateV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTopicCreateV3Request) (*CtgkafkaTopicCreateV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaTopicCreateV3Request
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
	var resp CtgkafkaTopicCreateV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaTopicCreateV3Request struct {
	RegionId          string `json:"regionId,omitempty"`          /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId        string `json:"prodInstId,omitempty"`        /*  实例ID。  */
	TopicName         string `json:"topicName,omitempty"`         /*  主题名称，英文字母、数字、下划线开头，且只能由英文字母、数字、中划线、下划线组成，长度为3-64个字符。  */
	PartitionNum      int32  `json:"partitionNum,omitempty"`      /*  分区数，取值范围[1, min(100, 实例剩余分区数量)]，实例剩余分区数量=实例分区上限-所有主题分区数之和。<br>您可以<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=330">查询实例</a>获取对应实例下partitionNum获取分区数上限。  */
	FactorNum         int32  `json:"factorNum,omitempty"`         /*  副本数，取值范围[1, 3]，单机版默认值1，集群版默认值3。  */
	PartitionCapacity int32  `json:"partitionCapacity,omitempty"` /*  分区容量限制，单位GB，取值-1或范围[1, 100]。-1表示无限制，默认值-1。  */
	RetentionTime     int32  `json:"retentionTime,omitempty"`     /*  消息保留时长，单位毫秒，取值-1或范围[3600000, 315360000000]，单位毫秒，-1表示永久保留。 默认值259200000。  */
	MinReplicas       int32  `json:"minReplicas,omitempty"`       /*  最小同步副本数，需小于等于factorNum，单机版默认值1，集群版默认值min(2, factorNum)。  */
	MaxMessage        int32  `json:"maxMessage,omitempty"`        /*  最大消息大小，单位字节，取值范围[1, 10485760]， 默认值1048588。  */
	NeedFlush         *bool  `json:"needFlush"`                   /*  是否同步刷盘。<br><li>true：是<br><li>false：否<br><li>默认值false  */
	TimestampType     string `json:"timestampType,omitempty"`     /*  消息时间戳类型。<br><li>CreateTime<br><li>LogAppendTime<br><li>默认值CreateTime  */
	Description       string `json:"description,omitempty"`       /*  主题描述，规则如下：<br><li>不能以+,-,@,= 特殊字符开头。 <br><li>长度不能大于200。  */
}

type CtgkafkaTopicCreateV3Response struct {
	StatusCode string                                  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                  `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTopicCreateV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTopicCreateV3ReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
