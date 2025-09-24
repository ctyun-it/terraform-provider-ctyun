package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsStopInstanceV41Api
/* 该接口提供用户关闭一台云主机功能，可选两种关机模式：普通关机、强制关机<br />请求下发后云主机变为关机中（stopping）状态，待异步任务完成后，云主机变为关机（stopped）状态<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项：</b><br />&emsp;&emsp;单台操作：当前接口只能操作单台云主机，关闭多台云主机请使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8304&data=87">关闭多台云主机</a>进行操作<br />&emsp;&emsp;异步接口：该接口为异步接口，请求过后会拿到任务ID（jobID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5543&data=87">查询一个异步任务的结果</a>来查询操作是否成功
 */type CtecsStopInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsStopInstanceV41Api(client *core.CtyunClient) *CtecsStopInstanceV41Api {
	return &CtecsStopInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/stop-instance",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsStopInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsStopInstanceV41Request) (*CtecsStopInstanceV41Response, error) {
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
	var resp CtecsStopInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsStopInstanceV41Request struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	Force      *bool  `json:"force"`                /*  是否强制关机，取值范围：<br />true：强制关机，<br />false：普通关机<br />注：默认值false  */
}

type CtecsStopInstanceV41Response struct {
	StatusCode  int32                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                 `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                 `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                 `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsStopInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsStopInstanceV41ReturnObjResponse struct {
	JobID string `json:"jobID,omitempty"` /*  关机任务ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5543&data=87">查询一个异步任务的结果</a>来查询操作是否成功  */
}
