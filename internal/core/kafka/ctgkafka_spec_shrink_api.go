package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaSpecShrinkApi
/* 规格缩容。
 */type CtgkafkaSpecShrinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaSpecShrinkApi(client *core.CtyunClient) *CtgkafkaSpecShrinkApi {
	return &CtgkafkaSpecShrinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/specShrink",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaSpecShrinkApi) Do(ctx context.Context, credential core.Credential, req *CtgkafkaSpecShrinkRequest) (*CtgkafkaSpecShrinkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaSpecShrinkRequest
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
	var resp CtgkafkaSpecShrinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaSpecShrinkRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	SpecName   string `json:"specName,omitempty"`   /*  实例的规格类型。<br>计算增强型的规格可选为：<li>kafka.2u4g.cluster<li>kafka.4u8g.cluster<li>kafka.8u16g.cluster<li>kafka.12u24g.cluster<li>kafka.16u32g.cluster<li>kafka.24u48g.cluster<li>kafka.32u64g.cluster<li>kafka.48u96g.cluster<li>kafka.64u128g.cluster <br>海光-计算增强型的规格可选为：<li>kafka.hg.2u4g.cluster<li>kafka.hg.4u8g.cluster<li>kafka.hg.8u16g.cluster<li>kafka.hg.16u32g.cluster<li>kafka.hg.32u64g.cluster <br>鲲鹏-计算增强型的规格可选为：<li>kafka.kp.2u4g.cluster<li>kafka.kp.4u8g.cluster<li>kafka.kp.8u16g.cluster<li>kafka.kp.16u32g.cluster<li>kafka.kp.32u64g.cluster  */
}

type CtgkafkaSpecShrinkResponse struct {
	StatusCode string                               `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                               `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaSpecShrinkReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                               `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaSpecShrinkReturnObjResponse struct {
}
