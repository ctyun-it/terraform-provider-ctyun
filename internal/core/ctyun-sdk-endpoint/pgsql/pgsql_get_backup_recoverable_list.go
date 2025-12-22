package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetRecoverableBackupListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetRecoverableBackupListApi(client *ctyunsdk.CtyunClient) *PgsqlGetRecoverableBackupListApi {
	return &PgsqlGetRecoverableBackupListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/backup/list",
		},
	}
}

func (this *PgsqlGetRecoverableBackupListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetRecoverableBackupListRequest, header *PgsqlGetRecoverableBackupListRequestHeader) (GetRecoverableTimeRangesResp *PgsqlGetRecoverableBackupListResponse, err error) {
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

	builder.AddParam("outerProdInstId", req.ProdInstId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetRecoverableTimeRangesResp = &PgsqlGetRecoverableBackupListResponse{}
	err = resp.Parse(GetRecoverableTimeRangesResp)
	if err != nil {
		return
	}
	return GetRecoverableTimeRangesResp, nil
}

type PgsqlGetRecoverableBackupListRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
}
type PgsqlGetRecoverableBackupListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGetRecoverableBackupListResponse struct {
	StatusCode int32                                            `json:"statusCode"`      // 接口状态码
	Error      *string                                          `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                           `json:"message"`         // 描述信息
	ReturnObj  []PgsqlGetRecoverableBackupListResponseReturnObj `json:"returnObj"`
}

type PgsqlGetRecoverableBackupListResponseReturnObj struct {
	PgBackupBaseInfoId int64  `json:"pgBackupBaseInfoId"`
	BackupName         string `json:"backupName"`
}
