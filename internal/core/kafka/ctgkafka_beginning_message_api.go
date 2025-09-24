package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaBeginningMessageApi
/* 查询分区最早消息位置。
 */type CtgkafkaBeginningMessageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaBeginningMessageApi(client *core.CtyunClient) *CtgkafkaBeginningMessageApi {
	return &CtgkafkaBeginningMessageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/beginningMessage",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaBeginningMessageApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaBeginningMessageRequest) (*CtgkafkaBeginningMessageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("partitionId", strconv.FormatInt(int64(req.PartitionId), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaBeginningMessageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaBeginningMessageRequest struct {
	RegionId    string `json:"regionId,omitempty"`    /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID。  */
	TopicName   string `json:"topicName,omitempty"`   /*  主题名称。  */
	PartitionId int32  `json:"partitionId,omitempty"` /*  主题分区ID。  */
}

type CtgkafkaBeginningMessageResponse struct {
	StatusCode int32                                      `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                     `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaBeginningMessageReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaBeginningMessageReturnObjResponse struct {
	Data []*CtgkafkaBeginningMessageReturnObjDataResponse `json:"data"` /*  主题信息。  */
}

type CtgkafkaBeginningMessageReturnObjDataResponse struct {
	TopicName string `json:"topicName,omitempty"` /*  主题名称。  */
	Partition int32  `json:"partition,omitempty"` /*  分区ID。  */
	Offset    int64  `json:"offset,omitempty"`    /*  最新消息偏移量。  */
	TimeStamp int64  `json:"timeStamp,omitempty"` /*  最新消息时间戳。  */
}
