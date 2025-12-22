package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetDatabaseSchemaListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetDatabaseSchemaListApi(client *ctyunsdk.CtyunClient) *PgsqlGetDatabaseSchemaListApi {
	return &PgsqlGetDatabaseSchemaListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/database/describe-schema",
		},
	}
}

func (this *PgsqlGetDatabaseSchemaListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetDatabaseSchemaListRequest, header *PgsqlGetDatabaseSchemaListRequestHeader) (GetDatabaseSchemaListResp *PgsqlGetDatabaseSchemaListResponse, err error) {
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
	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("outerProdInstId", req.ProdInstId)
	if req.SchemaName != nil {
		builder.AddHeader("dbName", *req.SchemaName)
	}
	builder.AddParam("pageNum", fmt.Sprintf("%d", req.PageNum))
	builder.AddParam("pageSize", fmt.Sprintf("%d", req.PageSize))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetDatabaseSchemaListResp = &PgsqlGetDatabaseSchemaListResponse{}
	err = resp.Parse(GetDatabaseSchemaListResp)
	if err != nil {
		return
	}
	return GetDatabaseSchemaListResp, nil
}

type PgsqlGetDatabaseSchemaListRequest struct {
	ProdInstId string  `json:"prodInstId"` // 外部实例ID，必填
	SchemaName *string `json:"schemaName,omitempty"`
	PageNum    int32   `json:"pageNum"`
	PageSize   int32   `json:"pageSize"`
}

type PgsqlGetDatabaseSchemaListRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlGetDatabaseSchemaListResponse struct {
	StatusCode int32                                         `json:"statusCode"`      // 接口状态码
	Error      *string                                       `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                        `json:"message"`         // 描述信息
	ReturnObj  []PgsqlGetDatabaseSchemaListResponseReturnObj `json:"returnObj"`
}

type PgsqlGetDatabaseSchemaListResponseReturnObj struct {
	ProdInstId    string `json:"prodInstId"`
	DBName        string `json:"dbName"`
	DBEncoding    string `json:"dbEncoding"`
	DBCollate     string `json:"dbCollate"`
	DbType        string `json:"dbCtype"`
	DBOwner       string `json:"dbOwner"`
	DBDescription string `json:"dbDescription"`
}
