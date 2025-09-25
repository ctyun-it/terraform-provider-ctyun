package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// Dcs2ChangeAutoRenewStatusApi
/* 自动续费开关，可设置实例到期后自动续费周期，该操作为幂等操作，重复执行返回相同结果。
 */type Dcs2ChangeAutoRenewStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewDcs2ChangeAutoRenewStatusApi(client *core.CtyunClient) *Dcs2ChangeAutoRenewStatusApi {
	return &Dcs2ChangeAutoRenewStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/res/spuInst/changeAutoRenewStatus",
			ContentType:  "application/json",
		},
	}
}

func (a *Dcs2ChangeAutoRenewStatusApi) Do(ctx context.Context, credential core.Credential, req *Dcs2ChangeAutoRenewStatusRequest) (*Dcs2ChangeAutoRenewStatusResponse, error) {
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
	var resp Dcs2ChangeAutoRenewStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type Dcs2ChangeAutoRenewStatusRequest struct {
	RegionId            string   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=49&api=7830&isNormal=1&vid=270">查询可用的资源池</a>  */
	AutoRenewCycleType  string   `json:"autoRenewCycleType,omitempty"`  /*  自动续订周期类型。<li>3：月<li>5：年  */
	AutoRenewCycleCount int32    `json:"autoRenewCycleCount,omitempty"` /*  自动续费周期(月)<br>autoRenew=true时必填，可选：1-6,12,24,36  */
	AutoRenewStatus     string   `json:"autoRenewStatus,omitempty"`     /*  自动续订状态。<li>1：自动续订<li>0：不自动续订  */
	Source              string   `json:"source,omitempty"`              /*  来源。<li>8：表示自动续订操作来源  */
	ProdInstIds         []string `json:"prodInstIds"`                   /*  待改变自动续订状态的实例列表  */
}

type Dcs2ChangeAutoRenewStatusResponse struct {
	StatusCode int32                                       `json:"statusCode,omitempty"` /*  响应状态码<li>800：成功<li>900：失败  */
	Message    string                                      `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *Dcs2ChangeAutoRenewStatusReturnObjResponse `json:"returnObj"`            /*  返回数据对象，数据见returnObj  */
	RequestId  string                                      `json:"requestId,omitempty"`  /*  请求 ID  */
	Code       string                                      `json:"code,omitempty"`       /*  响应码描述  */
	Error      string                                      `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type Dcs2ChangeAutoRenewStatusReturnObjResponse struct {
	StatusCode int32 `json:"statusCode,omitempty"` /*  响应状态码<li>800：自动续订下单成功<li>900：自动续订下单失败  */
}
