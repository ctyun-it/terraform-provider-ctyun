package mysql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateRecoveryJobApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateRecoveryJobApi(client *ctyunsdk.CtyunClient) *TeledbCreateRecoveryJobApi {
	return &TeledbCreateRecoveryJobApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v2/open-api/recovery/createNewRecoveryJob",
		},
	}
}

func (this *TeledbCreateRecoveryJobApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateRecoveryJobRequest, header *TeledbCreateRecoveryJobRequestHeader) (CreateRecoveryJobResp *TeledbCreateRecoveryJobResponse, err error) {
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
	if header.InstID == "" {
		err = fmt.Errorf("inst_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CreateRecoveryJobResp = &TeledbCreateRecoveryJobResponse{}
	err = resp.Parse(CreateRecoveryJobResp)
	if err != nil {
		return
	}
	return CreateRecoveryJobResp, nil
}

type TeledbCreateRecoveryJobRequest struct {
	SrcOuterProdInstId string   `json:"srcOuterProdInstId"`    // 实例ID，必填
	DstOuterProdInstId string   `json:"dstOuterProdInstId"`    // 实例ID，必填
	ToTimepoint        *string  `json:"toTimepoint,omitempty"` // 账户名称
	TaskId             *string  `json:"taskId,omitempty"`      // 用来恢复的备份任务集【taskId和toTimepoint不能同时为空，优先取toTimepoint】
	RestoreDbName      string   `json:"restoreDbName"`         // 恢复的库名
	RestoreTables      []string `json:"restoreTables"`         // 恢复的表名
	NewDbName          *string  `json:"newDbName"`             // 新的库名
	NewTables          []string `json:"newTables"`             // 新的表名
}
type TeledbCreateRecoveryJobRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbCreateRecoveryJobResponseReturnObj struct {
	Data string `json:"data"`
}

type TeledbCreateRecoveryJobResponse struct {
	StatusCode int32                                     `json:"statusCode"`      // 接口状态码
	Error      *string                                   `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                    `json:"message"`         // 描述信息
	ReturnObj  *TeledbCreateRecoveryJobResponseReturnObj `json:"returnObj"`
}
