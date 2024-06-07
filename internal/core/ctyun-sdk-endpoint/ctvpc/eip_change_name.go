package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EipChangeNameApi 修改弹性 IP 名字
// https://www.ctyun.cn/document/10026753/10045518
type EipChangeNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipChangeNameApi(client *ctyunsdk.CtyunClient) *EipChangeNameApi {
	return &EipChangeNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/change-name",
		},
	}
}

func (this *EipChangeNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipChangeNameRequest) (*EipChangeNameResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipChangeNameRealRequest{
		ClientToken: req.ClientToken,
		RegionID:    req.RegionId,
		EipID:       req.EipId,
		Name:        req.Name,
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

	return &EipChangeNameResponse{}, nil
}

type eipChangeNameRealRequest struct {
	ClientToken string `json:"clientToken"`
	RegionID    string `json:"regionID"`
	EipID       string `json:"eipID"`
	Name        string `json:"name"`
}

type EipChangeNameRequest struct {
	ClientToken string
	RegionId    string
	EipId       string
	Name        string
}

type EipChangeNameResponse struct {
}
