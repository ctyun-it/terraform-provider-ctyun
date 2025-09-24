package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateInstanceNameApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateInstanceNameApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateInstanceNameApi {
	return &PgsqlUpdateInstanceNameApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/product/modify-instance-name",
		},
	}
}

func (this *PgsqlUpdateInstanceNameApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateInstanceNameRequest, header *PgsqlUpdateInstanceNameRequestHeader) (updatedNameResp *PgsqlUpdateInstanceNameResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	if req.ProdInstId == "" {
		err = errors.New("missing required field: ProdInstId")
		return
	}
	if req.InstanceName == "" {
		err = errors.New("missing required field: InstanceName(实例名称)")
	}

	builder.AddParam("prodInstId", req.ProdInstId)
	builder.AddParam("instanceName", req.InstanceName)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	updatedNameResp = &PgsqlUpdateInstanceNameResponse{}
	err = resp.Parse(updatedNameResp)
	if err != nil {
		return
	}
	return updatedNameResp, nil
}

type PgsqlUpdateInstanceNameRequest struct {
	ProdInstId   string `json:"prodInstId"`
	InstanceName string `json:"instanceName"`
}

type PgsqlUpdateInstanceNameRequestHeader struct {
	ProjectID *string `json:"projectId,omitempty"`
	RegionID  string  `json:"regionId"`
}

type PgsqlUpdateInstanceNameResponse struct {
	StatusCode int32  `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}
