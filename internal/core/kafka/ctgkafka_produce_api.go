package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaProduceApi
/* 生产消息。
 */type CtgkafkaProduceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaProduceApi(client *core.CtyunClient) *CtgkafkaProduceApi {
	return &CtgkafkaProduceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/produce",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaProduceApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaProduceRequest) (*CtgkafkaProduceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaProduceRequest
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
	var resp CtgkafkaProduceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaProduceRequest struct {
	RegionId    string `json:"regionId,omitempty"`    /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID。  */
	TopicName   string `json:"topicName,omitempty"`   /*  主题名称。  */
	PartitionId int32  `json:"partitionId,omitempty"` /*  指定分区ID。  */
	Key         string `json:"key,omitempty"`         /*  消息key。  */
	Value       string `json:"value,omitempty"`       /*  消息内容。  */
	NumMessages int32  `json:"numMessages,omitempty"` /*  生产消息个数，取值区间[1, 10]，默认值1。  */
}

type CtgkafkaProduceResponse struct {
	StatusCode int32                             `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                            `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaProduceReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                            `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaProduceReturnObjResponse struct {
	Data   string `json:"data,omitempty"`   /*  返回数据。  */
	Result string `json:"result,omitempty"` /*  生产消息成功个数。  */
}
