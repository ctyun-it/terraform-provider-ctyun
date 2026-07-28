package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetBackupRecoveryListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetBackupRecoveryListApi(client *ctyunsdk.CtyunClient) *TeledbGetBackupRecoveryListApi {
	return &TeledbGetBackupRecoveryListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v2/open-api/recovery/record",
		},
	}
}

func (this *TeledbGetBackupRecoveryListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetBackupRecoveryListRequest, header *TeledbGetBackupRecoveryListRequestHeader) (GetBackupRecoveryListResp *TeledbGetBackupRecoveryListResponse, err error) {
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
	builder.AddHeader("outerProdInstId", req.OuterProdInstId)
	if req.TaskID != nil {
		builder.AddHeader("taskId", *req.TaskID)
	}
	if req.ProdInstName != nil {
		builder.AddHeader("prodInstName", *req.ProdInstName)
	}
	if req.StartTime != nil {
		builder.AddHeader("startTime", *req.StartTime)
	}
	if req.EndTime != nil {
		builder.AddHeader("endTime", *req.EndTime)
	}
	if req.ID != nil {
		builder.AddHeader("id", fmt.Sprintf("%d", *req.ID))
	}
	if req.PageNow == 0 {
		err = errors.New("page_no 为空")
		return
	}
	if req.PageSize == 0 {
		err = errors.New("page_size 为空")
		return
	}

	builder.AddHeader("pageSize", fmt.Sprintf("%d", req.PageSize))
	builder.AddHeader("pageNow", fmt.Sprintf("%d", req.PageNow))
	builder.AddParam("outerProdInstId", req.OuterProdInstId)
	if req.TaskID != nil {
		builder.AddParam("taskId", *req.TaskID)
	}
	if req.ProdInstName != nil {
		builder.AddParam("prodInstName", *req.ProdInstName)
	}
	if req.StartTime != nil {
		builder.AddParam("startTime", *req.StartTime)
	}
	if req.EndTime != nil {
		builder.AddParam("endTime", *req.EndTime)
	}
	if req.ID != nil {
		builder.AddParam("id", fmt.Sprintf("%d", *req.ID))
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetBackupRecoveryListResp = &TeledbGetBackupRecoveryListResponse{}
	err = resp.Parse(GetBackupRecoveryListResp)
	if err != nil {
		return
	}
	return GetBackupRecoveryListResp, nil
}

type TeledbGetBackupRecoveryListRequest struct {
	OuterProdInstId string  `json:"outerProdInstId"` // 外部实例ID，必填
	TaskID          *string `json:"taskId"`
	ProdInstName    *string `json:"prodInstName,omitempty"` // 实例名
	StartTime       *string `json:"startTime,omitempty"`    // 开始时间
	EndTime         *string `json:"endTime,omitempty"`      // 结束时间
	PageNow         int32   `json:"pageNow"`                // 当前页码，必填
	PageSize        int32   `json:"pageSize"`               // 分页大小，必填
	ID              *int64  `json:"id"`                     // 恢复任务主键Id
}

type TeledbGetBackupRecoveryListRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}

type TeledbGetBackupRecoveryListResponse struct {
	StatusCode int32                                         `json:"statusCode"`      // 接口状态码
	Error      *string                                       `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                        `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetBackupRecoveryListResponseReturnObj `json:"returnObj"`
}
type BrRecoveryRecordVo struct {
	RecoveryTimepoint  string `json:"recoveryTimepoint"`
	SrcProdInstId      int64  `json:"srcProdInstId"`
	DstProdInstId      int64  `json:"dstProdInstId"`
	DstProdInstName    string `json:"dstProdInstName"`
	ProdInstName       string `json:"prodInstName"`
	OuterDstProdInstId string `json:"outerDstProdInstId"`
	OuterSrcProdInstId string `json:"outerSrcProdInstId"`
	RecoveryScope      string `json:"recoveryScope"`
	OuterProdInstId    string `json:"outerProdInstId"`
	SrcProdInstName    string `json:"srcProdInstName"`
	ID                 int64  `json:"id"`
	StartTime          string `json:"startTime"`
	EndTime            string `json:"endTime"`
	TaskID             string `json:"taskId"`
	TaskStatus         int32  `json:"taskStatus"` // 任务状态：1进行中，2成功，3失败
	RecoveryName       string `json:"recoveryName"`
}

type TeledbGetBackupRecoveryListResponseReturnObj struct {
	PageNum           int32                `json:"pageNum"`           // 当前页
	PageSize          int32                `json:"pageSize"`          // 每页的数量
	Size              int32                `json:"size"`              // 当前页的数量
	StartRow          int32                `json:"startRow"`          // 当前页面第一个元素在数据库中的行号
	EndRow            int32                `json:"endRow"`            // 当前页面最后一个元素在数据库中的行号
	Total             int32                `json:"total"`             // 总记录数
	Pages             int32                `json:"pages"`             // 总页数
	PrePage           int32                `json:"prePage"`           // 前一页
	IsFirstPage       bool                 `json:"isFirstPage"`       // 是否为第一页
	IsLastPage        bool                 `json:"isLastPage"`        // 是否为最后一页
	HasPreviousPage   bool                 `json:"hasPreviousPage"`   // 是否有前一页
	HasNextPage       bool                 `json:"hasNextPage"`       // 是否有下一页
	NavigatePages     int32                `json:"navigatePages"`     // 导航页码数
	NavigatePageNums  []int32              `json:"navigatepageNums"`  // 所有导航页号
	List              []BrRecoveryRecordVo `json:"list"`              // 结果集
	NavigateLastPage  int32                `json:"navigateLastPage"`  // 页面上显示的最后一个页码
	NavigateFirstPage int32                `json:"navigateFirstPage"` // 页面显示的第一个页码
}
