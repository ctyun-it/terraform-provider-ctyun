package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDeleteBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDeleteBackupApi(client *ctyunsdk.CtyunClient) *PgsqlDeleteBackupApi {
	return &PgsqlDeleteBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/PG/v1/backup/delete",
		},
	}
}

func (this *PgsqlDeleteBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDeleteBackupRequest, header *PgsqlDeleteBackupRequestHeader) (DeleteBackupResp *PgsqlDeleteBackupResponse, err error) {
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
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}

	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("prodInstId", req.ProdInstId)
	builder.AddParam("backupId", fmt.Sprintf("%d", req.BackupId))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	DeleteBackupResp = &PgsqlDeleteBackupResponse{}
	err = resp.Parse(DeleteBackupResp)
	if err != nil {
		return
	}
	return DeleteBackupResp, nil
}

type PgsqlDeleteBackupRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	BackupId   int64  `json:"backupId"`
}
type PgsqlDeleteBackupRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlDeleteBackupResponse struct {
	StatusCode int32                              `json:"statusCode"`      // 接口状态码
	Error      *string                            `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                             `json:"message"`         // 描述信息
	ReturnObj  PgsqlDeleteBackupResponseReturnObj `json:"returnObj"`
}
type PgsqlDeleteBackupResponseReturnObj struct {
	Data string `json:"data"` // 备份记录id
}
