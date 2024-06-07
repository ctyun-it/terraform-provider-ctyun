package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// BandwidthCreateApi 创建共享带宽
// https://www.ctyun.cn/document/10026761/10040771
type BandwidthCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthCreateApi(client *ctyunsdk.CtyunClient) *BandwidthCreateApi {
	return &BandwidthCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/bandwidth/create",
		},
	}
}

func (this *BandwidthCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthCreateRequest) (*BandwidthCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(bandwidthCreateRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		CycleType:   req.CycleType,
		Bandwidth:   req.Bandwidth,
		CycleCount:  req.CycleCount,
		Name:        req.Name,
		ProjectID:   req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	response := &bandwidthCreateRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}

	return &BandwidthCreateResponse{
		MasterOrderId:        response.MasterOrderID,
		MasterOrderNo:        response.MasterOrderNO,
		MasterResourceId:     response.MasterResourceID,
		MasterResourceStatus: response.MasterResourceStatus,
		RegionId:             response.RegionID,
		BandwidthId:          response.BandwidthId,
	}, nil
}

type bandwidthCreateRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	CycleType   string `json:"cycleType"`
	Bandwidth   int64  `json:"bandwidth"`
	CycleCount  int64  `json:"cycleCount"`
	Name        string `json:"name"`
	ProjectID   string `json:"projectID,omitempty"`
}

type bandwidthCreateRealResponse struct {
	MasterOrderID        string `json:"masterOrderID"`
	MasterOrderNO        string `json:"masterOrderNO"`
	MasterResourceID     string `json:"masterResourceID"`
	MasterResourceStatus string `json:"masterResourceStatus"`
	RegionID             string `json:"regionID"`
	BandwidthId          string `json:"bandwidthId"`
}

type BandwidthCreateRequest struct {
	RegionId    string
	ClientToken string
	CycleType   string
	Bandwidth   int64
	CycleCount  int64
	Name        string
	ProjectId   string
}

type BandwidthCreateResponse struct {
	MasterOrderId        string
	MasterOrderNo        string
	MasterResourceId     string
	MasterResourceStatus string
	RegionId             string
	BandwidthId          string
}
