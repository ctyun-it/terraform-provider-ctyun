package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// SubnetCreateApi 创建子网
// https://www.ctyun.cn/document/10026755/10040804
type SubnetCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSubnetCreateApi(client *ctyunsdk.CtyunClient) *SubnetCreateApi {
	return &SubnetCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/create-subnet",
		},
	}
}

func (this *SubnetCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SubnetCreateRequest) (*SubnetCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := subnetCreateRealRequest{
		RegionID:        req.RegionId,
		ClientToken:     req.ClientToken,
		Name:            req.Name,
		VpcID:           req.VpcId,
		CIDR:            req.Cidr,
		Description:     req.Description,
		EnableIpv6:      req.EnableIpv6,
		DnsList:         req.DnsList,
		SubnetGatewayIP: req.SubnetGatewayIp,
		SubnetType:      req.SubnetType,
		ProjectID:       req.ProjectId,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &subnetCreateRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}

	return &SubnetCreateResponse{
		SubnetId: result.SubnetID,
	}, nil
}

type subnetCreateRealRequest struct {
	RegionID        string   `json:"regionID"`
	ClientToken     string   `json:"clientToken"`
	Name            string   `json:"name"`
	VpcID           string   `json:"vpcID"`
	CIDR            string   `json:"CIDR"`
	Description     string   `json:"description"`
	EnableIpv6      bool     `json:"enableIpv6"`
	DnsList         []string `json:"dnsList"`
	SubnetGatewayIP string   `json:"subnetGatewayIP,omitempty"`
	SubnetType      string   `json:"subnetType"`
	ProjectID       string   `json:"projectID,omitempty"`
}

type subnetCreateRealResponse struct {
	SubnetID string `json:"subnetID"`
}

type SubnetCreateRequest struct {
	RegionId        string
	ClientToken     string
	Name            string
	VpcId           string
	Cidr            string
	Description     string
	EnableIpv6      bool
	DnsList         []string
	SubnetGatewayIp string
	SubnetType      string
	ProjectId       string
}

type SubnetCreateResponse struct {
	SubnetId string
}
