package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EbsChangeSizeApi 云硬盘修改规格
// https://www.ctyun.cn/document/10027696/10110705
type EbsChangeSizeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsChangeSizeApi(client *ctyunsdk.CtyunClient) *EbsChangeSizeApi {
	return &EbsChangeSizeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/resize-ebs",
		},
	}
}

func (this *EbsChangeSizeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsChangeSizeRequest) (*EbsChangeSizeResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ebsChangeSizeRealRequest{
		RegionID:    req.RegionId,
		DiskID:      req.DiskId,
		DiskSize:    req.DiskSize,
		ClientToken: req.ClientToken,
	})
	if err != nil {
		return nil, err
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	response := &ebsChangeSizeRealResponse{}
	err = resp.ParseByStandardModelWithCheck(response)
	if err != nil {
		return nil, err
	}
	return &EbsChangeSizeResponse{
		MasterOrderId: response.MasterOrderID,
		MasterOrderNo: response.MasterOrderNO,
	}, nil
}

type ebsChangeSizeRealRequest struct {
	RegionID    string `json:"regionID"`
	DiskID      string `json:"diskID"`
	DiskSize    int    `json:"diskSize"`
	ClientToken string `json:"clientToken"`
}

type ebsChangeSizeRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
}

type EbsChangeSizeRequest struct {
	RegionId    string
	DiskId      string
	DiskSize    int
	ClientToken string
}

type EbsChangeSizeResponse struct {
	MasterOrderId string
	MasterOrderNo string
}
