package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListLiteInstanceV41Api
/* 该接口提供用户多台轻量型云主机信息查询功能，用户可以根据此接口的返回值得到多轻量型台云主机信息<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;分页查询：当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />
 */type CtecsListLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListLiteInstanceV41Api(client *core.CtyunClient) *CtecsListLiteInstanceV41Api {
	return &CtecsListLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/lite/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListLiteInstanceV41Request) (*CtecsListLiteInstanceV41Response, error) {
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
	var resp CtecsListLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListLiteInstanceV41Request struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListLiteInstanceV41Response struct {
	StatusCode  int32                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                     `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                     `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                     `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsListLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsListLiteInstanceV41ReturnObjResponse struct {
	CurrentCount int32                                               `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                               `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                               `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListLiteInstanceV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListLiteInstanceV41ReturnObjResultsResponse struct {
	ResourceID     string                                                     `json:"resourceID,omitempty"`     /*  云主机资源ID  */
	InstanceID     string                                                     `json:"instanceID,omitempty"`     /*  云主机ID  */
	DisplayName    string                                                     `json:"displayName,omitempty"`    /*  云主机显示名称  */
	InstanceName   string                                                     `json:"instanceName,omitempty"`   /*  云主机名称  */
	OsType         int32                                                      `json:"osType,omitempty"`         /*  操作系统类型，详见枚举值表格  */
	InstanceStatus string                                                     `json:"instanceStatus,omitempty"` /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	ExpiredTime    string                                                     `json:"expiredTime,omitempty"`    /*  到期时间  */
	CreatedTime    string                                                     `json:"createdTime,omitempty"`    /*  创建时间  */
	Addresses      *CtecsListLiteInstanceV41ReturnObjResultsAddressesResponse `json:"addresses"`                /*  网络地址信息  */
	Image          *CtecsListLiteInstanceV41ReturnObjResultsImageResponse     `json:"image"`                    /*  镜像信息  */
	Flavor         *CtecsListLiteInstanceV41ReturnObjResultsFlavorResponse    `json:"flavor"`                   /*  规格信息  */
	Bandwidth      int32                                                      `json:"bandwidth,omitempty"`      /*  带宽  */
	BootDiskSize   int32                                                      `json:"bootDiskSize,omitempty"`   /*  系统盘大小  */
}

type CtecsListLiteInstanceV41ReturnObjResultsAddressesResponse struct {
	AddressList []*CtecsListLiteInstanceV41ReturnObjResultsAddressesAddressListResponse `json:"addressList"` /*  网络地址列表  */
}

type CtecsListLiteInstanceV41ReturnObjResultsImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsListLiteInstanceV41ReturnObjResultsFlavorResponse struct {
	FlavorID   string `json:"flavorID,omitempty"`   /*  规格ID  */
	FlavorName string `json:"flavorName,omitempty"` /*  规格名称  */
	FlavorCPU  int32  `json:"flavorCPU,omitempty"`  /*  VCPU  */
	FlavorRAM  int32  `json:"flavorRAM,omitempty"`  /*  内存  */
}

type CtecsListLiteInstanceV41ReturnObjResultsAddressesAddressListResponse struct {
	Addr    string `json:"addr,omitempty"`    /*  地址  */
	Version int32  `json:"version,omitempty"` /*  IP版本  */
	RawType string `json:"type,omitempty"`    /*  网络类型  */
}
