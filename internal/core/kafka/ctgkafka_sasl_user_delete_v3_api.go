package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaSaslUserDeleteV3Api
/* 删除SASL用户。
 */type CtgkafkaSaslUserDeleteV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaSaslUserDeleteV3Api(client *core.CtyunClient) *CtgkafkaSaslUserDeleteV3Api {
	return &CtgkafkaSaslUserDeleteV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/saslUser/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaSaslUserDeleteV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaSaslUserDeleteV3Request) (*CtgkafkaSaslUserDeleteV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaSaslUserDeleteV3Request
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
	var resp CtgkafkaSaslUserDeleteV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaSaslUserDeleteV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	Username   string `json:"username,omitempty"`   /*  用户名称。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
}

type CtgkafkaSaslUserDeleteV3Response struct {
	StatusCode string                                     `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                     `json:"message,omitempty"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaSaslUserDeleteV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                     `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaSaslUserDeleteV3ReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  响应描述。  */
}
