package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaTopicStatusV3Api
/* 查询主题状态。
 */type CtgkafkaTopicStatusV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTopicStatusV3Api(client *core.CtyunClient) *CtgkafkaTopicStatusV3Api {
	return &CtgkafkaTopicStatusV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/status",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaTopicStatusV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTopicStatusV3Request) (*CtgkafkaTopicStatusV3Response, error) {
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
	var resp CtgkafkaTopicStatusV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaTopicStatusV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	TopicName  string `json:"topicName,omitempty"`  /*  主题名称。  */
}

type CtgkafkaTopicStatusV3Response struct {
	StatusCode string                                  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                  `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTopicStatusV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTopicStatusV3ReturnObjResponse struct {
	Total int32                                         `json:"total,omitempty"` /*  消息总数。  */
	Data  []*CtgkafkaTopicStatusV3ReturnObjDataResponse `json:"data"`            /*  主题分区状态信息。  */
}

type CtgkafkaTopicStatusV3ReturnObjDataResponse struct {
	PartitionId int32 `json:"partitionId,omitempty"` /*  分区ID。  */
	Total       int64 `json:"total,omitempty"`       /*  分区消息总数。  */
	Begin       int64 `json:"begin,omitempty"`       /*  分区消息最小偏移量。  */
	End         int64 `json:"end,omitempty"`         /*  分区消息最大偏移量。  */
}
