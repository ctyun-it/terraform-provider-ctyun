package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// EipDeleteApi 删除弹性 IP
// https://www.ctyun.cn/document/10026753/10040760
type EipDeleteApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipDeleteApi(client *ctyunsdk.CtyunClient) *EipDeleteApi {
	return &EipDeleteApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/delete",
		},
	}
}

func (this *EipDeleteApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipDeleteRequest) (*EipDeleteResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipDeleteRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		EipID:       req.EipId,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &eipDeleteRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}

	return &EipDeleteResponse{
		MasterOrderId: result.MasterOrderID,
		MasterOrderNo: result.MasterOrderNO,
		RegionId:      result.RegionID,
	}, nil
}

type eipDeleteRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	EipID       string `json:"eipID"`
}

type eipDeleteRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
}

type EipDeleteRequest struct {
	ClientToken string
	RegionId    string
	EipId       string
}

type EipDeleteResponse struct {
	MasterOrderId string
	MasterOrderNo string
	RegionId      string
}
