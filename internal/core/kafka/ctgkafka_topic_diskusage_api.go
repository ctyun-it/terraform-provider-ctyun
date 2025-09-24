package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaTopicDiskusageApi
/* 查询主题磁盘占用情况。
 */type CtgkafkaTopicDiskusageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTopicDiskusageApi(client *core.CtyunClient) *CtgkafkaTopicDiskusageApi {
	return &CtgkafkaTopicDiskusageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/diskusage",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaTopicDiskusageApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTopicDiskusageRequest) (*CtgkafkaTopicDiskusageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.MinSize != 0 {
		ctReq.AddParam("minSize", strconv.FormatInt(int64(req.MinSize), 10))
	}
	if req.Percentage != 0 {
		ctReq.AddParam("percentage", strconv.FormatInt(int64(req.Percentage), 10))
	}
	if req.Top != 0 {
		ctReq.AddParam("top", strconv.FormatInt(int64(req.Top), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaTopicDiskusageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaTopicDiskusageRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	MinSize    int32  `json:"minSize,omitempty"`    /*  查询每个Broker节点磁盘占用量超过minSize的分区，单位GB，最小值1。  */
	Percentage int32  `json:"percentage,omitempty"` /*  查询每个Broker节点磁盘占用量占比超过实例最大磁盘上限百分之percentage的分区，取值范围[1, 100]，minSize为空时有效。  */
	Top        int32  `json:"top,omitempty"`        /*  查询每个Broker节点磁盘占用量最多的前多少个分区，最小值1，默认值10。minSize和percent为空时有效。  */
}

type CtgkafkaTopicDiskusageResponse struct {
	StatusCode int32                                    `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string                                   `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTopicDiskusageReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                   `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTopicDiskusageReturnObjResponse struct {
	Data []*CtgkafkaTopicDiskusageReturnObjDataResponse `json:"data"` /*  每个Broker节点主题分区磁盘使用量信息。  */
}

type CtgkafkaTopicDiskusageReturnObjDataResponse struct {
	Total      int64                                                   `json:"total,omitempty"`      /*  磁盘总大小，单位GB。  */
	Used       int64                                                   `json:"used,omitempty"`       /*  磁盘已使用大小，单位GB。  */
	Free       int64                                                   `json:"free,omitempty"`       /*  磁盘剩余大小，单位GB。  */
	Rate       float64                                                 `json:"rate"`                 /*  磁盘已使用比例。  */
	BrokerName string                                                  `json:"brokerName,omitempty"` /*  Broker节点名。  */
	TopicList  []*CtgkafkaTopicDiskusageReturnObjDataTopicListResponse `json:"topicList"`            /*  主题分区占用磁盘量信息。  */
}

type CtgkafkaTopicDiskusageReturnObjDataTopicListResponse struct {
	Topic        string  `json:"topic,omitempty"`       /*  主题名称。  */
	PartitionId  int32   `json:"partitionId,omitempty"` /*  主题分区ID。  */
	Size         int64   `json:"size,omitempty"`        /*  主题分区占用磁盘大小，单位字节。  */
	BrokerName   string  `json:"brokerName,omitempty"`  /*  Broker节点名称。  */
	DiskUsedRate float64 `json:"diskUsedRate"`          /*  主题分区占用磁盘比例。  */
}
