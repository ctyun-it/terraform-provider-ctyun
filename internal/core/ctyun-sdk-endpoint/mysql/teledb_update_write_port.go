package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateWritePortApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateWritePortApi(client *ctyunsdk.CtyunClient) *TeledbUpdateWritePortApi {
	return &TeledbUpdateWritePortApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/instance/modify-instance-write-port",
		},
	}
}

func (this *TeledbUpdateWritePortApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateWritePortRequest, headers *TeledbUpdateWritePortRequestHeader) (updatedResponse *TeledbUpdateWritePortResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if headers.ProjectID != "" {
		builder.AddHeader("project-id", headers.ProjectID)
	}
	if headers.InstID != "" {
		builder.AddHeader("inst-id", headers.InstID)
	}
	if headers.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	if req.OuterProdInstId == "" {
		err = errors.New("missing required field: InstID")
		return
	}
	if req.WritePort == "" {
		err = errors.New("missing required field: WritePort")
		return
	}
	builder.AddHeader("regionId", headers.RegionID)
	builder.AddParam("outerProdInstId", req.OuterProdInstId)
	builder.AddParam("writePort", req.WritePort)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	updatedResponse = &TeledbUpdateWritePortResponse{}
	err = resp.Parse(updatedResponse)
	if err != nil {
		return
	}
	return updatedResponse, nil
}

type TeledbUpdateWritePortRequest struct {
	OuterProdInstId string `json:"outerProdInstId"`
	WritePort       string `json:"writePort"`
}

type TeledbUpdateWritePortRequestHeader struct {
	ProjectID string `json:"projectId"`
	InstID    string `json:"instId"`
	RegionID  string `json:"regionId"`
}

type TeledbUpdateWritePortResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
