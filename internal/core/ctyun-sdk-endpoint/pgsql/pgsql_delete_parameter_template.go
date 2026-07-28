package pgsql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDeleteParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDeleteParameterTemplateApi(client *ctyunsdk.CtyunClient) *PgsqlDeleteParameterTemplateApi {
	return &PgsqlDeleteParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/PG/v1/template/delete",
		},
	}
}

func (this *PgsqlDeleteParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDeleteParameterTemplateRequest, header *PgsqlDeleteParameterTemplateRequestHeader) (DeleteParameterTemplateResp *PgsqlDeleteParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("id", fmt.Sprintf("%d", req.TemplateId))
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	DeleteParameterTemplateResp = &PgsqlDeleteParameterTemplateResponse{}
	err = resp.Parse(DeleteParameterTemplateResp)
	if err != nil {
		return
	}
	return DeleteParameterTemplateResp, nil
}

type PgsqlDeleteParameterTemplateRequest struct {
	TemplateId int64 `json:"templateId"`
}
type PgsqlDeleteParameterTemplateRequestHeader struct {
	ProjectID string `json:"Project-Id"`
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type PgsqlDeleteParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
