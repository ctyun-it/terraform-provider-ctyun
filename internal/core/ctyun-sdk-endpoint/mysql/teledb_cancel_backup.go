package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCancelBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCancelBackupApi(client *ctyunsdk.CtyunClient) *TeledbCancelBackupApi {
	return &TeledbCancelBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v2/open-api/backup/cancel",
		},
	}
}

func (this *TeledbCancelBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCancelBackupRequest, header *TeledbCancelBackupRequestHeader) (CancelBackupResp *TeledbCancelBackupResponse, err error) {
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
	if header.InstID == "" {
		err = fmt.Errorf("inst_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("backupRecordId", fmt.Sprintf("%d", req.BackupRecordId))
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CancelBackupResp = &TeledbCancelBackupResponse{}
	err = resp.Parse(CancelBackupResp)
	if err != nil {
		return
	}
	return CancelBackupResp, nil
}

type TeledbCancelBackupRequest struct {
	BackupRecordId int64 `json:"backupRecordId"`
}
type TeledbCancelBackupRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCancelBackupResponse struct {
	StatusCode int32                               `json:"statusCode"`      // 接口状态码
	Error      *string                             `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                              `json:"message"`         // 描述信息
	ReturnObj  TeledbCancelBackupResponseReturnObj `json:"returnObj"`
}
type TeledbCancelBackupResponseReturnObj struct {
	Data string `json:"data"` // 备份记录id
}
