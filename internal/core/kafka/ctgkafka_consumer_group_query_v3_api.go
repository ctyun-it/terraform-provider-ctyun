package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaConsumerGroupQueryV3Api
/* 查询消费组。
 */type CtgkafkaConsumerGroupQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaConsumerGroupQueryV3Api(client *core.CtyunClient) *CtgkafkaConsumerGroupQueryV3Api {
	return &CtgkafkaConsumerGroupQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/consumerGroup/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaConsumerGroupQueryV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaConsumerGroupQueryV3Request) (*CtgkafkaConsumerGroupQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	if req.GroupName != "" {
		ctReq.AddParam("groupName", req.GroupName)
	}
	if req.PageNum != "" {
		ctReq.AddParam("pageNum", req.PageNum)
	}
	if req.PageSize != "" {
		ctReq.AddParam("pageSize", req.PageSize)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaConsumerGroupQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaConsumerGroupQueryV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	GroupName  string `json:"groupName,omitempty"`  /*  消费组名称，模糊查询。  */
	PageNum    string `json:"pageNum,omitempty"`    /*  分页中的页数，默认1，范围1-40000。  */
	PageSize   string `json:"pageSize,omitempty"`   /*  分页中的每页大小，默认10，范围1-40000。  */
}

type CtgkafkaConsumerGroupQueryV3Response struct {
	StatusCode string                                         `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                         `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaConsumerGroupQueryV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaConsumerGroupQueryV3ReturnObjResponse struct {
	Data  []*CtgkafkaConsumerGroupQueryV3ReturnObjDataResponse `json:"data"`            /*  返回数据。  */
	Total int32                                                `json:"total,omitempty"` /*  总记录数。  */
}

type CtgkafkaConsumerGroupQueryV3ReturnObjDataResponse struct {
	Id          int32  `json:"id,omitempty"`          /*  消费组ID。  */
	Name        string `json:"name,omitempty"`        /*  消费组名。  */
	Description string `json:"description,omitempty"` /*  消费组描述。  */
	Ctime       string `json:"ctime,omitempty"`       /*  创建时间。  */
}
