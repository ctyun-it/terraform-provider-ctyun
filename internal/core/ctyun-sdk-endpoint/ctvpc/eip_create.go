package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EipCreateApi 创建弹性 IP
// https://www.ctyun.cn/document/10026753/10040759
type EipCreateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipCreateApi(client *ctyunsdk.CtyunClient) *EipCreateApi {
	return &EipCreateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/create",
		},
	}
}

func (this *EipCreateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipCreateRequest) (*EipCreateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipCreateRealRequest{
		ClientToken:       req.ClientToken,
		RegionID:          req.RegionId,
		CycleType:         req.CycleType,
		CycleCount:        req.CycleCount,
		Name:              req.Name,
		Bandwidth:         req.Bandwidth,
		BandwidthID:       req.BandwidthId,
		DemandBillingType: req.DemandBillingType,
		ProjectID:         req.ProjectId,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &eipCreateRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}

	return &EipCreateResponse{
		MasterOrderId:        result.MasterOrderID,
		MasterOrderNo:        result.MasterOrderNO,
		MasterResourceId:     result.MasterResourceID,
		MasterResourceStatus: result.MasterResourceStatus,
		RegionId:             result.RegionID,
		EipId:                result.EipID,
	}, nil
}

type eipCreateRealRequest struct {
	ClientToken       string `json:"clientToken"`
	RegionID          string `json:"regionID"`
	CycleType         string `json:"cycleType"`
	CycleCount        int    `json:"cycleCount,omitempty"`
	Name              string `json:"name"`
	Bandwidth         int    `json:"bandwidth"`
	BandwidthID       string `json:"bandwidthID"`
	DemandBillingType string `json:"demandBillingType"`
	ProjectID         string `json:"projectID,omitempty"`
}

type eipCreateRealResponse struct {
	MasterOrderID        string `json:"masterOrderID"`
	MasterOrderNO        string `json:"masterOrderNO"`
	MasterResourceID     string `json:"masterResourceID"`
	MasterResourceStatus string `json:"masterResourceStatus"`
	RegionID             string `json:"regionID"`
	EipID                string `json:"eipID"`
}

type EipCreateRequest struct {
	ClientToken       string
	RegionId          string
	CycleType         string
	CycleCount        int
	Name              string
	Bandwidth         int
	BandwidthId       string
	DemandBillingType string
	ProjectId         string
}

type EipCreateResponse struct {
	MasterOrderId        string
	MasterOrderNo        string
	MasterResourceId     string
	MasterResourceStatus string
	RegionId             string
	EipId                string
}
