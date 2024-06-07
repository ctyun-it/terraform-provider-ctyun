package ctebs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EbsDeleteApi 云硬盘退订
// https://www.ctyun.cn/document/10027696/10110701
type EbsDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsDeleteApi(client *ctyunsdk.CtyunClient) *EbsDeleteApi {
	return &EbsDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/refund-ebs",
		},
	}
}

func (this *EbsDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsDeleteRequest) (*EbsDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ebsDeleteRealRequest{
		RegionID:    req.RegionId,
		DiskID:      req.DiskId,
		ClientToken: req.ClientToken,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	response := &ebsDeleteRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &EbsDeleteResponse{
		MasterOrderId: response.MasterOrderID,
		MasterOrderNo: response.MasterOrderNO,
	}, nil
}

type ebsDeleteRealRequest struct {
	RegionID    string `json:"regionID"`
	DiskID      string `json:"diskID"`
	ClientToken string `json:"clientToken"`
}

type ebsDeleteRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
}

type EbsDeleteRequest struct {
	RegionId    string
	DiskId      string
	ClientToken string
}

type EbsDeleteResponse struct {
	MasterOrderId string
	MasterOrderNo string
}
