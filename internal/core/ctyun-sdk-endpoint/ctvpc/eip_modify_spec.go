package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EipModifySpecApi 修改弹性 IP 带宽
// https://www.ctyun.cn/document/10026753/10040762
type EipModifySpecApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipModifySpecApi(client *ctyunsdk.CtyunClient) *EipModifySpecApi {
	return &EipModifySpecApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/modify-spec",
		},
	}
}

func (this *EipModifySpecApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipModifySpecRequest) (*EipModifySpecResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipModifySpecRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		EipID:       req.EipId,
		Bandwidth:   req.Bandwidth,
	}
	_, err := builder.WriteJson(requestContent)
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtvpc, builder)
	if err != nil {
		return nil, err
	}

	result := &eipModifySpecRealResponse{}
	err = response.ParseByStandardModelWithCheck(result)
	if err != nil {
		return nil, err
	}

	return &EipModifySpecResponse{
		MasterOrderId: result.MasterOrderID,
		MasterOrderNo: result.MasterOrderNO,
		RegionId:      result.RegionID,
	}, nil
}

type eipModifySpecRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	EipID       string `json:"eipID"`
	Bandwidth   int    `json:"bandwidth"`
}

type eipModifySpecRealResponse struct {
	MasterOrderID string `json:"masterOrderID"`
	MasterOrderNO string `json:"masterOrderNO"`
	RegionID      string `json:"regionID"`
}

type EipModifySpecRequest struct {
	ClientToken string
	RegionId    string
	EipId       string
	Bandwidth   int
}

type EipModifySpecResponse struct {
	MasterOrderId string
	MasterOrderNo string
	RegionId      string
}
