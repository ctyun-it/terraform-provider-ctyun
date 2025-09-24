package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
	"strings"
)

// SubnetListApi 查询子网
// https://www.ctyun.cn/document/10026755/10040797
type SubnetListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSubnetListApi(client *ctyunsdk.CtyunClient) *SubnetListApi {
	return &SubnetListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/vpc/new-list-subnet",
		},
	}
}

func (this *SubnetListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SubnetListRequest) (*SubnetListResponse, ctyunsdk.CtyunRequestError) {
	builder := this.
		WithCredential(&credential).
		AddParam("regionID", req.RegionId).
		AddParam("vpcID", req.VpcId).
		AddParam("subnetID", strings.Join(req.SubnetIds, ",")).
		AddParam("pageNumber", strconv.Itoa(req.PageNumber)).
		AddParam("pageSize", strconv.Itoa(req.PageSize))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	realResponse := &subnetListRealResponse{}
	err = response.ParseByStandardModelWithCheck(realResponse)
	if err != nil {
		return nil, err
	}
	var subnets []SubnetListSubnetsResponse
	for _, s := range realResponse.Subnets {
		subnets = append(subnets, SubnetListSubnetsResponse{
			SubnetId:          s.SubnetID,
			Name:              s.Name,
			Description:       s.Description,
			VpcId:             s.VpcID,
			Cidr:              s.CIDR,
			AvailableIpCount:  s.AvailableIPCount,
			GatewayIp:         s.GatewayIP,
			AvailabilityZones: s.AvailabilityZones,
			RouteTableId:      s.RouteTableID,
			NetworkAclId:      s.NetworkAclID,
			Start:             s.Start,
			End:               s.End,
			Ipv6Enabled:       s.Ipv6Enabled,
			Ipv6Cidr:          s.Ipv6CIDR,
			Ipv6Start:         s.Ipv6Start,
			Ipv6End:           s.Ipv6End,
			Ipv6GatewayIp:     s.Ipv6GatewayIP,
			DnsList:           s.DnsList,
			NtpList:           s.NtpList,
			Type:              s.Type,
			CreateAt:          s.CreateAt,
			UpdateAt:          s.UpdateAt,
		})
	}
	return &SubnetListResponse{
		Subnets:      subnets,
		CurrentCount: realResponse.CurrentCount,
		TotalPage:    realResponse.TotalPage,
	}, nil
}

type subnetListRealResponse struct {
	Subnets []struct {
		SubnetID          string   `json:"subnetID"`
		Name              string   `json:"name"`
		Description       string   `json:"description"`
		VpcID             string   `json:"vpcID"`
		CIDR              string   `json:"CIDR"`
		AvailableIPCount  int      `json:"availableIPCount"`
		GatewayIP         string   `json:"gatewayIP"`
		AvailabilityZones []string `json:"availabilityZones"`
		RouteTableID      string   `json:"routeTableID"`
		NetworkAclID      string   `json:"networkAclID"`
		Start             string   `json:"start"`
		End               string   `json:"end"`
		Ipv6Enabled       int      `json:"ipv6Enabled"`
		EnableIpv6        bool     `json:"enableIpv6"`
		Ipv6CIDR          string   `json:"ipv6CIDR"`
		Ipv6Start         string   `json:"ipv6Start"`
		Ipv6End           string   `json:"ipv6End"`
		Ipv6GatewayIP     string   `json:"ipv6GatewayIP"`
		DnsList           []string `json:"dnsList"`
		NtpList           []string `json:"ntpList"`
		Type              int      `json:"type"`
		CreateAt          string   `json:"createAt"`
		UpdateAt          string   `json:"updateAt"`
	} `json:"subnets"`
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
}

type SubnetListRequest struct {
	RegionId   string   // 资源池id
	VpcId      string   // 查询的vpcId
	SubnetIds  []string // 查询的vpcSubnetId
	PageNumber int      // 列表的页码，默认值为 1。
	PageSize   int      // 分页查询时每页的行数，最大值为 50，默认值为 10。
}

type SubnetListSubnetsResponse struct {
	SubnetId          string
	Name              string
	Description       string
	VpcId             string
	Cidr              string
	AvailableIpCount  int
	GatewayIp         string
	AvailabilityZones []string
	RouteTableId      string
	NetworkAclId      string
	Start             string
	End               string
	Ipv6Enabled       int
	EnableIpv6        bool
	Ipv6Cidr          string
	Ipv6Start         string
	Ipv6End           string
	Ipv6GatewayIp     string
	DnsList           []string
	NtpList           []string
	Type              int
	CreateAt          string
	UpdateAt          string
}

type SubnetListResponse struct {
	Subnets      []SubnetListSubnetsResponse
	CurrentCount int
	TotalCount   int
	TotalPage    int
}
