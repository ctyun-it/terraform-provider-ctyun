package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// VpcDeleteApi 删除VPC
// https://www.ctyun.cn/document/10026755/10040805
type VpcDeleteApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewVpcDeleteApi(client *ctyunsdk.CtyunClient) *VpcDeleteApi {
	return &VpcDeleteApi{
		client: client,
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/delete",
		},
	}
}

func (this *VpcDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *VpcDeleteRequest) (*VpcDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder, err := this.builder.WithCredential(&credential).WriteJson(vpcDeleteRealRequest{
		RegionID:    req.RegionId,
		VpcID:       req.VpcId,
		ClientToken: req.ClientToken,
		ProjectID:   req.ProjectId,
	})

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &VpcDeleteResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	return result, err
}

type vpcDeleteRealRequest struct {
	RegionID    string `json:"regionID"`
	VpcID       string `json:"vpcID"`
	ClientToken string `json:"clientToken"`
	ProjectID   string `json:"projectID,omitempty"`
}

type VpcDeleteRequest struct {
	ClientToken string // 客户端存根，用于保证订单幂等性。要求单个云平台账户内唯一
	RegionId    string // 资源池ID
	VpcId       string // VPC的ID
	ProjectId   string // 企业项目 ID，默认为"0"
}

type VpcDeleteResponse struct {
}
