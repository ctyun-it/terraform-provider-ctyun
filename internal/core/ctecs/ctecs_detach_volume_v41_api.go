package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetachVolumeV41Api
/* 支持云主机卸载云硬盘。 <br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsDetachVolumeV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetachVolumeV41Api(client *core.CtyunClient) *CtecsDetachVolumeV41Api {
	return &CtecsDetachVolumeV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/volume/detach",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetachVolumeV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetachVolumeV41Request) (*CtecsDetachVolumeV41Response, error) {
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
	var resp CtecsDetachVolumeV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetachVolumeV41Request struct {
	DiskID     string `json:"diskID,omitempty"`     /*  磁盘ID，您可以查看<a href="https://www.ctyun.cn/document/10027696/10027930">产品定义-云硬盘</a>来了解云硬盘 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7338&data=48">云硬盘列表查询</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=35&api=7332&data=48&isNormal=1&vid=45">创建云硬盘</a>  */
	RegionID   string `json:"regionID,omitempty"`   /*  如如本地语境支持保存regionID，那么建议传递。资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string `json:"instanceID,omitempty"` /*  云主机ID，用于共享盘指定卸载某个主机。  */
}

type CtecsDetachVolumeV41Response struct {
	StatusCode  int32                                    `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                   `json:"errorCode,omitempty"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                                   `json:"error,omitempty"`       /*  业务细分码，为product.module.code三段式码  */
	Message     string                                   `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                   `json:"description,omitempty"` /*  中文描述信息  */
	ErrorDetail *CtecsDetachVolumeV41ErrorDetailResponse `json:"errorDetail"`           /*  错误明细。一般情况下，会对订单侧(bss)的云硬盘订单业务相关的错误做明确的错误映射和提升，有唯一对应的errorCode<br />其他订单侧(bss)的错误，以Ebs.Order.ProcFailed的errorCode统一映射返回，并在errorDetail中返回订单侧的详细错误信息  */
	ReturnObj   *CtecsDetachVolumeV41ReturnObjResponse   `json:"returnObj"`             /*  返回参数  */
}

type CtecsDetachVolumeV41ErrorDetailResponse struct {
	BssErrCode       string `json:"bssErrCode,omitempty"`       /*  bss错误明细码，包含于bss格式化JSON错误信息中  */
	BssErrMsg        string `json:"bssErrMsg,omitempty"`        /*  bss错误信息，包含于bss格式化JSON错误信息中  */
	BssOrigErr       string `json:"bssOrigErr,omitempty"`       /*  无法明确解码bss错误信息时，原样透出的bss错误信息  */
	BssErrPrefixHint string `json:"bssErrPrefixHint,omitempty"` /*  bss格式化JSON错误信息的前置提示信息  */
}

type CtecsDetachVolumeV41ReturnObjResponse struct {
	DiskJobID     string `json:"diskJobID,omitempty"`     /*  异步任务ID，可通过公共接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5865&data=87&isNormal=1&vid=81">通用任务状态查询</a>该jobID的异步任务最终执行结果  */
	DiskRequestID string `json:"diskRequestID,omitempty"` /*  异步任务ID，可通过公共接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5865&data=87&isNormal=1&vid=81">通用任务状态查询</a>该jobID的异步任务最终执行结果  */
}
