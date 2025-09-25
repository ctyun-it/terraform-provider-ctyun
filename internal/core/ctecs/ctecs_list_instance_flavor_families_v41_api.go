package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListInstanceFlavorFamiliesV41Api
/* 该接口提供用户根据指定规格族查询云主机的名称、云主机 ID 及规格详情<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi 请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示 50 条数据
 */type CtecsListInstanceFlavorFamiliesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceFlavorFamiliesV41Api(client *core.CtyunClient) *CtecsListInstanceFlavorFamiliesV41Api {
	return &CtecsListInstanceFlavorFamiliesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/flavor/list-by-families",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceFlavorFamiliesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceFlavorFamiliesV41Request) (*CtecsListInstanceFlavorFamiliesV41Response, error) {
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
	var resp CtecsListInstanceFlavorFamiliesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceFlavorFamiliesV41Request struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName       string `json:"azName,omitempty"`       /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	FlavorFamily string `json:"flavorFamily,omitempty"` /*  规格族名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10138523">规格族</a>来了解规格族信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8326&data=87">查询云主机规格族列表<a>  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListInstanceFlavorFamiliesV41Response struct {
	StatusCode  int32                                                `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）默认值:800  */
	ErrorCode   string                                               `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                               `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                               `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                               `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListInstanceFlavorFamiliesV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListInstanceFlavorFamiliesV41ReturnObjResponse struct {
	CurrentCount int32                                                         `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                         `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                         `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListInstanceFlavorFamiliesV41ReturnObjResultsResponse `json:"results"`                /*  云主机列表  */
}

type CtecsListInstanceFlavorFamiliesV41ReturnObjResultsResponse struct {
	InstanceID   string                                                            `json:"instanceID,omitempty"`   /*  云主机ID  */
	InstanceName string                                                            `json:"instanceName,omitempty"` /*  云主机名称  */
	Flavor       *CtecsListInstanceFlavorFamiliesV41ReturnObjResultsFlavorResponse `json:"flavor"`                 /*  云主机规格详情  */
}

type CtecsListInstanceFlavorFamiliesV41ReturnObjResultsFlavorResponse struct {
	FlavorID     string `json:"flavorID,omitempty"`     /*  规格ID  */
	FlavorName   string `json:"flavorName,omitempty"`   /*  规格名称  */
	FlavorCPU    int32  `json:"flavorCPU,omitempty"`    /*  VCPU个数  */
	FlavorRAM    int32  `json:"flavorRAM,omitempty"`    /*  内存  */
	GpuType      string `json:"gpuType,omitempty"`      /*  GPU类型，取值范围：T4、V100、V100S、A10、A100、atlas 300i pro、mlu370-s4，支持类型会随着功能升级增加  */
	GpuCount     int32  `json:"gpuCount,omitempty"`     /*  GPU数目  */
	GpuVendor    string `json:"gpuVendor,omitempty"`    /*  GPU厂商  */
	VideoMemSize int32  `json:"videoMemSize,omitempty"` /*  GPU显存大小  */
}
