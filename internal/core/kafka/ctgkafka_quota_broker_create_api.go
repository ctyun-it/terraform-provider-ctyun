package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQuotaBrokerCreateApi
/* 创建集群流控配置。
 */type CtgkafkaQuotaBrokerCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQuotaBrokerCreateApi(client *core.CtyunClient) *CtgkafkaQuotaBrokerCreateApi {
	return &CtgkafkaQuotaBrokerCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/quota/broker/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaQuotaBrokerCreateApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQuotaBrokerCreateRequest) (*CtgkafkaQuotaBrokerCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaQuotaBrokerCreateRequest
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
	var resp CtgkafkaQuotaBrokerCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQuotaBrokerCreateRequest struct {
	RegionId         string `json:"regionId,omitempty"`         /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例ID，实例需是Kafka引擎类型且非单机版。  */
	ProducerByteRate int64  `json:"producerByteRate,omitempty"` /*  生产者流控上限，单位字节，单节点流控限制为producerByteRate/节点个数，取值范围[1*节点个数, 实例单节点最大流量*节点个数]，与consumerByteRate不能同时为空。<br>对应实例规格下的单节点流量规格：<li>kafka.2u4g.cluster:100MB<li>kafka.4u8g.cluster:200MB<li>kafka.8u16g.cluster:375MB<li>kafka.12u24g.cluster:625MB<li>kafka.16u32g.cluster:750MB<li>kafka.24u48g.cluster:1125MB<li>kafka.32u64g.cluster:1500MB<li>kafka.48u96g.cluster:2250MB<li>kafka.64u128g.cluster:3000MB <li>kafka.hg.2u4g.cluster:62MB<li>kafka.hg.4u8g.cluster:188MB<li>kafka.hg.8u16g.cluster:563MB<li>kafka.hg.16u32g.cluster:1063MB<li>kafka.hg.32u64g.cluster:2000MB<li>kafka.kp.2u4g.cluster:62MB<li>kafka.kp.4u8g.cluster:188MB<li>kafka.kp.8u16g.cluster:563MB<li>kafka.kp.16u32g.cluster:1063MB<li>kafka.kp.32u64g.cluster:2000MB  */
	ConsumerByteRate int64  `json:"consumerByteRate,omitempty"` /*  消费者流控上限，单位字节，单节点流控限制为consumerByteRate/节点个数，取值范围[1*节点个数, 实例单节点最大流量*节点个数]，与producerByteRate不能同时为空。<br>对应实例规格下的单节点流量规格：<li>kafka.2u4g.cluster:100MB<li>kafka.4u8g.cluster:200MB<li>kafka.8u16g.cluster:375MB<li>kafka.12u24g.cluster:625MB<li>kafka.16u32g.cluster:750MB<li>kafka.24u48g.cluster:1125MB<li>kafka.32u64g.cluster:1500MB<li>kafka.48u96g.cluster:2250MB<li>kafka.64u128g.cluster:3000MB <li>kafka.hg.2u4g.cluster:62MB<li>kafka.hg.4u8g.cluster:188MB<li>kafka.hg.8u16g.cluster:563MB<li>kafka.hg.16u32g.cluster:1063MB<li>kafka.hg.32u64g.cluster:2000MB<li>kafka.kp.2u4g.cluster:62MB<li>kafka.kp.4u8g.cluster:188MB<li>kafka.kp.8u16g.cluster:563MB<li>kafka.kp.16u32g.cluster:1063MB<li>kafka.kp.32u64g.cluster:2000MB  */
}

type CtgkafkaQuotaBrokerCreateResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string `json:"message,omitempty"`    /*  描述状态。  */
}
