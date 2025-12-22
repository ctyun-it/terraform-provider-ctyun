package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateParameterTemplateRemarkApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateParameterTemplateRemarkApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateParameterTemplateRemarkApi {
	return &PgsqlUpdateParameterTemplateRemarkApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/template/modifyDesc",
		},
	}
}

func (this *PgsqlUpdateParameterTemplateRemarkApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateParameterTemplateRemarkRequest, header *PgsqlUpdateParameterTemplateRemarkRequestHeader) (UpdateDatabaseRemarkResp *PgsqlUpdateParameterTemplateRemarkResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	UpdateDatabaseRemarkResp = &PgsqlUpdateParameterTemplateRemarkResponse{}
	err = resp.Parse(UpdateDatabaseRemarkResp)
	if err != nil {
		return
	}
	return UpdateDatabaseRemarkResp, nil
}

type PgsqlUpdateParameterTemplateRemarkRequest struct {
	TemplateId  int64  `json:"templateId"`  // 实例ID，必填
	Description string `json:"description"` // 账户名称
}
type PgsqlUpdateParameterTemplateRemarkRequestHeader struct {
	ProjectID string `json:"Project-Id"`
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type PgsqlUpdateParameterTemplateRemarkResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
