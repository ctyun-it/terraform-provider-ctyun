package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaUpdateUserTopicsAclApi
/* 更新用户ACL权限。
 */type CtgkafkaUpdateUserTopicsAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaUpdateUserTopicsAclApi(client *core.CtyunClient) *CtgkafkaUpdateUserTopicsAclApi {
	return &CtgkafkaUpdateUserTopicsAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/saslUser/updateUserTopicsAcl",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaUpdateUserTopicsAclApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaUpdateUserTopicsAclRequest) (*CtgkafkaUpdateUserTopicsAclResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaUpdateUserTopicsAclRequest
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
	var resp CtgkafkaUpdateUserTopicsAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaUpdateUserTopicsAclRequest struct {
	RegionId             string                                                    `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId           string                                                    `json:"prodInstId,omitempty"` /*  实例ID。  */
	AclOperationInfoList []*CtgkafkaUpdateUserTopicsAclAclOperationInfoListRequest `json:"aclOperationInfoList"` /*  要操作的ACL信息列表  */
}

type CtgkafkaUpdateUserTopicsAclAclOperationInfoListRequest struct {
	EventType  string `json:"eventType,omitempty"`  /*  操作事件类型<br><li>CREATE:创建<br><li>DELETE:删除<br>  */
	Permission string `json:"permission,omitempty"` /*  权限<br><li>ALLOW:允许<br><li>DENY:拒绝<br> 默认：ALLOW  */
	UserName   string `json:"userName,omitempty"`   /*  用户名  */
	Ip         string `json:"ip,omitempty"`         /*  ip或网段，* 表示所有ip  */
	Operation  string `json:"operation,omitempty"`  /*  操作<br><li>READ:消费<br><li>WRITE:生产  */
	Topic      string `json:"topic,omitempty"`      /*  topic名称  */
}

type CtgkafkaUpdateUserTopicsAclResponse struct {
	StatusCode string                                  `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                  `json:"message"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaUserTopicsAclReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                  `json:"error"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaUserTopicsAclReturnObjResponse struct {
	Data string `json:"data"` /*  返回数据。  */
}
