package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaBindElasticIpApi
/* 绑定弹性IP。
 */type CtgkafkaBindElasticIpApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaBindElasticIpApi(client *core.CtyunClient) *CtgkafkaBindElasticIpApi {
	return &CtgkafkaBindElasticIpApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/bindElasticIp",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaBindElasticIpApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaBindElasticIpRequest) (*CtgkafkaBindElasticIpResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaBindElasticIpRequest
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
	var resp CtgkafkaBindElasticIpResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaBindElasticIpRequest struct {
	RegionId       string `json:"regionId,omitempty"`       /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	PaasInstanceId string `json:"paasInstanceId,omitempty"` /*  实例ID。  */
	Ip             string `json:"ip,omitempty"`             /*  节点的IP地址，可以通过查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=48&api=7382&data=83&isNormal=1&vid=330">查询实例</a> 获取具体节点IP地址。  */
	ElasticIp      string `json:"elasticIp,omitempty"`      /*  绑定的弹性IP。  */
}

type CtgkafkaBindElasticIpResponse struct {
	StatusCode string                                  `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                  `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaBindElasticIpReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                  `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaBindElasticIpReturnObjResponse struct {
	Data *CtgkafkaBindElasticIpReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type CtgkafkaBindElasticIpReturnObjDataResponse struct{}
