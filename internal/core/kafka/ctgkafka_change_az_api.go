package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaChangeAzApi
/* 变更实例节点所在的可用区。
 */type CtgkafkaChangeAzApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaChangeAzApi(client *core.CtyunClient) *CtgkafkaChangeAzApi {
	return &CtgkafkaChangeAzApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/changeAz",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaChangeAzApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaChangeAzRequest) (*CtgkafkaChangeAzResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaChangeAzRequest
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
	var resp CtgkafkaChangeAzResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaChangeAzRequest struct {
	RegionId   string   `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string   `json:"prodInstId,omitempty"` /*  实例ID。  */
	AzInfoList []string `json:"azInfoList"`           /*  可用区信息。只能填一个（单可用区）或三个（多可用区），可用区信息可调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87&isNormal=1&vid=81">资源池可用区查询</a>API接口查询。  */
}

type CtgkafkaChangeAzResponse struct {
	StatusCode string                             `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                             `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaChangeAzReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                             `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaChangeAzReturnObjResponse struct {
	Data string `json:"data,omitempty"` /*  返回数据。  */
}
