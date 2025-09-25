package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaEnableAutoReassignApi
/* 节点扩容自动分区重平衡。
 */type CtgkafkaEnableAutoReassignApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaEnableAutoReassignApi(client *core.CtyunClient) *CtgkafkaEnableAutoReassignApi {
	return &CtgkafkaEnableAutoReassignApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/clusterExtend/enableAutoReassign",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaEnableAutoReassignApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaEnableAutoReassignRequest) (*CtgkafkaEnableAutoReassignResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaEnableAutoReassignRequest
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
	var resp CtgkafkaEnableAutoReassignResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaEnableAutoReassignRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	Enable     int32  `json:"enable,omitempty"`     /*  是否开启自动分区重平衡。<br><li>1：开启<br><li>2：关闭  */
}

type CtgkafkaEnableAutoReassignResponse struct {
	StatusCode string                                       `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                       `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaEnableAutoReassignReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaEnableAutoReassignReturnObjResponse struct {
	Data *CtgkafkaEnableAutoReassignReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type CtgkafkaEnableAutoReassignReturnObjDataResponse struct{}
