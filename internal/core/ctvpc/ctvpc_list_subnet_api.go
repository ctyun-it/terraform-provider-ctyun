package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListSubnetApi
/* 查询用户专有网络下子网列表
 */type CtvpcListSubnetApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListSubnetApi(client *core.CtyunClient) *CtvpcListSubnetApi {
	return &CtvpcListSubnetApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/list-subnet",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListSubnetApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListSubnetRequest) (*CtvpcListSubnetResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ClientToken != nil {
		ctReq.AddParam("clientToken", *req.ClientToken)
	}
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.SubnetID != nil {
		ctReq.AddParam("subnetID", *req.SubnetID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListSubnetResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListSubnetRequest struct {
	ClientToken *string /*  客户端存根，用于保证订单幂等性, 长度 1 - 64  */
	RegionID    string  /*  资源池 ID  */
	VpcID       *string /*  VPC 的 ID  */
	SubnetID    *string /*  多个 subnet 的 ID 之间用半角逗号（,）隔开。  */
	PageNumber  int32   /*  列表的页码，默认值为 1。  */
	PageNo      int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize    int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListSubnetResponse struct {
	StatusCode   int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    *CtvpcListSubnetReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	TotalCount   int32                             `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                             `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                             `json:"totalPage"`             /*  总页数  */
}

type CtvpcListSubnetReturnObjResponse struct {
	Subnets []*CtvpcListSubnetReturnObjSubnetsResponse `json:"subnets"` /*  subnets 组  */
}

type CtvpcListSubnetReturnObjSubnetsResponse struct {
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
	Ipv6Enabled       int32     `json:"ipv6Enabled"`             /*  是否配置了ipv6网段  */
	EnableIpv6        *bool     `json:"enableIpv6"`              /*  是否开启 ipv6  */
	Ipv6CIDR          *string   `json:"ipv6CIDR,omitempty"`      /*  子网 Ipv6 网段，掩码范围为 16-28 位  */
	Ipv6Start         *string   `json:"ipv6Start,omitempty"`     /*  子网内可用的起始 IPv6 地址  */
	Ipv6End           *string   `json:"ipv6End,omitempty"`       /*  子网内可用的结束 IPv6 地址  */
	Ipv6GatewayIP     *string   `json:"ipv6GatewayIP,omitempty"` /*  v6 网关地址  */
	DnsList           []*string `json:"dnsList"`                 /*  DNS 服务器地址:默认为空；必须为正确的 IPv4 格式；重新触发 DHCP 后生效，最大数组长度为 4  */
	NtpList           []*string `json:"ntpList"`                 /*  NTP 服务器地址: 默认为空，必须为正确的域名或 IPv4 格式；重新触发 DHCP 后生效，最大数组长度为 4  */
	RawType           int32     `json:"type"`                    /*  子网类型 :当前仅支持：0（普通子网）, 1（裸金属子网）  */
	UpdatedAt         *string   `json:"updatedAt,omitempty"`     /*  更新时间  */
}
