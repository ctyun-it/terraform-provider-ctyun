package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaMonitorMetricsV3Api
/* 查询监控指标。
 */type CtgkafkaMonitorMetricsV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaMonitorMetricsV3Api(client *core.CtyunClient) *CtgkafkaMonitorMetricsV3Api {
	return &CtgkafkaMonitorMetricsV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/metrics",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaMonitorMetricsV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaMonitorMetricsV3Request) (*CtgkafkaMonitorMetricsV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("regionCode", req.RegionCode)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("metricName", req.MetricName)
	ctReq.AddParam("startTime", strconv.FormatInt(int64(req.StartTime), 10))
	ctReq.AddParam("endTime", strconv.FormatInt(int64(req.EndTime), 10))
	if req.Labels != "" {
		ctReq.AddParam("labels", req.Labels)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaMonitorMetricsV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaMonitorMetricsV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	RegionCode string `json:"regionCode,omitempty"` /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	MetricName string `json:"metricName,omitempty"` /*  指标名称，例如节点存活数指标，对应的metricName为current_brokers <br>您可以<a href="https://www.ctyun.cn/document/10029624/10991364">查询可用的指标名称</a>，参考该文档表格的指标ID。  */
	StartTime  int32  `json:"startTime,omitempty"`  /*  查询数据起始时间，UNIX时间戳，单位秒（s）。  */
	EndTime    int32  `json:"endTime,omitempty"`    /*  查询数据截止时间，UNIX时间戳，单位秒（s）。startTime必须小于endTime 。  */
	Labels     string `json:"labels,omitempty"`     /*  指标的维度标签，格式为：key1=value1,key2=value2，<br><li>节点维度监控指标，lebels可传vpcIp，如lebels参数为："vpcIp=192.168.0.46"。<br><li>主题维度监控指标lebels可传topic，如lebels参数为："topic=mytopic"。<br><li>消费组维度监控指标，lebels可传group、topic，如lebels参数为："group=aaaa,topic=topic-w"。  */
}

type CtgkafkaMonitorMetricsV3Response struct {
	StatusCode string                                       `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                       `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  []*CtgkafkaMonitorMetricsV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaMonitorMetricsV3ReturnObjResponse struct {
	Metric *CtgkafkaMonitorMetricsV3ReturnObjMetricResponse   `json:"metric"` /*  metrics标签，标签如下：  <li>instance_id：实例ID<li>region_name：资源池名称 <li>vpc_ip：节点IP <li>topic：主题名 <li>group：消费组名称  */
	Values []*CtgkafkaMonitorMetricsV3ReturnObjValuesResponse `json:"values"` /*  指标数据列表。数据为一个数组，数组索引为0的数据时间戳，单位为秒；数组索引为1的数据为指标值。  */
}

type CtgkafkaMonitorMetricsV3ReturnObjMetricResponse struct{}

type CtgkafkaMonitorMetricsV3ReturnObjValuesResponse struct {
	MetricValues []string `json:"metricValues"` /*  标指标数据，数组第一个数据为时时间戳，数组第二个数据为指标值。  */
}
