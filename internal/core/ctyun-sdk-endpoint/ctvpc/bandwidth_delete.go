package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// BandwidthDeleteApi 删除共享带宽
// https://www.ctyun.cn/document/10026761/10040806
type BandwidthDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthDeleteApi(client *ctyunsdk.CtyunClient) *BandwidthDeleteApi {
	return &BandwidthDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/bandwidth/delete",
		},
	}
}

func (this *BandwidthDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthDeleteRequest) (*BandwidthDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&bandwidthDeleteRealRequest{
		RegionID:    req.RegionId,
		BandwidthID: req.BandwidthId,
		ClientToken: req.ClientToken,
		ProjectID:   req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	response := &bandwidthDeleteRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &BandwidthDeleteResponse{
		MasterOrderId: response.MasterOrderID,
		MasterOrderNo: response.MasterOrderNO,
		RegionId:      response.RegionID,
	}, nil
}

type bandwidthDeleteRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	BandwidthID string `json:"bandwidthID"`
	ProjectID   string `json:"projectID,omitempty"`
}

type bandwidthDeleteRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
}

type BandwidthDeleteRequest struct {
	RegionId    string
	ClientToken string
	BandwidthId string
	ProjectId   string
}

type BandwidthDeleteResponse struct {
	MasterOrderId string
	MasterOrderNo string
	RegionId      string
}
