package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SubnetDeleteApi 删除子网
// https://www.ctyun.cn/document/10026755/10040807
type SubnetDeleteApi struct {
	builder ctyunsdk.CtyunRequestBuilder
	client  *ctyunsdk.CtyunClient
}

func NewSubnetDeleteApi(client *ctyunsdk.CtyunClient) *SubnetDeleteApi {
	return &SubnetDeleteApi{
		client: client,
		builder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/delete-subnet",
		},
	}
}

func (this *SubnetDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SubnetDeleteRequest) (*SubnetDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.builder.WithCredential(&credential)
	builder, err := builder.WriteJson(subnetDeleteRealRequest{
		RegionID:    req.RegionId,
		SubnetID:    req.SubnetId,
		ProjectID:   req.ProjectId,
		ClientToken: req.ClientToken,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	err = response.ParseByStandardModelWithCheck(nil)
	if err != nil {
		return nil, err
	}
	return &SubnetDeleteResponse{}, nil
}

type subnetDeleteRealRequest struct {
	RegionID    string `json:"regionID"`
	SubnetID    string `json:"subnetID"`
	ProjectID   string `json:"projectID,omitempty"`
	ClientToken string `json:"clientToken"`
}

type SubnetDeleteRequest struct {
	RegionId    string
	SubnetId    string
	ProjectId   string
	ClientToken string
}

type SubnetDeleteResponse struct {
}
