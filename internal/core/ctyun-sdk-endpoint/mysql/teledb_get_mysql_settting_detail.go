package mysql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetBackupSettingDetailApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetBackupSettingDetailApi(client *ctyunsdk.CtyunClient) *TeledbGetBackupSettingDetailApi {
	return &TeledbGetBackupSettingDetailApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/backupConfig/get",
		},
	}
}

func (this *TeledbGetBackupSettingDetailApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetBackupSettingDetailRequest, header *TeledbGetBackupSettingDetailRequestHeader) (GetBackupSettingDetailResp *TeledbGetBackupSettingDetailResponse, err error) {
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

	if req.OuterProdInstId == "" {
		err = errors.New("instId 为空")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddParam("outerProdInstId", req.OuterProdInstId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetBackupSettingDetailResp = &TeledbGetBackupSettingDetailResponse{}
	err = resp.Parse(GetBackupSettingDetailResp)
	if err != nil {
		return
	}
	return GetBackupSettingDetailResp, nil
}

type TeledbGetBackupSettingDetailRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 外部实例ID，必填
}

type TeledbGetBackupSettingDetailRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}

type TeledbGetBackupSettingDetailResponse struct {
	StatusCode int32                                          `json:"statusCode"`      // 接口状态码
	Error      *string                                        `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                         `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetBackupSettingDetailResponseReturnObj `json:"returnObj"`
}

type TeledbGetBackupSettingDetailResponseReturnObj struct {
	ID                        int64   `json:"id"`                        // 主键ID
	ProdInstID                int64   `json:"prodInstId"`                // 实例ID
	FirstBackupTime           string  `json:"firstbackuptime"`           // 首次开始备份的时间
	BackupUnitTime            int64   `json:"backupunittime"`            // 每两次备份的时间间隔,以秒为单位
	ExpiredTime               int64   `json:"expiredTime"`               // 保留天数; 单位: 秒,默认1天
	FrequencyBackup           bool    `json:"frequencyBackup"`           // 高频备份 开启 关闭
	FrequencyBackupUnitTime   int64   `json:"frequencyBackupUnittime"`   // 高频备份频率 单位: 秒,默认1小时
	AllowEarliestTime         string  `json:"allowEarliestTime"`         // 允许最早开始备份时间 默认：00:00:00
	Status                    int32   `json:"status"`                    // 当前状态: 0表示没有任何操作, 1表示当前正在备份, 2表示当前正在恢复
	BackupOperation           bool    `json:"backupOperation"`           // 是否开启自动备份操作
	ProdInstName              string  `json:"prodInstName"`              // 实例名称
	OuterProdInstID           string  `json:"outerProdInstId"`           // 开通的实例id
	BackupRecoveryHeartbeatOk bool    `json:"backuprecoveryHeartbeatOk"` // 工具心跳
	SyncerHeartbeatOk         bool    `json:"syncerHeartbeatOk"`         // 实时同步心跳
	BucketName                string  `json:"bucketName"`                // 使用对象存储时桶名称
	CrossRegion               bool    `json:"crossRegion"`               // 是否开启跨域备份
	ResourceName              string  `json:"resourceName"`              // 资源名称
	StorePath                 string  `json:"storePath"`                 // 使用监控机时备份存储路径
	BackupInfoID              int64   `json:"backupInfoId"`              // 备份使用资源id
	Region                    string  `json:"region"`                    // 备份所在资源池id
	UseZOS                    int32   `json:"usezos"`                    // 是否使用对象存储，0代表否，1代表是
	TargetCrossRegion         string  `json:"targetCrossRegion"`         // 跨域备份的目标资源池id
	MinExpiredTime            int64   `json:"minExpiredTime"`            // 最小可以设置的备份过期时间，单位为秒
	MaxExpiredTime            int64   `json:"maxExpiredTime"`            // 最大可以设置的备份过期时间，单位为秒
	ConfigSyncerKeepAlive     bool    `json:"configSyncerKeepAlive"`     // 对应syncer是否开启了syncer保活
	TriggerDaysOfWeek         []int32 `json:"triggerDaysOfWeek"`         // 全备触发星期，单个元素取值范围为1~7，1代表周天，2代表周一以此类推
}
