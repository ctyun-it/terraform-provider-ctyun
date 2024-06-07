package ctvpc

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// EipAssociateApi 绑定弹性 IP
// https://www.ctyun.cn/document/10026753/10031946
type EipAssociateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewEipAssociateApi(client *ctyunsdk.CtyunClient) *EipAssociateApi {
	return &EipAssociateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/eip/associate",
		},
	}
}

func (this *EipAssociateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *EipAssociateRequest) (*EipAssociateResponse, ctyunsdk.CtyunRequestError) {
	builder := this.WithCredential(&credential)
	requestContent := eipAssociateRealRequest{
		RegionID:        req.RegionId,
		ClientToken:     req.ClientToken,
		EipID:           req.EipId,
		AssociationID:   req.AssociationId,
		AssociationType: req.AssociationType,
		ProjectID:       req.ProjectId,
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

	return &EipAssociateResponse{}, nil
}

type eipAssociateRealRequest struct {
	RegionID        string `json:"regionID"`
	ClientToken     string `json:"clientToken"`
	EipID           string `json:"eipID"`
	AssociationID   string `json:"associationID"`
	AssociationType int    `json:"associationType"`
	ProjectID       string `json:"projectID,omitempty"`
}

type EipAssociateRequest struct {
	RegionId        string
	ClientToken     string
	EipId           string
	AssociationId   string
	AssociationType int
	ProjectId       string
}

type EipAssociateResponse struct {
}
