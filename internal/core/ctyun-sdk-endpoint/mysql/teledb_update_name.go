package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateInstanceNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateInstanceNameApi(client *ctyunsdk.CtyunClient) *TeledbUpdateInstanceNameApi {
	return &TeledbUpdateInstanceNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/instance/instance-desc",
		},
	}
}

func (this *TeledbUpdateInstanceNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateInstanceNameRequest, header *TeledbUpdateInstanceNameRequestHeader) (updatedNameResp *TeledbUpdateInstanceNameResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.InstID != "" {
		builder.AddHeader("inst-id", header.InstID)
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	if req.OuterProdInstID == "" {
		err = errors.New("missing required field: outerProdInstId")
		return
	}
	if req.InstanceDescription == "" {
		err = errors.New("missing required field: name(实例名称)")
	}

	builder.AddParam("outerProdInstId", req.OuterProdInstID)
	builder.AddParam("instanceDescription", req.InstanceDescription)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	updatedNameResp = &TeledbUpdateInstanceNameResponse{}
	err = resp.Parse(updatedNameResp)
	if err != nil {
		return
	}
	return updatedNameResp, nil
}

type TeledbUpdateInstanceNameRequest struct {
	OuterProdInstID     string `json:"outerProdInstId"`
	InstanceDescription string `json:"instanceDescription"`
}

type TeledbUpdateInstanceNameRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
	InstID    string  `json:"instId"`
	RegionID  string  `json:"regionId"`
}

type TeledbUpdateInstanceNameResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
