package ctgkafka

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtgkafkaTransChargeTypeV3Api
/* 按需转包周期。
 */type CtgkafkaTransChargeTypeV3Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtgkafkaTransChargeTypeV3Api(client *core.CtyunClient) *CtgkafkaTransChargeTypeV3Api {
	return &CtgkafkaTransChargeTypeV3Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v3/instances/transChargeType",
			ContentType:  "application/json",
		},
	}
}

func (a *CtgkafkaTransChargeTypeV3Api) Do(ctx context.Context, credential core.Credential, req *CtgkafkaTransChargeTypeV3Request) (*CtgkafkaTransChargeTypeV3Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(struct {
		*CtgkafkaTransChargeTypeV3Request
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
	var resp CtgkafkaTransChargeTypeV3Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtgkafkaTransChargeTypeV3Request struct {
	RegionId   string `json:"regionId,omitempty"`   /*  实例的资源池ID。<br>获取方法如下：<br><li>方法一：通过查询<a href="https://www.ctyun.cn/document/10029624/11008434">分布式消息服务Kafka资源池附录文档</a>。<br><li>方法二：通过调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">查询可用的资源池</a>API接口查。  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例ID。  */
	CycleCnt   int32  `json:"cycleCnt,omitempty"`   /*  付费周期，单位为月，取值：1~6,12,24,36。  */
	AutoPay    *bool  `json:"autoPay"`              /*  是否自动支付。true：自动付费，默认值。false：手动付费。  */
}

type CtgkafkaTransChargeTypeV3Response struct {
	StatusCode string                                      `json:"statusCode,omitempty"` /*  接口系统层面状态码。成功："800"，失败："900"。  */
	Message    string                                      `json:"message,omitempty"`    /*  描述状态。  */
	ReturnObj  *CtgkafkaTransChargeTypeV3ReturnObjResponse `json:"returnObj"`            /*  返回对象。  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，描述错误信息。  */
}

type CtgkafkaTransChargeTypeV3ReturnObjResponse struct {
	Data *CtgkafkaTransChargeTypeV3ReturnObjDataResponse `json:"data"` /*  返回数据。  */
}

type CtgkafkaTransChargeTypeV3ReturnObjDataResponse struct{}
