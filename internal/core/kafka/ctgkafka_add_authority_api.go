package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaAddAuthorityApi
/* 创建SASL用户权限。
 */type CtgkafkaAddAuthorityApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaAddAuthorityApi(client *core.CtyunClient) *CtgkafkaAddAuthorityApi {
	return &CtgkafkaAddAuthorityApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/saslUser/addAuthority",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaAddAuthorityApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaAddAuthorityRequest) (*CtgkafkaAddAuthorityResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaAddAuthorityRequest
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
	var resp CtgkafkaAddAuthorityResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaAddAuthorityRequest struct {
	RegionId          string   `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId        string   `json:"prodInstId,omitempty"` /*  实例ID。  */
	Username          string   `json:"username,omitempty"`   /*  用户名称。  */
	AddWriteAuthority []string `json:"addWriteAuthority"`    /*  要增加生产权限的主题名称列表。与addReadAuthority至少一个不为空。  */
	AddReadAuthority  []string `json:"addReadAuthority"`     /*  要增加消费权限的主题名称列表。与addWriteAuthority至少一个不为空。  */
}

type CtgkafkaAddAuthorityResponse struct {
	StatusCode string                                 `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                 `json:"message,omitempty"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaAddAuthorityReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                 `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaAddAuthorityReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回创建描述。  */
}
