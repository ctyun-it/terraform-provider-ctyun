package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2RenewInstanceApi
/* 为包年包月的分布式缓存Redis实例续费。
 */type Dcs2RenewInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2RenewInstanceApi(client *core.CtyunClient) *Dcs2RenewInstanceApi {
	return &Dcs2RenewInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/lifeCycleServant/renewInstance",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2RenewInstanceApi) Do(ctx context.Context, credential core.Credential, req *Dcs2RenewInstanceRequest) (*Dcs2RenewInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp Dcs2RenewInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2RenewInstanceRequest struct {
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	ProdInstId string `json:"prodInstId,omitempty"` /*  实例Id  */
	AutoPay    *bool  `json:"autoPay"`              /*  是否自动支付(仅对包周期实例有效)：<li>true：自动付费<li>false：手动付费(默认值)<br>说明：选择为手动付费时，您需要在控制台的右上角选择我的订单，，然后单击左侧导航栏的订单管理 -> 待支付订单，找到目标订单进行支付。  */
	Period     int32  `json:"period,omitempty"`     /*  订购时长(月)，仅当包周期实例需要进行续费，取值范围：1-6,12,24,36。  */
}

type Dcs2RenewInstanceResponse struct {
	StatusCode int32                               `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                              `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2RenewInstanceReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                              `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                              `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                              `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2RenewInstanceReturnObjResponse struct {
	ErrorMessage      string                                                 `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted         *bool                                                  `json:"submitted"`              /*  是否提交  */
	OrderPlacedEvents []*Dcs2RenewInstanceReturnObjOrderPlacedEventsResponse `json:"orderPlacedEvents"`      /*  收费项  */
}

type Dcs2RenewInstanceReturnObjOrderPlacedEventsResponse struct {
	ErrorMessage string  `json:"errorMessage,omitempty"` /*  错误信息  */
	Submitted    *bool   `json:"submitted"`              /*  是否提交  */
	NewOrderId   string  `json:"newOrderId,omitempty"`   /*  订单ID  */
	NewOrderNo   string  `json:"newOrderNo,omitempty"`   /*  订单号  */
	TotalPrice   float64 `json:"totalPrice"`             /*  总价  */
}
