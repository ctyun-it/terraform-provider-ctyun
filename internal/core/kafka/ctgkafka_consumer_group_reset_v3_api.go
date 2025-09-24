package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaConsumerGroupResetV3Api
/* 重置消费组消费点。
 */type CtgkafkaConsumerGroupResetV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaConsumerGroupResetV3Api(client *core.CtyunClient) *CtgkafkaConsumerGroupResetV3Api {
	return &CtgkafkaConsumerGroupResetV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/consumerGroup/reset",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaConsumerGroupResetV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaConsumerGroupResetV3Request) (*CtgkafkaConsumerGroupResetV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaConsumerGroupResetV3Request
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
	var resp CtgkafkaConsumerGroupResetV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaConsumerGroupResetV3Request struct {
	RegionId           string                                                   `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId         string                                                   `json:"prodInstId,omitempty"` /*  实例ID。  */
	GroupName          string                                                   `json:"groupName,omitempty"`  /*  消费组名称。  */
	TopicName          string                                                   `json:"topicName,omitempty"`  /*  主题名称。  */
	RawType            int32                                                    `json:"type,omitempty"`       /*  类型，<li>0：重置到latest。 <li>1：按时间重置。<li>2：重置到earliest。<li>3：按位点重置，此类型参数partitionShiftList为必填。  */
	PartitionShiftList []*CtgkafkaConsumerGroupResetV3PartitionShiftListRequest `json:"partitionShiftList"`   /*  位点重置列表，当type为3时必填。  */
	Time               int64                                                    `json:"time,omitempty"`       /*  重置时间点毫秒时间戳，type=1时必填。  */
}

type CtgkafkaConsumerGroupResetV3PartitionShiftListRequest struct {
	Partition int32 `json:"partition,omitempty"` /*  主题分区号。  */
	ShiftBy   int64 `json:"shiftBy,omitempty"`   /*  主题分区消费位点向左或向右移动的相对位置，例如当前offset是1000，当shiftBy=-10重置后offset=990，当shiftBy=10重置后offset=1010。  */
}

type CtgkafkaConsumerGroupResetV3Response struct {
	ReturnObj  *CtgkafkaConsumerGroupResetV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Message    string                                         `json:"message,omitempty"`    /*  描述状态。  */
	StatusCode string                                         `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaConsumerGroupResetV3ReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
