package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetParameterTemplateListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetParameterTemplateListApi(client *ctyunsdk.CtyunClient) *PgsqlGetParameterTemplateListApi {
	return &PgsqlGetParameterTemplateListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/template/page",
		},
	}
}

func (this *PgsqlGetParameterTemplateListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetParameterTemplateListRequest, header *PgsqlGetParameterTemplateListRequestHeader) (GetParameterTemplateListResp *PgsqlGetParameterTemplateListResponse, err error) {
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
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNow))

	if req.Name != nil {
		builder.AddParam("name", *req.Name)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetParameterTemplateListResp = &PgsqlGetParameterTemplateListResponse{}
	err = resp.Parse(GetParameterTemplateListResp)
	if err != nil {
		return
	}
	return GetParameterTemplateListResp, nil
}

type PgsqlGetParameterTemplateListRequest struct {
	Name     *string `json:"name"`
	PageNow  int32   `json:"pageNow"`
	PageSize int32   `json:"pageSize"`
}
type PgsqlGetParameterTemplateListRequestHeader struct {
	ProjectID *string `json:"Project-Id"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGetParameterTemplateListResponse struct {
	StatusCode int32                                           `json:"statusCode"`      // 接口状态码
	Error      *string                                         `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                          `json:"message"`         // 描述信息
	ReturnObj  *PgsqlGetParameterTemplateListResponseReturnObj `json:"returnObj"`
}
type PgsqlGetParameterTemplateListResponseReturnObj struct {
	PageNum  int32                   `json:"pageNum"`  // 当前页
	PageSize int32                   `json:"pageSize"` // 每页的数量
	Size     int32                   `json:"size"`     // 当前页的数量
	Total    int32                   `json:"total"`    // 总记录数
	List     []ParameterTemplateInfo `json:"list"`     // 结果集
}

// 备份记录对象
type ParameterTemplateInfo struct {
	PgTemplateId    int64  `json:"pgTemplateId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Version         string `json:"version"`
	Modify          bool   `json:"modify"`
	UpdateTimestamp string `json:"updateTimestamp"`
}
