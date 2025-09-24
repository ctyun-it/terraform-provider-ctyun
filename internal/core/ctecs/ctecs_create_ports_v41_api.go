package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreatePortsV41Api
/* 创建弹性网卡<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsCreatePortsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreatePortsV41Api(client *core.CtyunClient) *CtecsCreatePortsV41Api {
	return &CtecsCreatePortsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/ports/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreatePortsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreatePortsV41Request) (*CtecsCreatePortsV41Response, error) {
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
	var resp CtecsCreatePortsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreatePortsV41Request struct {
	ClientToken             string   `json:"clientToken,omitempty"`             /*  客户端存根，用于保证订单幂等性。长度为1-64字符，要求单个云平台账户内唯一，使用同一个clientToken值，则代表为同一个请求。保留时间为24小时  */
	RegionID                string   `json:"regionID,omitempty"`                /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	SubnetID                string   `json:"subnetID,omitempty"`                /*  子网ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10028310">产品定义-子网</a>来了解子网<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=8659&data=94">查询子网列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4812&data=94">创建子网</a>  */
	PrimaryPrivateIp        string   `json:"primaryPrivateIp,omitempty"`        /*  弹性网卡的主私网IPv4地址  */
	Ipv6Addresses           []string `json:"ipv6Addresses"`                     /*  弹性网卡的主私网IPv6地址  */
	SecurityGroupIds        []string `json:"securityGroupIds"`                  /*  加入一个或多个安全组。安全组和弹性网卡必须在同一个专有网络VPC中。您可以查看<a href="https://www.ctyun.cn/document/10026755/10028520">安全组概述</a>来了解安全组<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4817&data=94">查询用户安全组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=4821&data=94">创建安全组</a>  */
	SecondaryPrivateIpCount int32    `json:"secondaryPrivateIpCount,omitempty"` /*  辅助私网IP地址数量，让ECS为您自动创建IP地址  */
	SecondaryPrivateIps     []string `json:"secondaryPrivateIps"`               /*  辅助私网IP地址，不能和secondaryPrivateIpCount同时指定  */
	Name                    string   `json:"name,omitempty"`                    /*  网卡名称，满足以下规则：支持拉丁字母、中文、数字，下划线，连字符，中文/英文字母开头，不能以http:/https:开头，长度2-32  */
	Description             string   `json:"description,omitempty"`             /*  网卡的描述，满足以下规则：支持拉丁字母、中文、数字, 特殊字符：\~!@#$%^&*()_-+= <>?:"{},./;'[\]·！@#￥%……&*（） —— -+={}\｜《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtecsCreatePortsV41Response struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为失败)  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsCreatePortsV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsCreatePortsV41ReturnObjResponse struct {
	VpcID                string   `json:"vpcID,omitempty"`                /*  vpc的ID  */
	SubnetID             string   `json:"subnetID,omitempty"`             /*  子网ID  */
	NetworkInterfaceID   string   `json:"networkInterfaceID,omitempty"`   /*  网卡ID  */
	NetworkInterfaceName string   `json:"networkInterfaceName,omitempty"` /*  网卡名称  */
	MacAddress           string   `json:"macAddress,omitempty"`           /*  mac地址  */
	Description          string   `json:"description,omitempty"`          /*  网卡描述  */
	Ipv6Address          []string `json:"ipv6Address"`                    /*  IPv6地址列表  */
	SecurityGroupIds     []string `json:"securityGroupIds"`               /*  安全组ID列表  */
	SecondaryPrivateIps  []string `json:"secondaryPrivateIps"`            /*  二级IP地址列表  */
	PrivateIpAddress     string   `json:"privateIpAddress,omitempty"`     /*  弹性网卡的主私有IP  */
	InstanceOwnerID      string   `json:"instanceOwnerID,omitempty"`      /*  绑定的实例的所有者ID  */
	InstanceType         string   `json:"instanceType,omitempty"`         /*  绑定的实例类型  */
	InstanceID           string   `json:"instanceID,omitempty"`           /*  绑定的实例ID  */
}
