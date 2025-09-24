package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaQueryMessageByOffsetV3Api
/* 按点位offset查询消息。
 */type CtgkafkaQueryMessageByOffsetV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQueryMessageByOffsetV3Api(client *core.CtyunClient) *CtgkafkaQueryMessageByOffsetV3Api {
	return &CtgkafkaQueryMessageByOffsetV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/queryMessageByOffset",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaQueryMessageByOffsetV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQueryMessageByOffsetV3Request) (*CtgkafkaQueryMessageByOffsetV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("topicName", req.TopicName)
	ctReq.AddParam("offset", strconv.FormatInt(int64(req.Offset), 10))
	ctReq.AddParam("partition", strconv.FormatInt(int64(req.Partition), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaQueryMessageByOffsetV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQueryMessageByOffsetV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	TopicName  string `json:"topicName,omitempty"`  /*  topic名称。  */
	Offset     int64  `json:"offset,omitempty"`     /*  offset位置，大于或等于0。  */
	Partition  int32  `json:"partition,omitempty"`  /*  分区号，大于或等于0。  */
}

type CtgkafkaQueryMessageByOffsetV3Response struct {
	StatusCode string                                           `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                           `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaQueryMessageByOffsetV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                           `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaQueryMessageByOffsetV3ReturnObjResponse struct {
	Data []*CtgkafkaQueryMessageByOffsetV3ReturnObjDataResponse `json:"data"` /*  消息数据列表。  */
}

type CtgkafkaQueryMessageByOffsetV3ReturnObjDataResponse struct {
	Topic     string `json:"topic,omitempty"`     /*  主题名称。  */
	Partition int32  `json:"partition,omitempty"` /*  分区号。  */
	Offset    int64  `json:"offset,omitempty"`    /*  offset位置。  */
	Timestamp int64  `json:"timestamp,omitempty"` /*  生产消息的时间戳，单位毫秒。  */
	Key       string `json:"key,omitempty"`       /*  消息key。  */
	Value     string `json:"value,omitempty"`     /*  消息内容。  */
	Size      string `json:"size,omitempty"`      /*  消息大小，单位字节。  */
}
