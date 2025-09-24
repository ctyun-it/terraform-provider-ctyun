package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListFlavorFamiliesV41Api
/* 该接口提供用户可用规格族列表查询功能，每种规格族代表不同种类的云主机规格，用户可以根据此接口的返回值了解自己可使用的规格族有哪些<br />&emsp;&emsp;云主机-二代机：X86云主机,包含通用型s2、内存优化型m2。s2、m2实例规格簇均为cpu共享型，上线时间较早<br />&emsp;&emsp;云主机-三代机：X86云主机,包含通用型S3、计算增强型c3、内存优化型m3。S3实例规格簇为cpu共享型，c3、m3实例规格簇为cpu独享,软硬件升级，性能增强<br />&emsp;&emsp;云主机-六代机：X86云主机,包含通用型s6、通用计算增强c6、内存优化型m6。S6实例规格簇为cpu共享型，c6、m6实例规格簇为cpu独享,性能优良，能承载不同业务需求<br />&emsp;&emsp;云主机-七代机：X86云主机,包含通用型s7、通用计算增强c7、内存优化型m7。通用型S7实例规格簇为cpu共享型，c7、m7实例规格簇为cpu独享,提供更大规格更优性能，能满足更高业务需要<br />&emsp;&emsp;国产化云主机：X86与ARM云主机,包含鲲鹏计算增强型kc1、海光计算增强型hc1、飞腾计算增强型fc1、鲲鹏内存优化型km1、海光内存优化型hm1、飞腾内存优化型fm1，对安全性有较高要求的政府或企业应用。<br />
 */ /* &emsp;&emsp;本地盘云主机：X86云主机，包含云主机规格（ip3），提供数据盘为本地盘的云主机<br />&emsp;&emsp;GPU云主机：包含图形加速基础型G5、图形加速基础型G7；计算加速型P2V、计算加速型PI7、计算加速型P8A、计算加速型PS4、计算加速型PI3、计算加速型PI2；图形加速基础型G6、图形加速基础型G7<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />
 */type CtecsListFlavorFamiliesV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListFlavorFamiliesV41Api(client *core.CtyunClient) *CtecsListFlavorFamiliesV41Api {
	return &CtecsListFlavorFamiliesV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/flavor-families/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListFlavorFamiliesV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListFlavorFamiliesV41Request) (*CtecsListFlavorFamiliesV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AzName != "" {
		ctReq.AddParam("azName", req.AzName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListFlavorFamiliesV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListFlavorFamiliesV41Request struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName   string /*  可用区名称，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
}

type CtecsListFlavorFamiliesV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                       `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                       `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListFlavorFamiliesV41ReturnObjResponse `json:"returnObj"`             /*  返回内容  */
}

type CtecsListFlavorFamiliesV41ReturnObjResponse struct {
	FlavorFamilyList []string `json:"flavorFamilyList"` /*  规格族列表  */
}
