package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpSpecShrinkApi
/* 规格缩容。
 */type AmqpSpecShrinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpSpecShrinkApi(client *core.CtyunClient) *AmqpSpecShrinkApi {
	return &AmqpSpecShrinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/specShrink",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpSpecShrinkApi) Do(ctx context.Context, credential core.Credential, req *AmqpSpecShrinkRequest) (*AmqpSpecShrinkResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpSpecShrinkRequest
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
	var resp AmqpSpecShrinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpSpecShrinkRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	SpecName   string `json:"specName,omitempty"`   /*  实例的规格类型。资源池所具备的规格可通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=55&api=20202&data=39&isNormal=1&vid=38">查询产品规格</a>接口获取。<br>计算增强型的规格可选为：<br>- rabbitmq.2u4g.cluster<br>- rabbitmq.4u8g.cluster<br>- rabbitmq.8u16g.cluster<br>- rabbitmq.12u24g.cluster<br>- rabbitmq.16u32g.cluster<br>- rabbitmq.24u48g.cluster<br>- rabbitmq.32u64g.cluster<br>- rabbitmq.48u96g.cluster<br>海光-计算增强型的规格可选为：<br>- rabbitmq.hg.2u4g.cluster<br>- rabbitmq.hg.4u8g.cluster<br>- rabbitmq.hg.8u16g.cluster<br>- rabbitmq.hg.16u32g.cluster <br>鲲鹏-计算增强型的规格可选为：<br>- rabbitmq.kp.2u4g.cluster<br>- rabbitmq.kp.4u8g.cluster<br>- rabbitmq.kp.8u16g.cluster<br>- rabbitmq.kp.16u32g.cluster  */
}

type AmqpSpecShrinkResponse struct {
	StatusCode string                           `json:"statusCode"` /*  响应状态码。<br>- 800：成功。<br>- 900：失败。  */
	Message    string                           `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpSpecShrinkReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                           `json:"error"`      /*  错误码，只有失败才显示，参见错误码说明。  */
}

type AmqpSpecShrinkReturnObjResponse struct {
	Data *AmqpSpecShrinkReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type AmqpSpecShrinkReturnObjDataResponse struct{}
