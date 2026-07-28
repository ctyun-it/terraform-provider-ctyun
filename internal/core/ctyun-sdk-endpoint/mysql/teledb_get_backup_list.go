package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetBackupListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetBackupListApi(client *ctyunsdk.CtyunClient) *TeledbGetBackupListApi {
	return &TeledbGetBackupListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/backup/list",
		},
	}
}

func (this *TeledbGetBackupListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetBackupListRequest, header *TeledbGetBackupListRequestHeader) (GetBackupListResp *TeledbGetBackupListResponse, err error) {
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
	builder.AddParam("outerProdInstId", req.OuterProdInstId)

	if req.ProdInstName != nil {
		builder.AddParam("prodInstName", *req.ProdInstName)
	}
	if req.BackupName != nil {
		builder.AddParam("backupName", *req.BackupName)
	}
	if req.BlockId != nil {
		builder.AddParam("blockId", fmt.Sprintf("%d", *req.BlockId))
	}
	if req.StartTime != nil {
		builder.AddParam("startTime", *req.StartTime)
	}
	if req.EndTime != nil {
		builder.AddParam("endTime", *req.EndTime)
	}
	if req.PageNow == 0 {
		err = errors.New("page_no 为空")
		return
	}
	if req.PageSize == 0 {
		err = errors.New("page_size 为空")
		return
	}

	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))
	builder.AddParam("pageNow", fmt.Sprintf("%d", req.PageNow))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetBackupListResp = &TeledbGetBackupListResponse{}
	err = resp.Parse(GetBackupListResp)
	if err != nil {
		return
	}
	return GetBackupListResp, nil
}

type TeledbGetBackupListRequest struct {
	OuterProdInstId string  `json:"outerProdInstId"`        // 外部实例ID，必填
	ProdInstName    *string `json:"prodInstName,omitempty"` // 实例名
	BackupName      *string `json:"backupName,omitempty"`   // 备份名
	BlockId         *int64  `json:"blockId,omitempty"`      // 备份周期ID
	StartTime       *string `json:"startTime,omitempty"`    // 开始时间
	EndTime         *string `json:"endTime,omitempty"`      // 结束时间
	PageNow         int32   `json:"pageNow"`                // 当前页码，必填
	PageSize        int32   `json:"pageSize"`               // 分页大小，必填
}
type TeledbGetBackupListRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbGetBackupListResponse struct {
	StatusCode int32                                 `json:"statusCode"`      // 接口状态码
	Error      *string                               `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetBackupListResponseReturnObj `json:"returnObj"`
}
type TeledbGetBackupListResponseReturnObj struct {
	PageNum           int32          `json:"pageNum"`           // 当前页
	PageSize          int32          `json:"pageSize"`          // 每页的数量
	Size              int32          `json:"size"`              // 当前页的数量
	StartRow          int32          `json:"startRow"`          // 当前页面第一个元素在数据库中的行号
	EndRow            int32          `json:"endRow"`            // 当前页面最后一个元素在数据库中的行号
	Total             int32          `json:"total"`             // 总记录数
	Pages             int32          `json:"pages"`             // 总页数
	PrePage           int32          `json:"prePage"`           // 前一页
	IsFirstPage       bool           `json:"isFirstPage"`       // 是否为第一页
	IsLastPage        bool           `json:"isLastPage"`        // 是否为最后一页
	HasPreviousPage   bool           `json:"hasPreviousPage"`   // 是否有前一页
	HasNextPage       bool           `json:"hasNextPage"`       // 是否有下一页
	NavigatePages     int32          `json:"navigatePages"`     // 导航页码数
	NavigatePageNums  []int32        `json:"navigatepageNums"`  // 所有导航页号
	List              []BackupPeriod `json:"list"`              // 结果集
	NavigateLastPage  int32          `json:"navigateLastPage"`  // 页面上显示的最后一个页码
	NavigateFirstPage int32          `json:"navigateFirstPage"` // 页面显示的第一个页码
}

// 备份周期对象
type BackupPeriod struct {
	ProdInstId      int64          `json:"prodInstId"`      // 实例id
	OuterProdInstId string         `json:"outerProdInstId"` // 外部实例id
	ProdInstName    string         `json:"prodInstName"`    // 实例名称
	FreezeStatus    string         `json:"freezeStatus"`
	KmsStatus       string         `json:"kmsStatus"`
	BlockId         int64          `json:"blockId"` // 备份周期id
	Records         []BackupRecord `json:"records"` // 备份记录集合
}

// 备份记录对象
type BackupRecord struct {
	BackupRecordId        int64  `json:"backupRecordId"`        // 备份记录id
	BackupTaskId          int64  `json:"backupTaskId"`          // 备份任务id
	TaskId                string `json:"taskId"`                // 任务id
	BackupName            string `json:"backupName"`            // 备份名称
	ProdInstId            int64  `json:"prodInstId"`            // 实例id
	OuterProdInstId       string `json:"outerProdInstId"`       // 外部实例id
	ProdInstName          string `json:"prodInstName"`          // 实例名称
	Description           string `json:"description"`           // 备份描述
	StorageType           string `json:"storageType"`           // 存储类型（s3/disk/region_s3）
	OpType                string `json:"opType"`                // 操作类型（auto/manual）
	TaskType              string `json:"taskType"`              // 备份类型（full/incr）
	TaskStatus            int32  `json:"taskStatus"`            // 任务状态（100/101/102/1/-1）
	BackedUpDataSize      int64  `json:"backedUpDataSize"`      // 备份大小（字节）
	BackedUpDataSizeHuman string `json:"backedUpDataSizeHuman"` // 展示大小（自动适配单位）
	BackupStartTime       string `json:"backupStartTime"`       // 备份开始时间
	BackupEndTime         string `json:"backupEndTime"`         // 备份结束时间
	Disabled              bool   `json:"disabled"`              // 禁用备份
}
