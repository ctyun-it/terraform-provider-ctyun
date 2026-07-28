package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaAclStrategyCreateApi
/* 创建ACL策略。
 */type CtgkafkaAclStrategyCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaAclStrategyCreateApi(client *core.CtyunClient) *CtgkafkaAclStrategyCreateApi {
	return &CtgkafkaAclStrategyCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/kafkaAclStrategy/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaAclStrategyCreateApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaAclStrategyCreateRequest) (*CtgkafkaAclStrategyCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaAclStrategyCreateRequest
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
	var resp CtgkafkaAclStrategyCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaAclStrategyCreateRequest struct {
	RegionId    string                                   `json:"regionId,omitempty"`    /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId  string                                   `json:"prodInstId,omitempty"`  /*  实例ID。  */
	Name        string                                   `json:"name,omitempty"`        /*  策略名称，规则如下：<br><li>以英文字母、数字、下划线开头，且只能由英文字母、数字、句点、中划线、下划线组成。<br><li>长度3-64。<br><li>名称不可重复。  */
	Rules       []*CtgkafkaAclStrategyCreateRulesRequest `json:"rules"`                 /*  acl规则  */
	UseNewTopic string                                   `json:"useNewTopic,omitempty"` /*  是否应用到新增主题 <br><li>1：是<br><li>2：否<br><li>默认值：2。  */
}

type CtgkafkaAclStrategyCreateRulesRequest struct {
	Permission string `json:"permission,omitempty"` /*  权限<br><li>ALLOW:允许<br><li>DENY:拒绝<br>  */
	UserName   string `json:"userName,omitempty"`   /*  用户名，必须是已经集群中创建了的用户  */
	Ip         string `json:"ip,omitempty"`         /*  ip或网段，多个用半角分号分开，默认*，表示所有ip  */
	Operation  string `json:"operation,omitempty"`  /*  操作<br><li>READ:消费<br><li>WRITE:生产  */
}

type CtgkafkaAclStrategyCreateResponse struct {
	StatusCode string `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string `json:"message"`    /*  提示信息。  */
	ReturnObj  string `json:"returnObj"`  /*  返回对象。  */
	Error      string `json:"error"`      /*  错误码，描述错误信息。  */
}
