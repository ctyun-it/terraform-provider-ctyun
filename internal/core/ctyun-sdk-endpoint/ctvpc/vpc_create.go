package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// VpcCreateApi 创建VPC
// https://www.ctyun.cn/document/10026755/10040800
type VpcCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewVpcCreateApi(client *ctyunsdk.CtyunClient) *VpcCreateApi {
	return &VpcCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/create",
		},
	}
}

func (this *VpcCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *VpcCreateRequest) (*VpcCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	realRequest := vpcCreateRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		CIDR:        req.Cidr,
		Name:        req.Name,
		Description: req.Description,
		EnableIpv6:  req.EnableIpv6,
		ProjectID:   req.ProjectId,
	}
	_, err := builder.WriteJson(realRequest)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	realResponse := &vpcCreateRealResponse{}
	err = response.ParseByStandardModelWithCheck(realResponse)
	if err != nil {
		return nil, err
	}
	return &VpcCreateResponse{
		VpcId: realResponse.VpcID,
	}, nil
}

type vpcCreateRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	CIDR        string `json:"CIDR"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EnableIpv6  bool   `json:"enableIpv6"`
	ProjectID   string `json:"projectID,omitempty"`
}

type vpcCreateRealResponse struct {
	VpcID string `json:"vpcID"`
}

type VpcCreateRequest struct {
	RegionId    string // 资源池id
	ClientToken string // 客户端存根，用于保证订单幂等性, 长度 1 - 64
	Name        string // 虚拟私有云名称
	Cidr        string // VPC的网段。建议您使用 192.168.0.0/16、172.16.0.0/12、10.0.0.0/8 三个 RFC 标准私网网段及其子网作为专有网络的主 IPv4 网段，网段掩码有效范围为 8~28 位
	Description string // 描述
	EnableIpv6  bool   // 是否开启 IPv6 网段。取值：false（默认值）:不开启，true: 开启
	ProjectId   string // 企业项目id
}

type VpcCreateResponse struct {
	VpcId string
}
