package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaGetAuthorityApi
/* 查询SASL用户权限。
 */type CtgkafkaGetAuthorityApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaGetAuthorityApi(client *core.CtyunClient) *CtgkafkaGetAuthorityApi {
	return &CtgkafkaGetAuthorityApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/saslUser/getAuthority",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaGetAuthorityApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaGetAuthorityRequest) (*CtgkafkaGetAuthorityResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("username", req.Username)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaGetAuthorityResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaGetAuthorityRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	Username   string `json:"username,omitempty"`   /*  用户名称。  */
}

type CtgkafkaGetAuthorityResponse struct {
	StatusCode string                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                 `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CtgkafkaGetAuthorityReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaGetAuthorityReturnObjResponse struct {
	Data *CtgkafkaGetAuthorityReturnObjDataResponse `json:"data"` /*  用户权限数据  */
}

type CtgkafkaGetAuthorityReturnObjDataResponse struct {
	Username             string   `json:"username,omitempty"`   /*  用户名  */
	ReadAuthorityTopics  []string `json:"readAuthorityTopics"`  /*  拥有消费权限的主题列表  */
	WriteAuthorityTopics []string `json:"writeAuthorityTopics"` /*  拥有生产权限的主题列表  */
}
