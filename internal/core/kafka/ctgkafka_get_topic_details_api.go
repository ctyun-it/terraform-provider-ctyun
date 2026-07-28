package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaGetTopicDetailsApi
/* 查询主题详细信息。 */
type CtgkafkaGetTopicDetailsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaGetTopicDetailsApi(client *core.CtyunClient) *CtgkafkaGetTopicDetailsApi {
	return &CtgkafkaGetTopicDetailsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/getDetails",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaGetTopicDetailsApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaGetTopicDetailsRequest) (*CtgkafkaGetTopicDetailsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaGetTopicDetailsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CtgkafkaGetTopicDetailsRequest struct {
	RegionId   string `json:"regionId"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId"` /*  实例ID。  */
	TopicName  string `json:"topicName"`  /*  主题名称。  */
}

type CtgkafkaGetTopicDetailsResponse struct {
	StatusCode *string                                   `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    *string                                   `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaGetTopicDetailsReturnObjResponse `json:"returnObj,omitempty"`  /*  返回对象。  */
	Error      *string                                   `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaGetTopicDetailsReturnObjResponse struct {
	Data *CtgkafkaGetTopicDetailsReturnObjDataResponse `json:"data,omitempty"` /*  主题详情信息。  */
}

type CtgkafkaGetTopicDetailsReturnObjDataResponse struct {
	TopicName       *string                                                      `json:"topicName,omitempty"`       /*  主题名称。  */
	PartitionList   []*CtgkafkaGetTopicDetailsReturnObjDataPartitionListResponse `json:"partitionList,omitempty"`   /*  分区列表信息。  */
	GroupSubscribed []*string                                                    `json:"groupSubscribed,omitempty"` /*  订阅主题的消费组列表。  */
}

type CtgkafkaGetTopicDetailsReturnObjDataPartitionListResponse struct {
	TopicName   *string                                                              `json:"topicName,omitempty"` /*  主题名称。  */
	PartitionId int32                                                                `json:"partitionId"`         /*  分区ID。  */
	Offsets     *CtgkafkaGetTopicDetailsReturnObjDataPartitionListOffsetsResponse    `json:"offsets,omitempty"`   /*  分区偏移量信息。  */
	Replicas    []*CtgkafkaGetTopicDetailsReturnObjDataPartitionListReplicasResponse `json:"replicas,omitempty"`  /*  副本信息。  */
}

type CtgkafkaGetTopicDetailsReturnObjDataPartitionListOffsetsResponse struct {
	Total      int64 `json:"total"`      /*  分区消息总数。  */
	Begin      int64 `json:"begin"`      /*  分区leader副本的最大偏移量。  */
	End        int64 `json:"end"`        /*  分区leader副本的最小偏移量。  */
	UpdateTime int64 `json:"updateTime"` /*  分区最近写入消息的毫秒时间戳。  */
	Hw         int64 `json:"hw"`         /*  分区消息高水位线，所有副本均已确认写入的最大偏移量。  */
}

type CtgkafkaGetTopicDetailsReturnObjDataPartitionListReplicasResponse struct {
	BrokerId int32 `json:"brokerId"`           /*  Broker节点ID。  */
	IsLeader *bool `json:"isLeader,omitempty"` /*  是否是主副本。  */
	InSync   *bool `json:"inSync,omitempty"`   /*  副本是否处于同步状态。  */
	Size     int64 `json:"size"`               /*  副本消息大小，单位字节。  */
	Lag      int64 `json:"lag"`                /*  该副本当前落后hw的消息数。  */
}
