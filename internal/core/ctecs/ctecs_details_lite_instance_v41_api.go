package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsLiteInstanceV41Api
/* 该接口提供用户轻量型云主机信息查询功能，用户可以根据此接口的返回值了解自己轻量型云主机的详细信息<br /><b>准备工作：</b><br/>&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br /><b>注意事项：</b><br />&emsp;&emsp;单台查询：当前接口只能查询单台云主机信息，查询多台云主机信息请使用接口<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11981&data=87">查询轻量型云主机列表</a>进行查询<br />
 */type CtecsDetailsLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsLiteInstanceV41Api(client *core.CtyunClient) *CtecsDetailsLiteInstanceV41Api {
	return &CtecsDetailsLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/lite/details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsLiteInstanceV41Request) (*CtecsDetailsLiteInstanceV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceID", req.InstanceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDetailsLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsLiteInstanceV41Request struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID string /*  轻量型云主机ID，您可以查看<a href="https://www.ctyun.cn/products/lite-ecs">轻量型云主机</a>了解轻量型云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11981&data=87">查询轻量型云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11980&data=87">创建轻量型云主机</a>  */
}

type CtecsDetailsLiteInstanceV41Response struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsDetailsLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsDetailsLiteInstanceV41ReturnObjResponse struct {
	ResourceID      string                                                         `json:"resourceID,omitempty"`     /*  云主机资源ID  */
	InstanceID      string                                                         `json:"instanceID,omitempty"`     /*  云主机ID  */
	DisplayName     string                                                         `json:"displayName,omitempty"`    /*  云主机显示名称  */
	InstanceName    string                                                         `json:"instanceName,omitempty"`   /*  云主机名称  */
	OsType          int32                                                          `json:"osType,omitempty"`         /*  操作系统类型，详见枚举值表格  */
	InstanceStatus  string                                                         `json:"instanceStatus,omitempty"` /*  云主机状态，请通过<a href="https://www.ctyun.cn/document/10026730/10741614">状态枚举值</a>查看云主机使用状态  */
	ExpiredTime     string                                                         `json:"expiredTime,omitempty"`    /*  到期时间  */
	UpdatedTime     string                                                         `json:"updatedTime,omitempty"`    /*  更新时间  */
	CreatedTime     string                                                         `json:"createdTime,omitempty"`    /*  创建时间  */
	AttachedVolume  []string                                                       `json:"attachedVolume"`           /*  附加卷  */
	Addresses       *CtecsDetailsLiteInstanceV41ReturnObjAddressesResponse         `json:"addresses"`                /*  网络地址信息  */
	SecGroupList    []*CtecsDetailsLiteInstanceV41ReturnObjSecGroupListResponse    `json:"secGroupList"`             /*  安全组信息  */
	NetworkCardList []*CtecsDetailsLiteInstanceV41ReturnObjNetworkCardListResponse `json:"networkCardList"`          /*  网卡信息  */
	Image           *CtecsDetailsLiteInstanceV41ReturnObjImageResponse             `json:"image"`                    /*  镜像信息  */
	Flavor          *CtecsDetailsLiteInstanceV41ReturnObjFlavorResponse            `json:"flavor"`                   /*  规格信息  */
	VpcID           string                                                         `json:"vpcID,omitempty"`          /*  vpc ID  */
	VpcName         string                                                         `json:"vpcName,omitempty"`        /*  vpc名称  */
	ZabbixName      string                                                         `json:"zabbixName,omitempty"`     /*  监控对象名称  */
	Bandwidth       int32                                                          `json:"bandwidth,omitempty"`      /*  带宽  */
	BootDiskSize    int32                                                          `json:"bootDiskSize,omitempty"`   /*  系统盘大小  */
}

type CtecsDetailsLiteInstanceV41ReturnObjAddressesResponse struct {
	AddressList []*CtecsDetailsLiteInstanceV41ReturnObjAddressesAddressListResponse `json:"addressList"` /*  网络地址列表  */
}

type CtecsDetailsLiteInstanceV41ReturnObjSecGroupListResponse struct {
	SecurityGroupID   string `json:"securityGroupID,omitempty"`   /*  安全组ID  */
	SecurityGroupName string `json:"securityGroupName,omitempty"` /*  安全组名称  */
}

type CtecsDetailsLiteInstanceV41ReturnObjNetworkCardListResponse struct {
	IPv4Address string `json:"IPv4Address,omitempty"` /*  IPv4地址  */
	IPv6Address string `json:"IPv6Address,omitempty"` /*  IPv6地址  */
	SubnetID    string `json:"subnetID,omitempty"`    /*  所处的子网ID  */
}

type CtecsDetailsLiteInstanceV41ReturnObjImageResponse struct {
	ImageID   string `json:"imageID,omitempty"`   /*  镜像ID  */
	ImageName string `json:"imageName,omitempty"` /*  镜像名称  */
}

type CtecsDetailsLiteInstanceV41ReturnObjFlavorResponse struct {
	FlavorID   string `json:"flavorID,omitempty"`   /*  规格ID  */
	FlavorName string `json:"flavorName,omitempty"` /*  规格名称  */
	FlavorCPU  int32  `json:"flavorCPU,omitempty"`  /*  VCPU  */
	FlavorRAM  int32  `json:"flavorRAM,omitempty"`  /*  内存  */
}

type CtecsDetailsLiteInstanceV41ReturnObjAddressesAddressListResponse struct {
	Addr    string `json:"addr,omitempty"`    /*  地址  */
	Version int32  `json:"version,omitempty"` /*  IP版本  */
	RawType string `json:"type,omitempty"`    /*  网络类型  */
}
