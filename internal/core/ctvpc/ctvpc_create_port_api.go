package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcCreatePortApi
/* 创建弹性网卡
 */type CtvpcCreatePortApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcCreatePortApi(client *core.CtyunClient) *CtvpcCreatePortApi {
	return &CtvpcCreatePortApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ports/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcCreatePortApi) Do(ctx context.Context, credential core.Credential, req *CtvpcCreatePortRequest) (*CtvpcCreatePortResponse, error) {
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
	var resp CtvpcCreatePortResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcCreatePortRequest struct {
	ClientToken             string    `json:"clientToken,omitempty"`      /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID                string    `json:"regionID,omitempty"`         /*  资源池ID  */
	SubnetID                string    `json:"subnetID,omitempty"`         /*  子网ID  */
	PrimaryPrivateIp        *string   `json:"primaryPrivateIp,omitempty"` /*  弹性网卡的主私有IP地址  */
	Ipv6Addresses           []*string `json:"ipv6Addresses"`              /*  为弹性网卡指定一个或多个IPv6地址  */
	SecurityGroupIds        []*string `json:"securityGroupIds"`           /*  加入一个或多个安全组。安全组和弹性网卡必须在同一个专有网络VPC中，最多同时支持 10 个  */
	SecondaryPrivateIpCount int32     `json:"secondaryPrivateIpCount"`    /*  指定私有IP地址数量，让ECS为您自动创建IP地址  */
	SecondaryPrivateIps     []*string `json:"secondaryPrivateIps"`        /*  指定私有IP地址，不能和secondaryPrivateIpCount同时指定  */
	Name                    *string   `json:"name,omitempty"`             /*  支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32  */
	Description             *string   `json:"description,omitempty"`      /*  支持拉丁字母、中文、数字, 特殊字符：~!@#$%^&*()_-+= <>?:{},./;'[]·~！@#￥%……&*（） —— -+={}\《》？：“”【】、；‘'，。、，不能以 http: / https: 开头，长度 0 - 128  */
}

type CtvpcCreatePortResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcCreatePortReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcCreatePortReturnObjResponse struct {
	VpcID                *string   `json:"vpcID,omitempty"`                /*  vpc的id  */
	SubnetID             *string   `json:"subnetID,omitempty"`             /*  子网id  */
	NetworkInterfaceID   *string   `json:"networkInterfaceID,omitempty"`   /*  网卡id  */
	NetworkInterfaceName *string   `json:"networkInterfaceName,omitempty"` /*  网卡名称  */
	MacAddress           *string   `json:"macAddress,omitempty"`           /*  mac地址  */
	Description          *string   `json:"description,omitempty"`          /*  网卡描述  */
	Ipv6Address          []*string `json:"ipv6Address"`                    /*  IPv6地址列表  */
	SecurityGroupIds     []*string `json:"securityGroupIds"`               /*  安全组ID列表  */
	SecondaryPrivateIps  []*string `json:"secondaryPrivateIps"`            /*  二级IP地址列表  */
	PrivateIpAddress     *string   `json:"privateIpAddress,omitempty"`     /*  弹性网卡的主私有IP  */
	InstanceOwnerID      *string   `json:"instanceOwnerID,omitempty"`      /*  绑定的实例的所有者ID  */
	InstanceType         *string   `json:"instanceType,omitempty"`         /*  设备类型 VM, BM, Other  */
	InstanceID           *string   `json:"instanceID,omitempty"`           /*  绑定的实例ID  */
}
