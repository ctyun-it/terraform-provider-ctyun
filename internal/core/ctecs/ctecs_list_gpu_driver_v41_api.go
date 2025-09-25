package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListGpuDriverV41Api
/* 该接口提供用户使用GPU 云主机规格查询可用驱动版本<br /><b>准备工作</b>：&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsListGpuDriverV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListGpuDriverV41Api(client *core.CtyunClient) *CtecsListGpuDriverV41Api {
	return &CtecsListGpuDriverV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/gpu-driver/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListGpuDriverV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListGpuDriverV41Request) (*CtecsListGpuDriverV41Response, error) {
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
	var resp CtecsListGpuDriverV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListGpuDriverV41Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	FlavorID string `json:"flavorID,omitempty"` /*  云主机规格ID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8327&data=87&isNormal=1&vid=81">云主机规格资源查询</a>来查看最新的天翼云具体资源池的云主机规格列表  */
}

type CtecsListGpuDriverV41Response struct {
	StatusCode  int32                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                  `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                  `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListGpuDriverV41ReturnObjResponse `json:"returnObj"`             /*  返回内容  */
}

type CtecsListGpuDriverV41ReturnObjResponse struct {
	GpuDriverList []string `json:"gpuDriverList"` /*  驱动列表  */
}
