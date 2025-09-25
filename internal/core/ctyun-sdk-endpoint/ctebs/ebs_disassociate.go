package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EbsDisassociateApi 云硬盘解绑
// https://www.ctyun.cn/document/10027696/10110704
type EbsDisassociateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsDisassociateApi(client *ctyunsdk.CtyunClient) *EbsDisassociateApi {
	return &EbsDisassociateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/detach-ebs",
		},
	}
}

func (this *EbsDisassociateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsDisassociateRequest) (*EbsDisassociateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := ebsDisassociateRealRequest{
		DiskID:   req.DiskId,
		RegionID: req.RegionId,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	resp := &ebsDisassociateRealResponse{}
	err = response.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}

	return &EbsDisassociateResponse{
		DiskJobId: resp.DiskJobID,
	}, nil
}

type ebsDisassociateRealRequest struct {
	DiskID   string `json:"diskID"`
	RegionID string `json:"regionID"`
}

type ebsDisassociateRealResponse struct {
	DiskJobID string `json:"diskJobID"`
}

type EbsDisassociateRequest struct {
	DiskId   string
	RegionId string
}

type EbsDisassociateResponse struct {
	DiskJobId string
}
