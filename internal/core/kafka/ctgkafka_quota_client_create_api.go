package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaQuotaClientCreateApi
/* 创建用户/客户端流控配置。
 */type CtgkafkaQuotaClientCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaQuotaClientCreateApi(client *core.CtyunClient) *CtgkafkaQuotaClientCreateApi {
	return &CtgkafkaQuotaClientCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/quota/user-client/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaQuotaClientCreateApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaQuotaClientCreateRequest) (*CtgkafkaQuotaClientCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaQuotaClientCreateRequest
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
	var resp CtgkafkaQuotaClientCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaQuotaClientCreateRequest struct {
	RegionId         string `json:"regionId,omitempty"`         /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId       string `json:"prodInstId,omitempty"`       /*  实例ID，实例需是Kafka引擎类型且非单机版。  */
	User             string `json:"user,omitempty"`             /*  用户名，defaultUser=false时有效，只能由字母，数字，中划线，下划线组成，长度为3~64个字符。  */
	Client           string `json:"client,omitempty"`           /*  客户端ID，defaultClient=false时有效，只能由英文字母开头，且只能由英文字母、数字、中划线、下划线组成，长度4~64个字符。  */
	DefaultUser      *bool  `json:"defaultUser"`                /*  是否使用默认用户，true表示对所有用户限流，默认false。  */
	DefaultClient    *bool  `json:"defaultClient"`              /*  是否使用默认客户端，true表示对所有客户端限流，默认false。  */
	ProducerByteRate int64  `json:"producerByteRate,omitempty"` /*  生产者流控上限，单位字节，取值范围[1048576, 1073741824]，即1MB至1GB，与consumerByteRate至少一个不为空。  */
	ConsumerByteRate int64  `json:"consumerByteRate,omitempty"` /*  消费者流控上限，单位字节，取值范围[1048576, 1073741824]，即1MB至1GB，与producerByteRate至少一个不为空。  */
}

type CtgkafkaQuotaClientCreateResponse struct {
	StatusCode int32  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功：800，失败：900。  */
	Message    string `json:"message,omitempty"`    /*  描述状态。  */
}
