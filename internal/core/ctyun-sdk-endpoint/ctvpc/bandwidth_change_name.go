package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// BandwidthChangeNameApi 修改共享带宽名称描述等
// https://www.ctyun.cn/document/10026761/10040814
type BandwidthChangeNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthChangeNameApi(client *ctyunsdk.CtyunClient) *BandwidthChangeNameApi {
	return &BandwidthChangeNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/bandwidth/modify-attribute",
		},
	}
}

func (this *BandwidthChangeNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthChangeNameRequest) (*BandwidthChangeNameResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&bandwidthChangeNameRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		BandwidthID: req.BandwidthId,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	response := &bandwidthChangeNameRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &BandwidthChangeNameResponse{}, nil
}

type bandwidthChangeNameRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	BandwidthID string `json:"bandwidthID"`
	Name        string `json:"name"`
}

type bandwidthChangeNameRealResponse struct {
}

type BandwidthChangeNameRequest struct {
	RegionId    string
	ClientToken string
	BandwidthId string
	Name        string
}

type BandwidthChangeNameResponse struct {
}
