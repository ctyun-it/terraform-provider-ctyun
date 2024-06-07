package ctebs

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EbsChangeNameApi 云硬盘修改名称
// https://www.ctyun.cn/document/10027696/10110706
type EbsChangeNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsChangeNameApi(client *ctyunsdk.CtyunClient) *EbsChangeNameApi {
	return &EbsChangeNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/update-attr-ebs",
		},
	}
}

func (this *EbsChangeNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsChangeNameRequest) (*EbsChangeNameResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	_, err := builder.WriteJson(&ebsChangeNameRealRequest{
		RegionID: req.RegionId,
		DiskID:   req.DiskId,
		DiskName: req.DiskName,
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
	return nil, nil
}

type ebsChangeNameRealRequest struct {
	RegionID string `json:"regionID"`
	DiskID   string `json:"diskID"`
	DiskName string `json:"diskName"`
}

type EbsChangeNameRequest struct {
	RegionId string
	DiskId   string
	DiskName string
}

type EbsChangeNameResponse struct {
}
