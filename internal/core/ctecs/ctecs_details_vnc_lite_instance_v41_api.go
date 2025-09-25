package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsVncLiteInstanceV41Api
/* 该接口提供用户查询一台轻量型云主机的Web管理终端地址<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />获取到轻量型云主机访问地址如下所示：```http://203.193.231.250:60001/vnc_auto.html?token=f0b6ed3f-e049-4961-8c87-07815e058662&instance_name=VM-fbb407b6(aca51c6e-3c72-47ce-894e-1e2fbfdbabf9)```有如上信息后，在轻量型云主机为开机状态时，直接浏览器输入访问地址进行访问
 */type CtecsDetailsVncLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsVncLiteInstanceV41Api(client *core.CtyunClient) *CtecsDetailsVncLiteInstanceV41Api {
	return &CtecsDetailsVncLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/lite/vnc/details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsVncLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsVncLiteInstanceV41Request) (*CtecsDetailsVncLiteInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDetailsVncLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsVncLiteInstanceV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  轻量型云主机ID，您可以查看<a href="https://www.ctyun.cn/products/lite-ecs">轻量型云主机</a>了解轻量型云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11981&data=87">查询轻量型云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11980&data=87">创建轻量型云主机</a>  */
}

type CtecsDetailsVncLiteInstanceV41Response struct {
	StatusCode  int32                                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                           `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                           `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                           `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                           `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDetailsVncLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDetailsVncLiteInstanceV41ReturnObjResponse struct {
	Token string `json:"token,omitempty"` /*  token  */
}
