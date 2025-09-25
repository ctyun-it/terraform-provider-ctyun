package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtgkafkaTopicQueryV3Api
/* 查询主题。
 */type CtgkafkaTopicQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTopicQueryV3Api(client *core.CtyunClient) *CtgkafkaTopicQueryV3Api {
	return &CtgkafkaTopicQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/topic/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaTopicQueryV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTopicQueryV3Request) (*CtgkafkaTopicQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.TopicName != "" {
		ctReq.AddParam("topicName", req.TopicName)
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
	var resp CtgkafkaTopicQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaTopicQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	TopicName  string `json:"topicName,omitempty"`  /*  主题名称，模糊查询。  */
	PageNum    int32  `json:"pageNum,omitempty"`    /*  分页中的页数，默认1，范围1-40000。  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaTopicQueryV3Response struct {
	StatusCode string                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                 `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTopicQueryV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTopicQueryV3ReturnObjResponse struct {
	Total              int32                                        `json:"total,omitempty"`              /*  总记录数。  */
	MaxPartitions      int32                                        `json:"maxPartitions,omitempty"`      /*  分区数量上限。  */
	RemainPartitions   int32                                        `json:"remainPartitions,omitempty"`   /*  剩余分区数量。  */
	TopicMaxPartitions int32                                        `json:"topicMaxPartitions,omitempty"` /*  主题最大分区数量。  */
	Data               []*CtgkafkaTopicQueryV3ReturnObjDataResponse `json:"data"`                         /*  主题信息。  */
}

type CtgkafkaTopicQueryV3ReturnObjDataResponse struct {
	ProdInstId   string                                            `json:"prodInstId,omitempty"`   /*  集群ID。  */
	Name         string                                            `json:"name,omitempty"`         /*  主题名。  */
	PartitionNum int32                                             `json:"partitionNum,omitempty"` /*  主题分区数量。  */
	Factor       int32                                             `json:"factor,omitempty"`       /*  主题副本数。  */
	Configs      *CtgkafkaTopicQueryV3ReturnObjDataConfigsResponse `json:"configs"`                /*  主题配置信息，Map的key为配置名称，value为配置值。  */
}

type CtgkafkaTopicQueryV3ReturnObjDataConfigsResponse struct{}
