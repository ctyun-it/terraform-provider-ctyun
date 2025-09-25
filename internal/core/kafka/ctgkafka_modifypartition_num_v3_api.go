package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaModifypartitionNumV3Api
/* 修改主题分区数。
 */type CtgkafkaModifypartitionNumV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaModifypartitionNumV3Api(client *core.CtyunClient) *CtgkafkaModifypartitionNumV3Api {
	return &CtgkafkaModifypartitionNumV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/topic/modifyPartitionNum",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaModifypartitionNumV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaModifypartitionNumV3Request) (*CtgkafkaModifypartitionNumV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaModifypartitionNumV3Request
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
	var resp CtgkafkaModifypartitionNumV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaModifypartitionNumV3Request struct {
	RegionId     string `json:"regionId,omitempty"`     /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例ID。  */
	TopicName    string `json:"topicName,omitempty"`    /*  主题名称，<br> 您可以<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7396&data=83&isNormal=1&vid=330">查询主题</a> 获取对应实例下的主题名称列表。  */
	PartitionNum int32  `json:"partitionNum,omitempty"` /*  新分区数量，取值范围[旧分区数量+1，min(100, 实例剩余分区数量)]。实例剩余分区数量=实例分区上限-所有主题分区数之和。<br>您可以<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=330">查询实例</a>查询对应实例下partitionNum获取分区数上限。<br>您可以<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7396&data=83&isNormal=1&vid=330">查询主题</a>查询对应实例下主题的分区数量。  */
}

type CtgkafkaModifypartitionNumV3Response struct {
	StatusCode string                                         `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                         `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaModifypartitionNumV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                         `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaModifypartitionNumV3ReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
