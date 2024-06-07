package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// BandwidthDisassociateEipApi 从共享带宽移出EIP
// https://www.ctyun.cn/document/10026761/10040772
type BandwidthDisassociateEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewBandwidthDisassociateEipApi(client *ctyunsdk.CtyunClient) *BandwidthDisassociateEipApi {
	return &BandwidthDisassociateEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/bandwidth/disassociate-eip",
		},
	}
}

func (this *BandwidthDisassociateEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *BandwidthDisassociateEipRequest) (*BandwidthDisassociateEipResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := bandwidthDisassociateEipRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		EipIds:      req.EipIds,
		BandwidthID: req.BandwidthId,
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

	return &BandwidthDisassociateEipResponse{}, nil
}

type bandwidthDisassociateEipRealRequest struct {
	RegionID    string   `json:"regionID"`
	ClientToken string   `json:"clientToken"`
	BandwidthID string   `json:"bandwidthID"`
	EipIds      []string `json:"eipIDs"`
}

type BandwidthDisassociateEipRequest struct {
	RegionId    string
	ClientToken string
	EipIds      []string
	BandwidthId string
}

type BandwidthDisassociateEipResponse struct {
}
