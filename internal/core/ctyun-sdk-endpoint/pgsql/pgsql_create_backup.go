package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlCreateBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlCreateBackupApi(client *ctyunsdk.CtyunClient) *PgsqlCreateBackupApi {
	return &PgsqlCreateBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/backup/create",
		},
	}
}

func (this *PgsqlCreateBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlCreateBackupRequest, header *PgsqlCreateBackupRequestHeader) (CreateBackupResp *PgsqlCreateBackupResponse, err error) {
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

	if req.BackupName == "" {
		err = fmt.Errorf("backup_name is required")
		return
	}

	builder.AddHeader("regionId", header.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CreateBackupResp = &PgsqlCreateBackupResponse{}
	err = resp.Parse(CreateBackupResp)
	if err != nil {
		return
	}
	return CreateBackupResp, nil
}

type PgsqlCreateBackupRequest struct {
	ProdInstId string  `json:"prodInstId"` // 实例ID，必填
	BackupName string  `json:"backupName"` // 账户名称
	Desc       *string `json:"desc"`       // 账户密码（安全考虑需要用base64加密后传输）
}

type PgsqlCreateBackupRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlCreateBackupResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
	//ReturnObj  *PgsqlCreateBackupResponseReturnObj `json:"returnObj"`
}

//type PgsqlCreateBackupResponseReturnObj struct {
//	Data int64 `json:"data"` // 备份记录id
//}
