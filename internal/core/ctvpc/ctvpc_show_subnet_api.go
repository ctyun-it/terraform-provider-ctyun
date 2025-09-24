package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowSubnetApi
/* 查询用户专有网络 VPC 下子网详情。
 */type CtvpcShowSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowSubnetApi(client *core.CtyunClient) *CtvpcShowSubnetApi {
	return &CtvpcShowSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/query-subnet",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowSubnetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowSubnetRequest) (*CtvpcShowSubnetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("subnetID", req.SubnetID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowSubnetRequest struct {
	RegionID string /*  资源池 ID  */
	SubnetID string /*  subnet 的 ID  */
}

type CtvpcShowSubnetResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowSubnetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcShowSubnetReturnObjResponse struct {
	SubnetID          *string   `json:"subnetID,omitempty"`      /*  subnet ID  */
	Name              *string   `json:"name,omitempty"`          /*  名称  */
	Description       *string   `json:"description,omitempty"`   /*  描述  */
	VpcID             *string   `json:"vpcID,omitempty"`         /*  VpcID  */
	AvailabilityZones []*string `json:"availabilityZones"`       /*  子网所在的可用区名  */
	RouteTableID      *string   `json:"routeTableID,omitempty"`  /*  子网路由表 ID  */
	NetworkAclID      *string   `json:"networkAclID,omitempty"`  /*  子网 aclID  */
	CIDR              *string   `json:"CIDR,omitempty"`          /*  子网网段，掩码范围为 16-28 位  */
	GatewayIP         *string   `json:"gatewayIP,omitempty"`     /*  子网网关  */
	DhcpIP            *string   `json:"dhcpIP,omitempty"`        /*  dhcpIP  */
	Start             *string   `json:"start,omitempty"`         /*  子网网段起始 IP  */
	End               *string   `json:"end,omitempty"`           /*  子网网段结束 IP  */
	AvailableIPCount  int32     `json:"availableIPCount"`        /*  子网内可用 IPv4 数目  */
	Ipv6Enabled       int32     `json:"ipv6Enabled"`             /*  是否配置了ipv6网段，1 表示开启，0 表示未开启  */
	EnableIpv6        *bool     `json:"enableIpv6"`              /*  是否开启 ipv6  */
	Ipv6CIDR          *string   `json:"ipv6CIDR,omitempty"`      /*  子网 Ipv6 网段，掩码范围为 16-28 位  */
	Ipv6Start         *string   `json:"ipv6Start,omitempty"`     /*  子网内可用的起始 IPv6 地址  */
	Ipv6End           *string   `json:"ipv6End,omitempty"`       /*  子网内可用的结束 IPv6 地址  */
	Ipv6GatewayIP     *string   `json:"ipv6GatewayIP,omitempty"` /*  v6 网关地址  */
	DnsList           []*string `json:"dnsList"`                 /*  DNS 服务器地址:默认为空；必须为正确的 IPv4 格式；重新触发 DHCP 后生效，最大数组长度为 4  */
	NtpList           []*string `json:"ntpList"`                 /*  NTP 服务器地址: 默认为空，必须为正确的域名或 IPv4 格式；重新触发 DHCP 后生效，最大数组长度为 4  */
	RawType           int32     `json:"type"`                    /*  子网类型 :当前仅支持：0（普通子网）, 1（裸金属子网）  */
	CreatedAt         *string   `json:"createdAt,omitempty"`     /*  创建时间  */
	UpdatedAt         *string   `json:"updatedAt,omitempty"`     /*  更新时间  */
}
