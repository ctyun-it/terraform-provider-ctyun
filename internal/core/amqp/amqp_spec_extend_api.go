package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpSpecExtendApi
/* 规格扩容。
 */type AmqpSpecExtendApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpSpecExtendApi(client *core.CtyunClient) *AmqpSpecExtendApi {
	return &AmqpSpecExtendApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/specExtend",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpSpecExtendApi) Do(ctx context.Context, credential core.Credential, req *AmqpSpecExtendRequest) (*AmqpSpecExtendResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpSpecExtendRequest
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
	var resp AmqpSpecExtendResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpSpecExtendRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	SpecName   string `json:"specName,omitempty"`   /*  实例的规格类型。<br>计算增强型的规格可选为：<li>rabbitmq.4u8g.cluster<li>rabbitmq.8u16g.cluster<li>rabbitmq.12u24g.cluster<li>rabbitmq.16u32g.cluster<li>rabbitmq.24u48g.cluster<li>rabbitmq.32u64g.cluster<li>rabbitmq.48u96g.cluster<li>rabbitmq.64u128g.cluster <br>海光-计算增强型的规格可选为：<li>rabbitmq.hg.4u8g.cluster<li>rabbitmq.hg.8u16g.cluster<li>rabbitmq.hg.16u32g.cluster<li>rabbitmq.hg.32u64g.cluster <br>鲲鹏-计算增强型的规格可选为：<li>rabbitmq.kp.4u8g.cluster<li>rabbitmq.kp.8u16g.cluster<li>rabbitmq.kp.16u32g.cluster<li>rabbitmq.kp.32u64g.cluster  */
	AutoPay    *bool  `json:"autoPay"`              /*  是否自动支付，当实例为按需计费模式不生效。true：自动付费，默认值。false：手动付费。  */
}

type AmqpSpecExtendResponse struct {
	StatusCode string                           `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                           `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpSpecExtendReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                           `json:"error"`      /*  错误码，描述错误信息。  */
}

type AmqpSpecExtendReturnObjResponse struct {
	Data *AmqpSpecExtendReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type AmqpSpecExtendReturnObjDataResponse struct{}
