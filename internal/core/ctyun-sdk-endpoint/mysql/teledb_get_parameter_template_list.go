package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGetParameterTemplateListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGetParameterTemplateListApi(client *ctyunsdk.CtyunClient) *TeledbGetParameterTemplateListApi {
	return &TeledbGetParameterTemplateListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/RDS2/v1/open-api/parameter/parameter-groups",
		},
	}
}

func (this *TeledbGetParameterTemplateListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGetParameterTemplateListRequest, header *TeledbGetParameterTemplateListRequestHeader) (GetParameterTemplateListResp *TeledbGetParameterTemplateListResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	builder.AddHeader("regionId", header.RegionID)

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

	if req.ParameterGroupName != nil {
		builder.AddParam("parameterGroupName", *req.ParameterGroupName)
	}
	if req.EngineVersion != nil {
		builder.AddParam("engineVersion", *req.EngineVersion)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	GetParameterTemplateListResp = &TeledbGetParameterTemplateListResponse{}
	err = resp.Parse(GetParameterTemplateListResp)
	if err != nil {
		return
	}
	return GetParameterTemplateListResp, nil
}

type TeledbGetParameterTemplateListRequest struct {
	ParameterGroupName *string `json:"parameterGroupName"`
	EngineVersion      *string `json:"engineVersion"`
	PageNow            int32   `json:"pageNow"`
	PageSize           int32   `json:"pageSize"`
}
type TeledbGetParameterTemplateListRequestHeader struct {
	ProjectID *string `json:"Project-Id"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type TeledbGetParameterTemplateListResponse struct {
	StatusCode int32                                            `json:"statusCode"`      // 接口状态码
	Error      *string                                          `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                           `json:"message"`         // 描述信息
	ReturnObj  *TeledbGetParameterTemplateListResponseReturnObj `json:"returnObj"`
}
type TeledbGetParameterTemplateListResponseReturnObj struct {
	PageNum           int32                   `json:"pageNum"`           // 当前页
	PageSize          int32                   `json:"pageSize"`          // 每页的数量
	Size              int32                   `json:"size"`              // 当前页的数量
	StartRow          int32                   `json:"startRow"`          // 当前页面第一个元素在数据库中的行号
	EndRow            int32                   `json:"endRow"`            // 当前页面最后一个元素在数据库中的行号
	Total             int32                   `json:"total"`             // 总记录数
	Pages             int32                   `json:"pages"`             // 总页数
	PrePage           int32                   `json:"prePage"`           // 前一页
	IsFirstPage       bool                    `json:"isFirstPage"`       // 是否为第一页
	IsLastPage        bool                    `json:"isLastPage"`        // 是否为最后一页
	HasPreviousPage   bool                    `json:"hasPreviousPage"`   // 是否有前一页
	HasNextPage       bool                    `json:"hasNextPage"`       // 是否有下一页
	NavigatePages     int32                   `json:"navigatePages"`     // 导航页码数
	NavigatePageNums  []int32                 `json:"navigatepageNums"`  // 所有导航页号
	List              []ParameterTemplateInfo `json:"list"`              // 结果集
	NavigateLastPage  int32                   `json:"navigateLastPage"`  // 页面上显示的最后一个页码
	NavigateFirstPage int32                   `json:"navigateFirstPage"` // 页面显示的第一个页码
}

// 备份记录对象
type ParameterTemplateInfo struct {
	ParameterGroupName string `json:"parameterGroupName"`
	MysqlEngine        string `json:"mysqlEngine"`
	CreateTime         int64  `json:"createTime"`
	Restart            bool   `json:"restart"`
	Description        string `json:"description"`
	ID                 int64  `json:"id"`
	IsDefault          int32  `json:"isdefault"`
	UserId             int64  `json:"userId"`
}
