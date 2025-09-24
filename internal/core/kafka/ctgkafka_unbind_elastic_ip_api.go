package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUnbindElasticIpApi
/* 解绑弹性IP。
 */type CtgkafkaUnbindElasticIpApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUnbindElasticIpApi(client *core.CtyunClient) *CtgkafkaUnbindElasticIpApi {
	return &CtgkafkaUnbindElasticIpApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/unbindElasticIp",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUnbindElasticIpApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUnbindElasticIpRequest) (*CtgkafkaUnbindElasticIpResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUnbindElasticIpRequest
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
	var resp CtgkafkaUnbindElasticIpResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUnbindElasticIpRequest struct {
	RegionId       string `json:"regionId,omitempty"`       /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	PaasInstanceId string `json:"paasInstanceId,omitempty"` /*  实例ID。  */
	ElasticIp      string `json:"elasticIp,omitempty"`      /*  解绑的弹性IP。  */
}

type CtgkafkaUnbindElasticIpResponse struct {
	StatusCode string                                    `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                    `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaUnbindElasticIpReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUnbindElasticIpReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回提示。  */
}
