package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetBackupTaskListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetBackupTaskListApi(client *ctyunsdk.CtyunClient) *PgsqlGetBackupTaskListApi {
	return &PgsqlGetBackupTaskListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/backup/task-list",
		},
	}
}

func (this *PgsqlGetBackupTaskListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetBackupTaskListRequest, header *PgsqlGetBackupTaskListRequestHeader) (GetRecoverableTimeRangesResp *PgsqlGetBackupTaskListResponse, err error) {
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
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNum))
	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	if req.BackupName != nil {
		builder.AddParam("backupName", *req.BackupName)
	}
	if req.SelectType != nil {
		builder.AddParam("selectType", *req.SelectType)
	}
	if req.StartTime != nil {
		builder.AddParam("startTime", *req.StartTime)
	}
	if req.EndTime != nil {
		builder.AddParam("endTime", *req.EndTime)
	}
	if req.CreateTime != nil {
		builder.AddParam("createTime", *req.CreateTime)
	}
	if req.UpdateTime != nil {
		builder.AddParam("updateTime", *req.UpdateTime)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetRecoverableTimeRangesResp = &PgsqlGetBackupTaskListResponse{}
	err = resp.Parse(GetRecoverableTimeRangesResp)
	if err != nil {
		return
	}
	return GetRecoverableTimeRangesResp, nil
}

type PgsqlGetBackupTaskListRequest struct {
	ProdInstId string  `json:"prodInstId"` // 实例ID，必填
	PageNum    int32   `json:"pageNum"`
	PageSize   int32   `json:"pageSize"`
	BackupName *string `json:"backupName,omitempty"`
	SelectType *string `json:"selectType,omitempty"`
	StartTime  *string `json:"startTime,omitempty"`
	EndTime    *string `json:"endTime,omitempty"`
	CreateTime *string `json:"createTime,omitempty"`
	UpdateTime *string `json:"updateTime,omitempty"`
}
type PgsqlGetBackupTaskListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGetBackupTaskListResponse struct {
	StatusCode int32                                    `json:"statusCode"`      // 接口状态码
	Error      *string                                  `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                   `json:"message"`         // 描述信息
	ReturnObj  *PgsqlGetBackupTaskListResponseReturnObj `json:"returnObj"`
}

type PgsqlGetBackupTaskListResponseReturnObjList struct {
	Id         int64  `json:"id"`
	ProdInstId string `json:"prodInstId"`
	BackupName string `json:"backupName"`
	Type       int32  `json:"type"`
	Result     string `json:"result"`
	Status     int32  `json:"status"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	DataLen    string `json:"dataLen"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type PgsqlGetBackupTaskListResponseReturnObj struct {
	PageNum   int32                                         `json:"pageNum"`
	PageSize  int32                                         `json:"pageSize"`
	PageTotal int32                                         `json:"pageTotal"`
	Total     int32                                         `json:"total"`
	List      []PgsqlGetBackupTaskListResponseReturnObjList `json:"list"`
}
