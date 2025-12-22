package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDeleteDatabaseApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDeleteDatabaseApi(client *ctyunsdk.CtyunClient) *PgsqlDeleteDatabaseApi {
	return &PgsqlDeleteDatabaseApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/database/drop",
		},
	}
}

func (this *PgsqlDeleteDatabaseApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDeleteDatabaseRequest, header *PgsqlDeleteDatabaseRequestHeader) (DeleteDatabaseResp *PgsqlDeleteDatabaseResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
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

	//builder.AddParam("dbName", fmt.Sprintf("%s", req.DBName))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	DeleteDatabaseResp = &PgsqlDeleteDatabaseResponse{}
	err = resp.Parse(DeleteDatabaseResp)
	if err != nil {
		return
	}
	return DeleteDatabaseResp, nil
}

type PgsqlDeleteDatabaseRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	DBName     string `json:"dbName"`
	//SchemaName string `json:"schemaName"`
}
type PgsqlDeleteDatabaseRequestHeader struct {
	ProjectID string `json:"projectID"`
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type PgsqlDeleteDatabaseResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
