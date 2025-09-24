package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaQueryMessageApi
/* 指定时间范围查询消息。
 */type CtgkafkaQueryMessageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQueryMessageApi(client *core.CtyunClient) *CtgkafkaQueryMessageApi {
	return &CtgkafkaQueryMessageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/message",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaQueryMessageApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQueryMessageRequest) (*CtgkafkaQueryMessageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("partitionId", strconv.FormatInt(int64(req.PartitionId), 10))
	ctReq.AddParam("startTime", strconv.FormatInt(int64(req.StartTime), 10))
	ctReq.AddParam("endTime", strconv.FormatInt(int64(req.EndTime), 10))
	if req.Content != "" {
		ctReq.AddParam("content", req.Content)
	}
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaQueryMessageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQueryMessageRequest struct {
	RegionId    string `json:"regionId,omitempty"`    /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID。  */
	TopicName   string `json:"topicName,omitempty"`   /*  topic名称。  */
	PartitionId int32  `json:"partitionId,omitempty"` /*  分区号，大于或等于0。  */
	StartTime   int64  `json:"startTime,omitempty"`   /*  时间戳，单位毫秒，endTime必须大于startTime。  */
	EndTime     int64  `json:"endTime,omitempty"`     /*  时间戳，单位毫秒，endTime必须大于startTime。  */
	Content     string `json:"content,omitempty"`     /*  消息包含内容  */
	PageNum     int32  `json:"pageNum,omitempty"`     /*  分页中的页数，默认1，范围1-40000。  */
	PageSize    int32  `json:"pageSize,omitempty"`    /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaQueryMessageResponse struct {
	StatusCode string                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                 `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaQueryMessageReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaQueryMessageReturnObjResponse struct {
	Data *CtgkafkaQueryMessageReturnObjDataResponse `json:"data"` /*  消息分页数据对象。  */
}

type CtgkafkaQueryMessageReturnObjDataResponse struct {
	Messages []*CtgkafkaQueryMessageReturnObjDataMessagesResponse `json:"messages"` /*  消息信息。  */
}

type CtgkafkaQueryMessageReturnObjDataMessagesResponse struct {
	Topic     string `json:"topic,omitempty"`     /*  主题名称。  */
	Partition int32  `json:"partition,omitempty"` /*  分区号。  */
	Offset    int64  `json:"offset,omitempty"`    /*  offset位置。  */
	Timestamp int64  `json:"timestamp,omitempty"` /*  生产消息的时间戳，单位毫秒。  */
	Key       string `json:"key,omitempty"`       /*  消息Key。  */
	Value     string `json:"value,omitempty"`     /*  消息内容。  */
	Size      int32  `json:"size,omitempty"`      /*  消息大小。  */
}
