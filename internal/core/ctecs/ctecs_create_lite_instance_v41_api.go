package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateLiteInstanceV41Api
/* 该接口提供用户创建一台包年包月的轻量型云主机<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u6784%u9020%u8BF7%u6C42&data=87&vid=81">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=%u8BA4%u8BC1%u9274%u6743&data=87&vid=81">认证鉴权</a><br />&emsp;&emsp;计费模式：确认开通云主机的计费模式及计费项，详细查看<a href="https://www.ctyun.cn/document/10114925/10115664">计费方式及计费项</a><br />&emsp;&emsp;产品选型：购买轻量型云主机前，请先阅读<a href="https://www.ctyun.cn/document/10114925/10268652">实例套餐</a>了解轻量型云主机的规格套餐，并通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11998&data=87">查询轻量型云主机的规格套餐资源</a>接口，获取当前资源池可用轻量型云主机规格信息<br /><b>注意事项：</b><br />&emsp;&emsp;成本估算：了解云主机的计费项，详细查看<a href="https://www.ctyun.cn/document/10114925/10115703">价格总览</a>进行成本估算<br />&emsp;&emsp;用户配额：确认个人在不同资源池下资源配额，可以通过<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9714&data=87">用户配额查询</a>接口进行查询<br />&emsp;&emsp;异步接口：该接口为异步接口，下单过后会拿到主订单ID（masterOrderID），后续可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9607&data=87&isNormal=1">根据masterOrderID查询云主机ID</a>，使用主订单ID来对订单情况与开通成功后的资源ID进行查询<br />
 */type CtecsCreateLiteInstanceV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateLiteInstanceV41Api(client *core.CtyunClient) *CtecsCreateLiteInstanceV41Api {
	return &CtecsCreateLiteInstanceV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/lite/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateLiteInstanceV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateLiteInstanceV41Request) (*CtecsCreateLiteInstanceV41Response, error) {
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
	var resp CtecsCreateLiteInstanceV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateLiteInstanceV41Request struct {
	ClientToken     string                                           `json:"clientToken,omitempty"`     /*  客户端存根，用于保证订单幂等性，保留时间为24小时  */
	RegionID        string                                           `json:"regionID,omitempty"`        /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	AzName          string                                           `json:"azName,omitempty"`          /*  可用区名称，如果是4.0资源池，必须提供可用区名称。您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解可用区 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br />注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）  */
	DisplayName     string                                           `json:"displayName,omitempty"`     /*  云主机显示名称，长度为2~15字符  */
	FlavorSetType   string                                           `json:"flavorSetType,omitempty"`   /*  规格套餐类型，取值范围：<br />fix：固定套餐，<br />band：带宽套餐<br />选择带宽套餐必须设置系统盘大小和带宽大小，选择固定套餐会忽略传入的系统盘大小和带宽大小  */
	FlavorName      string                                           `json:"flavorName,omitempty"`      /*  规格套餐名称，<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11998&data=87">查询轻量型云主机的规格套餐资源</a>  */
	ImageID         string                                           `json:"imageID,omitempty"`         /*  镜像ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10030151">镜像概述</a>来了解云主机镜像<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89">查询可以使用的镜像资源</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4765&data=89">创建私有镜像（云主机系统盘）</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=5230&data=89">创建私有镜像（云主机数据盘）</a><br />注：同一镜像名称在不同资源池的镜像ID是不同的，调用前需确认镜像ID是否归属当前资源池  */
	CycleCount      int32                                            `json:"cycleCount,omitempty"`      /*  订购时长，该参数需要与cycleType一同使用<br />注：最长订购周期为60个月（5年）；cycleType与cycleCount一起填写；按量付费（即onDemand为true）时，无需填写该参数（填写无效）  */
	CycleType       string                                           `json:"cycleType,omitempty"`       /*  订购周期类型，取值范围：<br />MONTH：按月，<br />YEAR：按年<br />注：cycleType与cycleCount一起填写；按量付费（即onDemand为true）时，无需填写该参数（填写无效）  */
	IpVersion       string                                           `json:"ipVersion,omitempty"`       /*  弹性IP版本，取值范围：<br />ipv4：v4地址，<br />ipv6：v6地址<br />不指定默认为ipv4。注：请先确认该资源池是否支持ipv6<br/>  */
	BootDiskSize    int32                                            `json:"bootDiskSize,omitempty"`    /*  系统盘大小，带宽套餐时填写，固定套餐时填写会忽略，单位：GB，取值范围获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11998&data=87">查询轻量型云主机的规格套餐资源</a>  */
	Bandwidth       int32                                            `json:"bandwidth,omitempty"`       /*  带宽大小带，带宽套餐时填写，固定套餐时填写会忽略，单位：Mbit/s。取值范围获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=11998&data=87">查询轻量型云主机的规格套餐资源</a>  */
	DataDiskList    []*CtecsCreateLiteInstanceV41DataDiskListRequest `json:"dataDiskList"`              /*  数据盘信息列表，注：同一云主机下最多可挂载5块数据盘  */
	UserPassword    string                                           `json:"userPassword,omitempty"`    /*  用户密码，满足以下规则：<br />长度在8～30个字符；<br />必须包含大写字母、小写字母、数字以及特殊符号中的三项；<br />特殊符号可选：()`~!@#$%^&*_-+=｜{}[]:;'<>,.?/\且不能以斜线号 / 开头；<br />不能包含3个及以上连续字符；<br />Linux镜像不能包含镜像用户名（root）、用户名的倒序（toor）、用户名大小写变化（如RoOt、rOot等）；<br />Windows镜像不能包含镜像用户名（Administrator）、用户名大小写变化（adminiSTrator等）  */
	AutoRenewStatus int32                                            `json:"autoRenewStatus,omitempty"` /*  本参数表示是否自动续订,取值范围：<br />0：不续费，<br />1：自动续费<br />注：按月购买，自动续订周期为1个月；按年购买，自动续订周期为1年  */
}

type CtecsCreateLiteInstanceV41DataDiskListRequest struct {
	DiskType string `json:"diskType,omitempty"` /*  云硬盘类型，取值范围：<br />SATA：普通IO，<br />SAS：高IO，<br />SSD：超高IO，<br />FAST-SSD：极速型SSD<br />您可以查看<a href="https://www.ctyun.cn/document/10027696/10162918">磁盘类型及性能介绍</a>磁盘类型相关信息  */
	DiskSize int32  `json:"diskSize,omitempty"` /*  磁盘容量大小单位为GB，取值范围：[10-32768]  */
}

type CtecsCreateLiteInstanceV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                       `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                       `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsCreateLiteInstanceV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCreateLiteInstanceV41ReturnObjResponse struct {
	MasterOrderID    string `json:"masterOrderID,omitempty"`    /*  主订单ID。调用方在拿到masterOrderID之后，可以使用materOrderID进一步确认订单状态及资源状态<br />查询订单状态及资源UUID：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9607&data=87&isNormal=1">根据masterOrderID查询云主机ID</a>  */
	MasterOrderNO    string `json:"masterOrderNO,omitempty"`    /*  订单号  */
	MasterResourceID string `json:"masterResourceID,omitempty"` /*  主资源ID  */
	RegionID         string `json:"regionID,omitempty"`         /*  资源所属资源池ID  */
}
