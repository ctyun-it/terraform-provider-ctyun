package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PostgresqlCollectorPolicyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPostgresqlCollectorPolicyApi(client *ctyunsdk.CtyunClient) *PostgresqlCollectorPolicyApi {
	return &PostgresqlCollectorPolicyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/collector/sql-collector-policy",
		},
	}
}

func (this *PostgresqlCollectorPolicyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PostgresqlCollectorPolicyRequest, header *PostgresqlCollectorPolicyRequestHeader) (createResponse *PostgresqlCollectorPolicyResponse, err error) {
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
		err = errors.New("missing required field: RegionID")
		return
	}
	builder.AddHeader("regionId", *header.RegionId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response := PostgresqlCollectorPolicyResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

type PostgresqlCollectorPolicyRequest struct {
	ProdInstId         string `json:"prodInstId"`            // 实例id
	SqlCollectorStatus string `json:"sqlCollectorStatus"`    // 开启或关闭SQL洞察（SQL审计），取值：enable | disabled
	LogInterval        *int32 `json:"logInterval,omitempty"` // 定时收集SQL日志的间隔.单位：分钟，默认5分钟，取值：5，10，30，60
}

type PostgresqlCollectorPolicyRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
	RegionId  *string `json:"regionId,omitempty"` // 资源池id
}

type PostgresqlCollectorPolicyResponse struct {
	StatusCode int32    `json:"statusCode"`
	Error      string   `json:"error"`     // 错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    string   `json:"message"`   // 描述信息
	ReturnObj  struct{} `json:"returnObj"` // 返回对象，根据示例为空对象
}
