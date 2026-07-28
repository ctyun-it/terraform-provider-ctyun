package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PostgresqlCollectorRetentionApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPostgresqlCollectorRetentionApi(client *ctyunsdk.CtyunClient) *PostgresqlCollectorRetentionApi {
	return &PostgresqlCollectorRetentionApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/collector/sql-collector-retention",
		},
	}
}

func (this *PostgresqlCollectorRetentionApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PostgresqlCollectorRetentionRequest, header *PostgresqlCollectorRetentionRequestHeader) (response *PostgresqlCollectorRetentionResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}

	// 添加请求头参数
	if header.ProjectId != nil {
		builder.AddHeader("Project-Id", *header.ProjectId)
	}
	if header.RegionId == nil {
		err = errors.New("missing required field: RegionId")
		return
	}
	builder.AddHeader("regionId", *header.RegionId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response = &PostgresqlCollectorRetentionResponse{}
	err = resp.Parse(response)
	if err != nil {
		return
	}
	return response, nil
}

type PostgresqlCollectorRetentionRequest struct {
	ProdInstId  string `json:"prodInstId"`  // 实例id
	RetainValue int32  `json:"retainValue"` // SQL洞察日志保存时长，单位：天，取值：1,2,3
}

type PostgresqlCollectorRetentionRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
	RegionId  *string `json:"regionId,omitempty"` // 资源池id
}

type PostgresqlCollectorRetentionResponse struct {
	StatusCode int32    `json:"statusCode"`
	Error      string   `json:"error"`     // 错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    string   `json:"message"`   // 描述信息
	ReturnObj  struct{} `json:"returnObj"` // 返回对象
}
