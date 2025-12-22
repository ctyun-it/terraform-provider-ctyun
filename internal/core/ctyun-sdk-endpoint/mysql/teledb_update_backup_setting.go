package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateBackupSettingApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateBackupSettingApi(client *ctyunsdk.CtyunClient) *TeledbUpdateBackupSettingApi {
	return &TeledbUpdateBackupSettingApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v2/open-api/backupConfig/update",
		},
	}
}

func (this *TeledbUpdateBackupSettingApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateBackupSettingRequest, header *TeledbUpdateBackupSettingRequestHeader) (UpdateBackupSettingResp *TeledbUpdateBackupSettingResponse, err error) {
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
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	UpdateBackupSettingResp = &TeledbUpdateBackupSettingResponse{}
	err = resp.Parse(UpdateBackupSettingResp)
	if err != nil {
		return
	}
	return UpdateBackupSettingResp, nil
}

type TeledbUpdateBackupSettingRequest struct {
	ExpiredTime             int64   `json:"expiredTime"`
	FrequencyBackup         bool    `json:"frequencyBackup"`
	FrequencyBackupUnitTime *int64  `json:"frequencyBackupUnittime"`
	AllowEarliestTime       string  `json:"allowEarliestTime"`
	OuterProdInstId         string  `json:"outerProdInstId"`
	TriggerDaysOfWeek       []int32 `json:"triggerDaysOfWeek"`
}
type TeledbUpdateBackupSettingRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbUpdateBackupSettingResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
