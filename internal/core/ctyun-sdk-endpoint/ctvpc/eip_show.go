package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EipShowApi 查看弹性 IP 详情
// https://www.ctyun.cn/document/10026753/10042983
type EipShowApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipShowApi(client *ctyunsdk.CtyunClient) *EipShowApi {
	return &EipShowApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/eip/show",
		},
	}
}

func (this *EipShowApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipShowRequest) (*EipShowResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	builder.AddParam("regionID", req.RegionId)
	builder.AddParam("eipID", req.EipId)

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &eipShowRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}
	return &EipShowResponse{
		Id:               result.ID,
		Name:             result.Name,
		EipAddress:       result.EipAddress,
		AssociationId:    result.AssociationID,
		AssociationType:  result.AssociationType,
		PrivateIpAddress: result.PrivateIpAddress,
		Bandwidth:        result.Bandwidth,
		BandwidthId:      result.BandwidthID,
		BandwidthType:    result.BandwidthType,
		Status:           result.Status,
		Tags:             result.Tags,
		CreatedAt:        result.CreatedAt,
		UpdatedAt:        result.UpdatedAt,
		ExpiredAt:        result.ExpiredAt,
	}, nil
}

type eipShowRealResponse struct {
	ID               string `json:"ID"`
	Name             string `json:"name"`
	EipAddress       string `json:"eipAddress"`
	AssociationID    string `json:"associationID"`
	AssociationType  string `json:"associationType"`
	PrivateIpAddress string `json:"privateIpAddress"`
	Bandwidth        int    `json:"bandwidth"`
	BandwidthID      string `json:"bandwidthID"`
	BandwidthType    string `json:"bandwidthType"`
	Status           string `json:"status"`
	Tags             string `json:"tags"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	ExpiredAt        string `json:"expiredAt"`
}

type EipShowRequest struct {
	RegionId string
	EipId    string
}

type EipShowResponse struct {
	Id               string
	Name             string
	EipAddress       string
	AssociationId    string
	AssociationType  string
	PrivateIpAddress string
	Bandwidth        int
	BandwidthId      string
	BandwidthType    string
	Status           string
	Tags             string
	CreatedAt        string
	UpdatedAt        string
	ExpiredAt        string
}
