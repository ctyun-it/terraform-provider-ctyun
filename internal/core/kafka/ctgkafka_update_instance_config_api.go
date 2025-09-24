package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUpdateInstanceConfigApi
/* 修改实例配置，包括静态配置和动态配置，动态配置修改后即时生效，静态配置修改后需重启实例生效。
 */type CtgkafkaUpdateInstanceConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUpdateInstanceConfigApi(client *core.CtyunClient) *CtgkafkaUpdateInstanceConfigApi {
	return &CtgkafkaUpdateInstanceConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/updateInstanceConfig",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUpdateInstanceConfigApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUpdateInstanceConfigRequest) (*CtgkafkaUpdateInstanceConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUpdateInstanceConfigRequest
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
	var resp CtgkafkaUpdateInstanceConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUpdateInstanceConfigRequest struct {
	RegionId       string                                               `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId     string                                               `json:"prodInstId,omitempty"` /*  实例ID。  */
	DynamicConfigs []*CtgkafkaUpdateInstanceConfigDynamicConfigsRequest `json:"dynamicConfigs"`       /*  实例动态配置，参数说明可参考<a href="https://www.ctyun.cn/document/10029624/10551853">修改配置参数文档</a>，修改后即时生效。  */
	StaticConfigs  []*CtgkafkaUpdateInstanceConfigStaticConfigsRequest  `json:"staticConfigs"`        /*  实例静态配置，参数说明可参考<a href="https://www.ctyun.cn/document/10029624/10551853">修改配置参数文档</a>，静态配置修改后需重启实例方能生效。  */
}

type CtgkafkaUpdateInstanceConfigDynamicConfigsRequest struct {
	Name  string `json:"name,omitempty"`  /*  需修改配置的名称  */
	Value string `json:"value,omitempty"` /*  需修改配置的值  */
}

type CtgkafkaUpdateInstanceConfigStaticConfigsRequest struct {
	Name  string `json:"name,omitempty"`  /*  需修改配置的名称  */
	Value string `json:"value,omitempty"` /*  需修改配置的值  */
}

type CtgkafkaUpdateInstanceConfigResponse struct {
	StatusCode string `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string `json:"message,omitempty"`    /*  描述状态。  */
	Error      string `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
	ReturnObj  *struct {
		Data string `json:"data"`
	} `json:"returnObj"`
}
