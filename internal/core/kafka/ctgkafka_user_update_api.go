package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUserUpdateApi
/* 更新Kafka实例的SASL用户，SASL用户可连接开启SASL的Kafka实例。
 */type CtgkafkaUserUpdateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUserUpdateApi(client *core.CtyunClient) *CtgkafkaUserUpdateApi {
	return &CtgkafkaUserUpdateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/saslUser/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUserUpdateApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUserUpdateRequest) (*CtgkafkaUserUpdateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUserUpdateRequest
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
	var resp CtgkafkaUserUpdateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUserUpdateRequest struct {
	RegionId    string `json:"regionId,omitempty"`    /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId  string `json:"prodInstId,omitempty"`  /*  实例ID。  */
	Username    string `json:"username,omitempty"`    /*  用户名称。  */
	NewPassword string `json:"newPassword,omitempty"` /*  密码，规则如下：<br><li>长度8-26字符。<br><li>必须同时包含大写字母、小写字母、数字和英文格式特殊符号(@%^*_+!$-=.)中的至少三种类型。<br><li>不能有空格。  */
	Description string `json:"description,omitempty"` /*  用户描述，规则如下：<br><li>不能以+,-,@,= 特殊字符开头。 <br><li>长度不能大于200。  */
}

type CtgkafkaUserUpdateResponse struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                               `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CtgkafkaUserUpdateReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUserUpdateReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回创建描述。  */
}
