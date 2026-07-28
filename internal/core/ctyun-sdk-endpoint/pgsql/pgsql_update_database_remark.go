package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateDatabaseRemarkApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateDatabaseRemarkApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateDatabaseRemarkApi {
	return &PgsqlUpdateDatabaseRemarkApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/database/modify-description",
		},
	}
}

func (this *PgsqlUpdateDatabaseRemarkApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateDatabaseRemarkRequest, header *PgsqlUpdateDatabaseRemarkRequestHeader) (UpdateDatabaseRemarkResp *PgsqlUpdateDatabaseRemarkResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	if req.ProdInstId == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddHeader("regionId", header.RegionID)
	builder.AddParam("prodInstId", req.ProdInstId)
	builder.AddParam("dbName", req.DBName)
	builder.AddParam("description", *req.Description)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	UpdateDatabaseRemarkResp = &PgsqlUpdateDatabaseRemarkResponse{}
	err = resp.Parse(UpdateDatabaseRemarkResp)
	if err != nil {
		return
	}
	return UpdateDatabaseRemarkResp, nil
}

type PgsqlUpdateDatabaseRemarkRequest struct {
	ProdInstId  string  `json:"prodInstId"` // 实例ID，必填
	DBName      string  `json:"dbName"`     // 账户名称
	Description *string `json:"description,omitempty"`
}
type PgsqlUpdateDatabaseRemarkRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlUpdateDatabaseRemarkResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
