package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaAclQueryV3Api
/* 查询SASL用户下资源ACL列表，该接口为旧接口，推荐使用查询SASL用户权限接口。
 */type CtgkafkaAclQueryV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaAclQueryV3Api(client *core.CtyunClient) *CtgkafkaAclQueryV3Api {
	return &CtgkafkaAclQueryV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v3/acl/query",
			ContentType:  "application/x-www-form-urlencoded",
		},
	}
}

func (a *CtgkafkaAclQueryV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaAclQueryV3Request) (*CtgkafkaAclQueryV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("prodInstId", req.ProdInstId)
	ctReq.AddParam("username", req.Username)
	ctReq.AddParam("resourceType", req.ResourceType)
	ctReq.AddParam("resourceName", req.ResourceName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtgkafkaAclQueryV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaAclQueryV3Request struct {
	RegionId     string `json:"regionId,omitempty"`     /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId   string `json:"prodInstId,omitempty"`   /*  实例ID。  */
	Username     string `json:"username,omitempty"`     /*  用户名称。  */
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型，可选值：<br><li>TOPIC：主题资源<br><li>GROUP：消费组资源，云原生引擎无GROUP资源  */
	ResourceName string `json:"resourceName,omitempty"` /*  资源名称。  */
}

type CtgkafkaAclQueryV3Response struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                               `json:"message,omitempty"`    /*  提示信息。  */
	ReturnObj  *CtgkafkaAclQueryV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaAclQueryV3ReturnObjResponse struct {
	Data []*CtgkafkaAclQueryV3ReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type CtgkafkaAclQueryV3ReturnObjDataResponse struct {
	Resource *CtgkafkaAclQueryV3ReturnObjDataResourceResponse `json:"resource"` /*  资源信息。  */
	Entry    *CtgkafkaAclQueryV3ReturnObjDataEntryResponse    `json:"entry"`    /*  资源条目信息。  */
}

type CtgkafkaAclQueryV3ReturnObjDataResourceResponse struct {
	ResourceType string `json:"resourceType,omitempty"` /*  资源类型。  */
	Name         string `json:"name,omitempty"`         /*  资源名称。  */
}

type CtgkafkaAclQueryV3ReturnObjDataEntryResponse struct {
	Data *CtgkafkaAclQueryV3ReturnObjDataEntryDataResponse `json:"data"` /*  资源ACL条目。  */
}

type CtgkafkaAclQueryV3ReturnObjDataEntryDataResponse struct {
	Principal string `json:"principal,omitempty"` /*  权限主体。  */
	Host      string `json:"host,omitempty"`      /*  拥有权限的主机地址。  */
}
