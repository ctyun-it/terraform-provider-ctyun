package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateBackupApi(client *ctyunsdk.CtyunClient) *TeledbCreateBackupApi {
	return &TeledbCreateBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v2/open-api/backup/createManual",
		},
	}
}

func (this *TeledbCreateBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateBackupRequest, header *TeledbCreateBackupRequestHeader) (CreateBackupResp *TeledbCreateBackupResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OuterProdInstId == "" || header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	if header.InstID == "" {
		err = fmt.Errorf("inst_id is required")
		return
	}
	if req.BackupName == "" {
		err = fmt.Errorf("backup_name is required")
		return
	}
	if req.TaskType == "" {
		err = fmt.Errorf("task_type is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CreateBackupResp = &TeledbCreateBackupResponse{}
	err = resp.Parse(CreateBackupResp)
	if err != nil {
		return
	}
	return CreateBackupResp, nil
}

type TeledbCreateBackupRequest struct {
	OuterProdInstId string  `json:"outerProdInstId"` // 实例ID，必填
	BackupName      string  `json:"backupName"`      // 账户名称
	Description     *string `json:"description"`     // 账户密码（安全考虑需要用base64加密后传输）
	TaskType        string  `json:"taskType"`        // 备份类型,默认全量物理备份 全量物理备份:full 全量逻辑备份:logic_full
}
type TeledbCreateBackupRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCreateBackupResponse struct {
	StatusCode int32                                `json:"statusCode"`      // 接口状态码
	Error      *string                              `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                               `json:"message"`         // 描述信息
	ReturnObj  *TeledbCreateBackupResponseReturnObj `json:"returnObj"`
}
type TeledbCreateBackupResponseReturnObj struct {
	Data int64 `json:"data"` // 备份记录id
}
