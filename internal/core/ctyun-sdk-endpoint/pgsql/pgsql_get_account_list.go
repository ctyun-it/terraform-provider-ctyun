package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetAccountListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetAccountListApi(client *ctyunsdk.CtyunClient) *PgsqlGetAccountListApi {
	return &PgsqlGetAccountListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/inst-user/users",
		},
	}
}

func (this *PgsqlGetAccountListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetAccountListRequest, header *PgsqlGetAccountListRequestHeader) (GetAccountInfoResp *PgsqlGetAccountListResponse, err error) {
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
	if req.Username != nil {
		builder.AddParam("username", *req.Username)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetAccountInfoResp = &PgsqlGetAccountListResponse{}
	err = resp.Parse(GetAccountInfoResp)
	if err != nil {
		return
	}
	return GetAccountInfoResp, nil
}

type PgsqlGetAccountListRequest struct {
	ProdInstId string  `json:"prodInstId"` // 实例ID，必填
	PageNum    int32   `json:"pageNum"`
	PageSize   int32   `json:"pageSize"`
	Username   *string `json:"username,omitempty"`
}
type PgsqlGetAccountListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGetAccountListResponseReturnObjList struct {
	List      []PgsqlGetAccountListResponseReturnObj `json:"list"`
	PageNum   int32                                  `json:"pageNum"`
	PageSize  int32                                  `json:"pageSize"`
	PageTotal int32                                  `json:"pageTotal"`
	Total     int32                                  `json:"total"`
}

type PgsqlGetAccountListResponse struct {
	StatusCode int32                                     `json:"statusCode"`      // 接口状态码
	Error      *string                                   `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                    `json:"message"`         // 描述信息
	ReturnObj  *PgsqlGetAccountListResponseReturnObjList `json:"returnObj"`
}

type PgsqlGetAccountListResponseReturnObj struct {
	Username      string `json:"username"`
	RolSuper      bool   `json:"rolsuper"`
	RolInherit    bool   `json:"rolinherit"`
	RolCreateRole bool   `json:"rolcreaterole"`
	RolCreateDB   bool   `json:"rolcreatedb"`
	RolCanLogin   bool   `json:"rolcanlogin"`
	RolConnLimit  int32  `json:"rolconnlimit"`
	RolByPassRls  bool   `json:"rolbypassrls"` // 用户是否绕过每个行级安全策略
}
