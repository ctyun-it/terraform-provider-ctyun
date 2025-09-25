package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// BandwidthChangeSpecApi 修改共享带宽的带宽峰值
// https://www.ctyun.cn/document/10026761/10040813
type BandwidthChangeSpecApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthChangeSpecApi(client *ctyunsdk.CtyunClient) *BandwidthChangeSpecApi {
	return &BandwidthChangeSpecApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/bandwidth/modify-spec",
		},
	}
}

func (this *BandwidthChangeSpecApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthChangeSpecRequest) (*BandwidthChangeSpecResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&bandwidthChangeSpecRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		BandwidthID: req.BandwidthId,
		Bandwidth:   req.Bandwidth,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	response := &bandwidthChangeSpecRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &BandwidthChangeSpecResponse{
		MasterOrderId: response.MasterOrderID,
		MasterOrderNo: response.MasterOrderNO,
		RegionId:      response.RegionID,
	}, nil
}

type bandwidthChangeSpecRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	BandwidthID string `json:"bandwidthID"`
	Bandwidth   int    `json:"bandwidth"`
}

type bandwidthChangeSpecRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
}

type BandwidthChangeSpecRequest struct {
	RegionId    string
	ClientToken string
	BandwidthId string
	Bandwidth   int
}

type BandwidthChangeSpecResponse struct {
	MasterOrderId string
	MasterOrderNo string
	RegionId      string
}
