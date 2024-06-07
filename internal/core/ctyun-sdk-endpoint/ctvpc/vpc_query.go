package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// VpcQueryApi 查询VPC
// https://www.ctyun.cn/document/10026755/10040783
type VpcQueryApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewVpcQueryApi(client *ctyunsdk.CtyunClient) *VpcQueryApi {
	return &VpcQueryApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/vpc/query",
		},
	}
}

func (this *VpcQueryApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *VpcQueryRequest) (*VpcQueryResponse, ctyunsdk.CtyunRequestError) {
	builder := this.
		WithCredential(&credential).
		AddParam("clientToken", req.ClientToken).
		AddParam("regionID", req.RegionId).
		AddParam("projectID", req.ProjectId).
		AddParam("vpcID", req.VpcId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	resp := &vpcQueryRealResponse{}
	err = response.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}
	return &VpcQueryResponse{
		VpcId:          resp.VpcID,
		Name:           resp.Name,
		Description:    resp.Description,
		Cidr:           resp.CIDR,
		Ipv6Enabled:    resp.Ipv6Enabled,
		Ipv6Cidrs:      resp.Ipv6CIDRS,
		SubnetIds:      resp.SubnetIDs,
		NatGatewayIds:  resp.NatGatewayIDs,
		SecondaryCidrs: resp.SecondaryCIDRs,
	}, nil
}

type vpcQueryRealResponse struct {
	VpcID          string   `json:"vpcID"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	CIDR           string   `json:"CIDR"`
	Ipv6Enabled    bool     `json:"ipv6Enabled"`
	Ipv6CIDRS      []string `json:"ipv6CIDRS"`
	SubnetIDs      []string `json:"subnetIDs"`
	NatGatewayIDs  []string `json:"natGatewayIDs"`
	SecondaryCIDRs []string `json:"secondaryCIDRs"`
}

type VpcQueryRequest struct {
	ClientToken string
	RegionId    string
	ProjectId   string
	VpcId       string
}

type VpcQueryResponse struct {
	VpcId          string
	Name           string
	Description    string
	Cidr           string
	Ipv6Enabled    bool
	Ipv6Cidrs      []string
	SubnetIds      []string
	NatGatewayIds  []string
	SecondaryCidrs []string
}
