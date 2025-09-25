package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDeleteEipV41Api
/* 调用此接口可删除未绑定云产品实例的弹性公网IP<br /><b>准备工作</b>：<br /> &emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsDeleteEipV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDeleteEipV41Api(client *core.CtyunClient) *CtecsDeleteEipV41Api {
	return &CtecsDeleteEipV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/eip/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDeleteEipV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDeleteEipV41Request) (*CtecsDeleteEipV41Response, error) {
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
	var resp CtecsDeleteEipV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDeleteEipV41Request struct {
	ClientToken string `json:"clientToken,omitempty"` /*  客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一，使用同一个ClientToken值，则代表为同一个请求。保留时间为24小时  */
	RegionID    string `json:"regionID,omitempty"`    /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	EipID       string `json:"eipID,omitempty"`       /*  弹性公网IP的ID，您可以查看<a href="https://www.ctyun.cn/document/10026753/10026909">产品定义-弹性IP</a>来了解弹性公网IP <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=13&api=8652&data=101">查询指定地域已创建的弹性IP</a><br />  */
}

type CtecsDeleteEipV41Response struct {
	StatusCode  int32                               `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                              `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                              `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                              `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDeleteEipV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDeleteEipV41ReturnObjResponse struct {
	MasterOrderID        string `json:"masterOrderID,omitempty"`        /*  订单ID。调用方在拿到masterOrderID之后，可以使用materOrderID进一步确认订单状态及资源状态  */
	MasterOrderNO        string `json:"masterOrderNO,omitempty"`        /*  订单号  */
	RegionID             string `json:"regionID,omitempty"`             /*  资源池ID  */
	MasterResourceStatus string `json:"masterResourceStatus,omitempty"` /*  资源状态  */
}
