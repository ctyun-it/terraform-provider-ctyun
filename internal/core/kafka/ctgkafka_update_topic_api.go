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
	RegionId                    string `json:"regionId,omitempty"`          /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId                  string `json:"prodInstId,omitempty"`        /*  实例ID。  */
	TopicName                   string `json:"topicName,omitempty"`         /*  主题名称。  */
	PartitionNum                int32  `json:"partitionNum,omitempty"`      /*  主题分区数，取值范围[1,100]。  */
	PartitionCapacity           int32  `json:"partitionCapacity,omitempty"` /*  分区容量限制，单位GB，取值-1或范围[1, 100]。-1表示无限制。不传入则不修改。  */
	RetentionTime               int32  `json:"retentionTime,omitempty"`     /*  消息保留时长，单位毫秒，取值-1或范围[36000, 315360000000]，单位毫秒，-1表示永久保留。不传入则不修改。  */
	MinReplicas                 int32  `json:"minReplicas,omitempty"`       /*  最小同步副本数，需小于等于factorNum。  */
	MaxMessage                  int32  `json:"maxMessage,omitempty"`        /*  最大消息大小，单位字节，取值范围[1, 10485760]。不传入则不修改。  */
	NeedFlush                   *bool  `json:"needFlush"`                   /*  是否同步刷盘。<br><li>true：是<br><li>false：否<br>  */
	TimestampType               string `json:"timestampType,omitempty"`     /*  消息时间戳类型，不传入则不修改。<br><li>CreateTime<br><li>LogAppendTime  */
	Description                 string `json:"description,omitempty"`       /*  主题描述，规则如下：<br><li>不能以+,-,@,= 特殊字符开头。 <br><li>长度不能大于200。  */
	StrategyName                string `json:"strategyName,omitempty"`      /*  预设策略名称。  */
	CleanupPolicy               string `json:"cleanupPolicy,omitempty"`     /*  日志保留策略。<br><li>delete<br><li>compact  */
	UncleanLeaderElectionEnable *bool  `json:"uncleanLeaderElectionEnable"` /*  是否允许不同步的副本参与leader选举。<br><li>false<br><li>true  */
	SegmentMs                   int64  `json:"segmentMs,omitempty"`         /*  日志滚动时间，单位ms。 取值范围[86400000, 7776000000]  */
	SegmentBytes                int64  `json:"segmentBytes,omitempty"`      /*  分片大小，单位byte。 取值范围[268435456, 10737418240]  */
	RemoteStorageEnable         *bool  `json:"remoteStorageEnable"`         /*  是否开启对象存储。<br><li>true：是<br><li>false：否<br>  */
	LocalRetentionMs            int64  `json:"localRetentionMs,omitempty"`  /*  本地保留时长，单位ms。 取值范围[180000, 315360000000]  */
}

type CtgkafkaUpdateTopicResponse struct {
	StatusCode string                                `json:"statusCode"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                `json:"message"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaUpdateTopicReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                `json:"error"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUpdateTopicReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据。  */
}
