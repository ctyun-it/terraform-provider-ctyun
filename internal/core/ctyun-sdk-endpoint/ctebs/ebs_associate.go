package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EbsAssociateApi 云硬盘绑定
// https://www.ctyun.cn/document/10027696/10110703
type EbsAssociateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEbsAssociateApi(client *ctyunsdk.CtyunClient) *EbsAssociateApi {
	return &EbsAssociateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/ebs/attach-ebs",
		},
	}
}

func (this *EbsAssociateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EbsAssociateRequest) (*EbsAssociateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := ebsAssociateRealRequest{
		RegionID:   req.RegionId,
		DiskID:     req.DiskId,
		InstanceID: req.InstanceId,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtebs, builder)
	if err != nil {
		return nil, err
	}

	resp := &ebsAssociateRealResponse{}
	err = response.ParseByStandardModelWithCheck(resp)
	if err != nil {
		return nil, err
	}

	return &EbsAssociateResponse{
		DiskJobId: resp.DiskJobID,
	}, nil
}

type ebsAssociateRealRequest struct {
	RegionID   string `json:"regionID"`
	DiskID     string `json:"diskID"`
	InstanceID string `json:"instanceID"`
}

type ebsAssociateRealResponse struct {
	DiskJobID string `json:"diskJobID"`
}

type EbsAssociateRequest struct {
	RegionId   string
	DiskId     string
	InstanceId string
}

type EbsAssociateResponse struct {
	DiskJobId string
}
