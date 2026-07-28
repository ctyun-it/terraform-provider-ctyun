package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbDeleteBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbDeleteBackupApi(client *ctyunsdk.CtyunClient) *TeledbDeleteBackupApi {
	return &TeledbDeleteBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/RDS2/v2/open-api/backup/manualDelete",
		},
	}
}

func (this *TeledbDeleteBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbDeleteBackupRequest, header *TeledbDeleteBackupRequestHeader) (DeleteBackupResp *TeledbDeleteBackupResponse, err error) {
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
	if req.BlockID == 0 {
		err = fmt.Errorf("block_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("outerProdInstId", req.OuterProdInstId)
	builder.AddParam("blockId", fmt.Sprintf("%d", req.BlockID))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	DeleteBackupResp = &TeledbDeleteBackupResponse{}
	err = resp.Parse(DeleteBackupResp)
	if err != nil {
		return
	}
	return DeleteBackupResp, nil
}

type TeledbDeleteBackupRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	BlockID         int64  `json:"blockId"`
}
type TeledbDeleteBackupRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbDeleteBackupResponse struct {
	StatusCode int32                               `json:"statusCode"`      // 接口状态码
	Error      *string                             `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                              `json:"message"`         // 描述信息
	ReturnObj  TeledbDeleteBackupResponseReturnObj `json:"returnObj"`
}
type TeledbDeleteBackupResponseReturnObj struct {
	Data string `json:"data"` // 备份记录id
}
