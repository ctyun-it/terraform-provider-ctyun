package pgsql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlCreateParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlCreateParameterTemplateApi(client *ctyunsdk.CtyunClient) *PgsqlCreateParameterTemplateApi {
	return &PgsqlCreateParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/template/create",
		},
	}
}

func (this *PgsqlCreateParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlCreateParameterTemplateRequest, header *PgsqlCreateParameterTemplateRequestHeader) (CreateParameterTemplateResp *PgsqlCreateParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CreateParameterTemplateResp = &PgsqlCreateParameterTemplateResponse{}
	err = resp.Parse(CreateParameterTemplateResp)
	if err != nil {
		return
	}
	return CreateParameterTemplateResp, nil
}

type PgsqlCreateParameterTemplateRequest struct {
	SourceTemplateId int64   `json:"sourceTemplateId"`
	Name             string  `json:"name"`                  // 实例ID，必填
	Description      *string `json:"description,omitempty"` //
}
type PgsqlCreateParameterTemplateRequestHeader struct {
	ProjectID *string `json:"Project-Id,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlCreateParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
