package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaConsumerGroupStatusV3Api
/* 查询消费组状态。
 */type CtgkafkaConsumerGroupStatusV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaConsumerGroupStatusV3Api(client *core.CtyunClient) *CtgkafkaConsumerGroupStatusV3Api {
	return &CtgkafkaConsumerGroupStatusV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumerGroup/status",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaConsumerGroupStatusV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaConsumerGroupStatusV3Request) (*CtgkafkaConsumerGroupStatusV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("groupName", req.GroupName)
	ctReq.AddParam("topicName", req.TopicName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaConsumerGroupStatusV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaConsumerGroupStatusV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	GroupName  string `json:"groupName,omitempty"`  /*  消费组名称。  */
	TopicName  string `json:"topicName,omitempty"`  /*  topic名称。  */
}

type CtgkafkaConsumerGroupStatusV3Response struct {
	StatusCode string                                          `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                          `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaConsumerGroupStatusV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaConsumerGroupStatusV3ReturnObjResponse struct {
	Total int32                                                 `json:"total,omitempty"` /*  总记录数。  */
	Data  []*CtgkafkaConsumerGroupStatusV3ReturnObjDataResponse `json:"data"`            /*  消费组所订阅的主题状态列表。  */
}

type CtgkafkaConsumerGroupStatusV3ReturnObjDataResponse struct {
	TopicName   string                                                     `json:"topicName,omitempty"`   /*  topic名称。  */
	PartitionId int32                                                      `json:"partitionId,omitempty"` /*  主题的分区序号。  */
	Offsets     *CtgkafkaConsumerGroupStatusV3ReturnObjDataOffsetsResponse `json:"offsets"`               /*  分区的消费进度状态。  */
}

type CtgkafkaConsumerGroupStatusV3ReturnObjDataOffsetsResponse struct {
	Begin   int64 `json:"begin,omitempty"`   /*  最小点位。  */
	End     int64 `json:"end,omitempty"`     /*  最大点位。  */
	Current int64 `json:"current,omitempty"` /*  消费组消费当前位点。  */
	Lag     int64 `json:"lag,omitempty"`     /*  消息堆积。  */
	Total   int64 `json:"total,omitempty"`   /*  分区消息数总量。  */
	Hw      int64 `json:"hw,omitempty"`      /*  消息高水位线，所有副本均已确认写入的最高偏移量。  */
}
