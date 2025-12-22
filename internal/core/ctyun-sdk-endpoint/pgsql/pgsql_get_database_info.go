package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGetDatabaseSchemaApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGetDatabaseSchemaApi(client *ctyunsdk.CtyunClient) *PgsqlGetDatabaseSchemaApi {
	return &PgsqlGetDatabaseSchemaApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/database/list",
		},
	}
}

func (this *PgsqlGetDatabaseSchemaApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGetDatabaseSchemaRequest, header *PgsqlGetDatabaseSchemaRequestHeader) (GetDatabaseSchemaResp *PgsqlGetDatabaseSchemaResponse, err error) {
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

	builder.AddParam("prodInstId", req.ProdInstId)
	if req.DBName != nil {
		builder.AddParam("dbName", *req.DBName)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	GetDatabaseSchemaResp = &PgsqlGetDatabaseSchemaResponse{}
	err = resp.Parse(GetDatabaseSchemaResp)
	if err != nil {
		return
	}
	return GetDatabaseSchemaResp, nil
}

type PgsqlGetDatabaseSchemaRequest struct {
	ProdInstId string  `json:"prodInstId"` // 外部实例ID，必填
	DBName     *string `json:"dbName,omitempty"`
}

type PgsqlGetDatabaseSchemaRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}

type PgsqlGetDatabaseSchemaResponse struct {
	StatusCode int32                                     `json:"statusCode"`      // 接口状态码
	Error      *string                                   `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string                                    `json:"message"`         // 描述信息
	ReturnObj  []PgsqlGetDatabaseSchemaResponseReturnObj `json:"returnObj"`
}

type PgsqlGetDatabaseSchemaResponseReturnObj struct {
	ProdInstId    string `json:"prodInstId"`
	DBName        string `json:"dbName"`
	DBEncoding    string `json:"dbEncoding"`
	DBCollate     string `json:"dbCollate"`
	DbType        string `json:"dbCtype"`
	DBOwner       string `json:"dbOwner"`
	DBDescription string `json:"dbDescription"`
}
