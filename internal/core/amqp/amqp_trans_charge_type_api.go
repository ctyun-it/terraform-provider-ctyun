package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// AmqpTransChargeTypeApi
/* 按需转包周期。
 */type AmqpTransChargeTypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewAmqpTransChargeTypeApi(client *core.CtyunClient) *AmqpTransChargeTypeApi {
	return &AmqpTransChargeTypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/transChargeType",
			ContentType:  "application/json",
		},
	}
}

func (a *AmqpTransChargeTypeApi) Do(ctx context.Context, credential core.Credential, req *AmqpTransChargeTypeRequest) (*AmqpTransChargeTypeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*AmqpTransChargeTypeRequest
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
	var resp AmqpTransChargeTypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type AmqpTransChargeTypeRequest struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。您可以通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	CycleCnt   int32  `json:"cycleCnt,omitempty"`   /*  付费周期，单位为月，取值：1~6,12,24,36。  */
	AutoPay    *bool  `json:"autoPay"`              /*  是否自动支付。true：自动付费，默认值。false：手动付费。  */
}

type AmqpTransChargeTypeResponse struct {
	StatusCode string                                `json:"statusCode"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                `json:"message"`    /*  描述状态。  */
	ReturnObj  *AmqpTransChargeTypeReturnObjResponse `json:"returnObj"`  /*  返回对象。  */
	Error      string                                `json:"error"`      /*  错误码，描述错误信息。  */
}

type AmqpTransChargeTypeReturnObjResponse struct {
	Data []*AmqpTransChargeTypeReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type AmqpTransChargeTypeReturnObjDataResponse struct {
	MasterOrderId   string `json:"masterOrderId"`
	MasterOrderNo   string `json:"masterOrderNo"`
	MasterOrderType string `json:"masterOrderType"`
}
