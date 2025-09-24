package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaReassignmentsApi
/* 创建分区重平衡任务。
 */type CtgkafkaReassignmentsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaReassignmentsApi(client *core.CtyunClient) *CtgkafkaReassignmentsApi {
	return &CtgkafkaReassignmentsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/reassignments",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaReassignmentsApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaReassignmentsRequest) (*CtgkafkaReassignmentsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaReassignmentsRequest
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
	var resp CtgkafkaReassignmentsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaReassignmentsRequest struct {
	RegionId         string                                        `json:"regionId,omitempty"`     /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId       string                                        `json:"prodInstId,omitempty"`   /*  实例ID，实例节点数需要大于1个。  */
	TopicName        string                                        `json:"topicName,omitempty"`    /*  主题名称<br> 您可以<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7396&data=83&isNormal=1&vid=330">查询主题</a> 获取对应实例下的主题名称列表。  */
	Throttle         int32                                         `json:"throttle,omitempty"`     /*  限流配置，取值-1或范围[2048, 314572800]，-1表示无限制，默认值-1。  */
	RawType          int32                                         `json:"type,omitempty"`         /*  重平衡类型，默认0。<br><li>0：自动<br><li>1：手动  */
	Brokers          []int32                                       `json:"brokers"`                /*  Broker ID列表，type=0时需传入，元素个数取值范围[副本数, 分区数*副本数] <br>您可以<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=330">查询实例列表</a> 对应实例下nodeList属性下的serverSeq获取到brokerId。  */
	PartitionBrokers *CtgkafkaReassignmentsPartitionBrokersRequest `json:"partitionBrokers"`       /*  手动重平衡分区ID与Broker ID列表的对应信息，Map的key表示分区ID，value为一个整型素组，表示将对应分区消息重分配到对应broker下，type=1时需传入。  */
	ScheduleDate     string                                        `json:"scheduleDate,omitempty"` /*  任务调度开始执行时间，格式yyyy-MM-dd HH:mm:ss。不传值或空值表示立即执行。  */
}

type CtgkafkaReassignmentsPartitionBrokersRequest struct{}

type CtgkafkaReassignmentsResponse struct {
	StatusCode int32                                   `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                  `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaReassignmentsReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaReassignmentsReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
