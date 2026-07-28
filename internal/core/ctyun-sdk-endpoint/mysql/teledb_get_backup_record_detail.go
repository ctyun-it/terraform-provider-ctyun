package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetBackupRecordDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetBackupRecordDetailApi(client *ctyunsdk.CtyunClient) *TeledbGetBackupRecordDetailApi {
	return &TeledbGetBackupRecordDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/backup/get",
		},
	}
}

func (this *TeledbGetBackupRecordDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetBackupRecordDetailRequest, header *TeledbGetBackupRecordDetailRequestHeader) (GetBackupRecordDetailResp *TeledbGetBackupRecordDetailResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)
	if req.Id == "" {
		err = fmt.Errorf("id is required")
		return
	}
	builder.AddParam("id", req.Id)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetBackupRecordDetailResp = &TeledbGetBackupRecordDetailResponse{}
	err = resp.Parse(GetBackupRecordDetailResp)
	if err != nil {
		return
	}
	return GetBackupRecordDetailResp, nil
}

type TeledbGetBackupRecordDetailRequest struct {
	Id string `json:"id"` // backupRecordId
}
type TeledbGetBackupRecordDetailRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbGetBackupRecordDetailResponse struct {
	StatusCode int32                                         `json:"statusCode"`      // 接口状态码
	Error      *string                                       `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                        `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetBackupRecordDetailResponseReturnObj `json:"returnObj"`
}

type TeledbGetBackupRecordDetailResponseReturnObj struct {
	BackupRecordId   int64  `json:"backupRecordId"` // 备份记录id
	ProdInstName     string `json:"prodInstName"`
	OpType           string `json:"opType"`
	Description      string `json:"description"`     // 备份描述
	BackupStartTime  string `json:"backupStartTime"` // 备份开始时间
	OuterProdInstId  string `json:"outerProdInstId"` // 外部实例id
	TaskType         string `json:"taskType"`        // 备份类型（full/incr）
	BackedUpDataSize int32  `json:"backedUpDataSize"`
	BackupTaskId     int64  `json:"backupTaskId"`
	TaskId           string `json:"TaskId"`
	StorageType      string `json:"storageType"` // 存储类型（s3/disk/region_s3）
	ProdInstId       int64  `json:"prodInstId"`  // 实例id
	Disabled         bool   `json:"disabled"`    // 禁用备份
	BackupName       string `json:"backupName"`  // 备份名称
	TaskStatus       int32  `json:"taskStatus"`  // 任务状态（100/101/102/1/-1）
}
