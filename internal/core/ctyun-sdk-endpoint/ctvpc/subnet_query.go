package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"time"
)

// SubnetQueryApi 查询子网详情
// https://www.ctyun.cn/document/10026755/10040792
type SubnetQueryApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSubnetQueryApi(client *ctyunsdk.CtyunClient) *SubnetQueryApi {
	return &SubnetQueryApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/vpc/query-subnet",
		},
	}
}

func (this *SubnetQueryApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SubnetQueryRequest) (*SubnetQueryResponse, ctyunsdk.CtyunRequestError) {
	builder := this.
		WithCredential(&credential).
		AddParam("clientToken", req.ClientToken).
		AddParam("regionID", req.RegionId).
		AddParam("projectID", req.ProjectId).
		AddParam("subnetID", req.SubnetId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	realResponse := &subnetQueryRealResponse{}
	err = response.ParseByStandardModelWithCheck(realResponse)
	if err != nil {
		return nil, err
	}
	return &SubnetQueryResponse{
		SubnetId:          realResponse.SubnetID,
		Name:              realResponse.Name,
		Description:       realResponse.Description,
		VpcId:             realResponse.VpcID,
		AvailabilityZones: realResponse.AvailabilityZones,
		RouteTableId:      realResponse.RouteTableID,
		NetworkAclId:      realResponse.NetworkAclID,
		Cidr:              realResponse.CIDR,
		Gateway:           realResponse.Gateway,
		Start:             realResponse.Start,
		End:               realResponse.End,
		Ipv6Enabled:       realResponse.Ipv6Enabled,
		EnableIpv6:        realResponse.EnableIpv6,
		AvailableIpCount:  realResponse.AvailableIpCount,
		Ipv6Cidr:          realResponse.Ipv6CIDR,
		Ipv6Start:         realResponse.Ipv6Start,
		Ipv6End:           realResponse.Ipv6End,
		Ipv6GatewayIp:     realResponse.Ipv6GatewayIP,
		DnsList:           realResponse.DnsList,
		NtpList:           realResponse.NtpList,
		Type:              realResponse.Type,
		CreatedAt:         realResponse.CreatedAt,
		UpdatedAt:         realResponse.UpdatedAt,
	}, nil
}

type subnetQueryRealResponse struct {
	SubnetID          string    `json:"subnetID"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	VpcID             string    `json:"vpcID"`
	AvailabilityZones []string  `json:"availabilityZones"`
	RouteTableID      string    `json:"routeTableID"`
	NetworkAclID      string    `json:"networkAclID"`
	CIDR              string    `json:"CIDR"`
	Gateway           string    `json:"gateway"`
	Start             string    `json:"start"`
	End               string    `json:"end"`
	AvailableIpCount  int       `json:"availableIpCount"`
	Ipv6Enabled       int       `json:"ipv6Enabled"`
	EnableIpv6        bool      `json:"enableIpv6"`
	Ipv6CIDR          string    `json:"ipv6CIDR"`
	Ipv6Start         string    `json:"ipv6Start"`
	Ipv6End           string    `json:"ipv6End"`
	Ipv6GatewayIP     string    `json:"ipv6GatewayIP"`
	DnsList           []string  `json:"dnsList"`
	NtpList           []string  `json:"ntpList"`
	Type              int       `json:"type"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type SubnetQueryRequest struct {
	ClientToken string
	RegionId    string
	ProjectId   string
	SubnetId    string
}

type SubnetQueryResponse struct {
	SubnetId          string
	Name              string
	Description       string
	VpcId             string
	AvailabilityZones []string
	RouteTableId      string
	NetworkAclId      string
	Cidr              string
	Gateway           string
	Start             string
	End               string
	Ipv6Enabled       int
	EnableIpv6        bool
	AvailableIpCount  int
	Ipv6Cidr          string
	Ipv6Start         string
	Ipv6End           string
	Ipv6GatewayIp     string
	DnsList           []string
	NtpList           []string
	Type              int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
