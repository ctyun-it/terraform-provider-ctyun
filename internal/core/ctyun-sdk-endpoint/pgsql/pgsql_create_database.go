package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlCreateDatabaseApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlCreateDatabaseApi(client *ctyunsdk.CtyunClient) *PgsqlCreateDatabaseApi {
	return &PgsqlCreateDatabaseApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/database/create",
		},
	}
}

func (this *PgsqlCreateDatabaseApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlCreateDatabaseRequest, header *PgsqlCreateDatabaseRequestHeader) (CreateDatabaseResp *PgsqlCreateDatabaseResponse, err error) {
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

	if req.DBName == "" {
		err = fmt.Errorf("database_name is required")
		return
	}

	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CreateDatabaseResp = &PgsqlCreateDatabaseResponse{}
	err = resp.Parse(CreateDatabaseResp)
	if err != nil {
		return
	}
	return CreateDatabaseResp, nil
}

type PgsqlCreateDatabaseRequest struct {
	ProdInstId    string  `json:"prodInstId"` // 实例ID，必填
	DBName        string  `json:"dbName"`     // 数据库名称
	DBEncoding    string  `json:"dbEncoding"` // 字符集
	DBCollate     *string `json:"dbCollate,omitempty"`
	DBType        *string `json:"dbCtype,omitempty"`
	DBOwner       *string `json:"dbOwner,omitempty"`
	DBDescription *string `json:"dbDescription,omitempty"`
}

type PgsqlCreateDatabaseRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlCreateDatabaseResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
	//ReturnObj  *PgsqlCreateDatabaseResponseReturnObj `json:"returnObj"`
}
