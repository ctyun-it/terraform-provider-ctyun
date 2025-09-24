package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsShowPortsV41Api
/* 查询网卡信息<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsShowPortsV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsShowPortsV41Api(client *core.CtyunClient) *CtecsShowPortsV41Api {
	return &CtecsShowPortsV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/ports/show",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsShowPortsV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsShowPortsV41Request) (*CtecsShowPortsV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("networkInterfaceID", req.NetworkInterfaceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsShowPortsV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsShowPortsV41Request struct {
	RegionID           string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	NetworkInterfaceID string /*  网卡ID，您可以查看<a href="https://www.ctyun.cn/document/10026755/10197673">弹性网卡-弹性网卡基本知识</a>来了解弹性网卡<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5802&data=94">查询弹性网卡列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=18&api=5789&data=94">创建弹性网卡</a>  */
}

type CtecsShowPortsV41Response struct {
	StatusCode  int32                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                              `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                              `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                              `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                              `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsShowPortsV41ReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtecsShowPortsV41ReturnObjResponse struct {
	NetworkInterfaceName string                                           `json:"networkInterfaceName,omitempty"` /*  虚拟网名称  */
	NetworkInterfaceID   string                                           `json:"networkInterfaceID,omitempty"`   /*  虚拟网id  */
	VpcID                string                                           `json:"vpcID,omitempty"`                /*  所属vpc  */
	SubnetID             string                                           `json:"subnetID,omitempty"`             /*  所属子网id  */
	Role                 int32                                            `json:"role,omitempty"`                 /*  网卡类型: 0 主网卡， 1 弹性网卡  */
	MacAddress           string                                           `json:"macAddress,omitempty"`           /*  mac地址  */
	PrimaryPrivateIp     string                                           `json:"primaryPrivateIp,omitempty"`     /*  主ip  */
	Ipv6Addresses        []string                                         `json:"ipv6Addresses"`                  /*  ipv6地址  */
	InstanceID           string                                           `json:"instanceID,omitempty"`           /*  关联的设备id  */
	InstanceType         string                                           `json:"instanceType,omitempty"`         /*  设备类型 VM, BM, Other  */
	Description          string                                           `json:"description,omitempty"`          /*  描述  */
	SecurityGroupIds     []string                                         `json:"securityGroupIds"`               /*  安全组ID列表  */
	SecondaryPrivateIps  []string                                         `json:"secondaryPrivateIps"`            /*  辅助私网IP  */
	AdminStatus          string                                           `json:"adminStatus,omitempty"`          /*  是否启用DOWN, UP  */
	AssociatedEip        *CtecsShowPortsV41ReturnObjAssociatedEipResponse `json:"associatedEip"`                  /*  关联的eip信息  */
}

type CtecsShowPortsV41ReturnObjAssociatedEipResponse struct {
	Id   string `json:"id,omitempty"`   /*  弹性公网IP的ID  */
	Name string `json:"name,omitempty"` /*  弹性公网IP的名称  */
}
