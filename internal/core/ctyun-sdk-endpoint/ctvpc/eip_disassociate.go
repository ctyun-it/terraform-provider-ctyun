package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EipDisassociateApi 解绑弹性 IP
// https://www.ctyun.cn/document/10026753/10031948
type EipDisassociateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipDisassociateApi(client *ctyunsdk.CtyunClient) *EipDisassociateApi {
	return &EipDisassociateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/disassociate",
		},
	}
}

func (this *EipDisassociateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipDisassociateRequest) (*EipDisassociateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipDisassociateRealRequest{
		RegionID:    req.RegionId,
		ClientToken: req.ClientToken,
		EipID:       req.EipId,
		ProjectID:   req.ProjectId,
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

	return &EipDisassociateResponse{}, nil
}

type eipDisassociateRealRequest struct {
	RegionID    string `json:"regionID"`
	ClientToken string `json:"clientToken"`
	EipID       string `json:"eipID"`
	ProjectID   string `json:"projectID,omitempty"`
}

type EipDisassociateRequest struct {
	RegionId    string
	ClientToken string
	EipId       string
	ProjectId   string
}

type EipDisassociateResponse struct {
}
