package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListFlavorLiteInstanceV41Api
/* 该接口提供用户可用规格列表查询功能，可返回轻量型云主机规格的详细信息，用户可以根据此接口的返回值了解自己可使用的云主机规格有哪些<br />您可以通过<a href="https://www.ctyun.cn/document/10114925/10268652">实例套餐</a>了解轻量型云主机的规格套餐信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsListFlavorLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListFlavorLiteInstanceV41Api(client *core.CtyunClient) *CtecsListFlavorLiteInstanceV41Api {
	return &CtecsListFlavorLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/lite/flavor/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListFlavorLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListFlavorLiteInstanceV41Request) (*CtecsListFlavorLiteInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	if req.FlavorSetType != "" {
		ctReq.AddParam("flavorSetType", req.FlavorSetType)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListFlavorLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListFlavorLiteInstanceV41Request struct {
	RegionID      string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName        string /*  可用区名称，如果是4.0资源池，必须提供可用区名称。您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，无需填写本字段）  */
	FlavorSetType string /*  规格套餐类型，取值范围：<br />fix：固定套餐，<br />band：带宽套餐  */
}

type CtecsListFlavorLiteInstanceV41Response struct {
	StatusCode  int32                                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                           `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                           `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                           `json:"message,omitempty"`     /*  英文描述信息   */
	Description string                                           `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListFlavorLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsListFlavorLiteInstanceV41ReturnObjResponse struct {
	CurrentCount int32                                                     `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                     `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                     `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListFlavorLiteInstanceV41ReturnObjResultsResponse `json:"results"`                /*  规格列表  */
}

type CtecsListFlavorLiteInstanceV41ReturnObjResultsResponse struct {
	FlavorSetType          string  `json:"flavorSetType,omitempty"`      /*  规格套餐类型，取值范围：<br />fix：固定套餐，<br />band：带宽套餐  */
	FlavorName             string  `json:"flavorName,omitempty"`         /*  规格套餐名称  */
	FlavorSeries           string  `json:"flavorSeries,omitempty"`       /*  规格系列  */
	FlavorType             string  `json:"flavorType,omitempty"`         /*  规格类型  */
	FlavorRAM              int32   `json:"flavorRAM,omitempty"`          /*  内存大小，单位为G  */
	FlavorCPU              int32   `json:"flavorCPU,omitempty"`          /*  VCPU个数  */
	FlavorBandwidth        int32   `json:"flavorBandwidth,omitempty"`    /*  固定套餐带宽大小，当flavorSetType为fix时展示  */
	FlavorBootDiskSize     int32   `json:"flavorBootDiskSize,omitempty"` /*  固定套餐系统盘大小，当flavorSetType为fix时展示  */
	FlavorBandwidthList    []int32 `json:"flavorBandwidthList"`          /*  带宽套餐带宽大小取值列表，当flavorSetType为band时展示  */
	FlavorBootDiskSizeList []int32 `json:"flavorBootDiskSizeList"`       /*  带宽套餐系统盘大小取值列表，当flavorSetType为band时展示  */
}
