package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// SubnetUpdateApi 修改子网属性
// https://www.ctyun.cn/document/10026755/10040815
type SubnetUpdateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewSubnetUpdateApi(client *ctyunsdk.CtyunClient) *SubnetUpdateApi {
	return &SubnetUpdateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/vpc/update-subnet",
		},
	}
}

func (this *SubnetUpdateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *SubnetUpdateRequest) (*SubnetUpdateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := subnetUpdateRealRequest{
		RegionID:    req.RegionId,
		SubnetID:    req.SubnetId,
		Name:        req.Name,
		Description: req.Description,
		DnsList:     req.DnsList,
	}
	_, err := builder.WriteJson(requestContent)
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
	return &SubnetUpdateResponse{}, nil
}

type subnetUpdateRealRequest struct {
	RegionID    string   `json:"regionID"`
	SubnetID    string   `json:"subnetID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DnsList     []string `json:"dnsList"`
}

type SubnetUpdateRequest struct {
	RegionId    string
	SubnetId    string
	Name        string
	Description string
	DnsList     []string
}

type SubnetUpdateResponse struct {
}
