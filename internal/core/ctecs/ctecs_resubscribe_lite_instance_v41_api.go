package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsResubscribeLiteInstanceV41Api
/* 该接口提供用户续订一台包年或者包月的轻量型云主机功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项</b>：<br />&emsp;&emsp;成本估算：了解云主机的计费项，详细查看<a href="https://www.ctyun.cn/document/10114925/10115703">价格总览</a>进行成本估算<br />&emsp;&emsp;异步接口：该接口为异步接口，下单过后会拿到主订单ID（masterOrderID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9607&data=87&isNormal=1&vid=81">根据masterOrderID查询云主机ID</a>，使用主订单ID来对订单情况与开通成功后的资源ID进行查询<br />
 */type CtecsResubscribeLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsResubscribeLiteInstanceV41Api(client *core.CtyunClient) *CtecsResubscribeLiteInstanceV41Api {
	return &CtecsResubscribeLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/lite/resubscribe",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsResubscribeLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsResubscribeLiteInstanceV41Request) (*CtecsResubscribeLiteInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsResubscribeLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsResubscribeLiteInstanceV41Request struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性，保留时间为24小时  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID  string `json:"instanceID,omitempty"`  /*  轻量型云主机ID，您可以查看<a href="https://www.ctyun.cn/products/lite-ecs">轻量型云主机</a>了解轻量型云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11981&data=87">查询轻量型云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11980&data=87">创建轻量型云主机</a>  */
	CycleCount  int32  `json:"cycleCount,omitempty"`  /*  订购时长，<br />注：最长订购周期为36个月（3年）  */
	CycleType   string `json:"cycleType,omitempty"`   /*  订购周期类型，取值范围：<br />MONTH：按月，<br />YEAR：按年  */
}

type CtecsResubscribeLiteInstanceV41Response struct {
	StatusCode  int32                                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                            `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                            `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                            `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                            `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsResubscribeLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsResubscribeLiteInstanceV41ReturnObjResponse struct {
	MasterOrderID string `json:"masterOrderID,omitempty"` /*  订单ID。调用方在拿到masterOrderID之后，在若干错误情况下，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO string `json:"masterOrderNO,omitempty"` /*  订单号  */
	RegionID      string `json:"regionID,omitempty"`      /*  资源所属资源池ID  */
}
