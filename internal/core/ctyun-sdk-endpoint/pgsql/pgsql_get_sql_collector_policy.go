package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PostgresqlGetCollectorPolicyApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPostgresqlGetCollectorPolicyApi(client *ctyunsdk.CtyunClient) *PostgresqlGetCollectorPolicyApi {
	return &PostgresqlGetCollectorPolicyApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/PG/v1/collector/sql-collector-policy",
		},
	}
}

func (this *PostgresqlGetCollectorPolicyApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PostgresqlGetCollectorPolicyRequest, header *PostgresqlGetCollectorPolicyRequestHeader) (response *PostgresqlGetCollectorPolicyResponse, err error) {
	builder := this.WithCredential(&credential)

	// 添加Query参数
	builder.AddParam("prodInstId", req.ProdInstId)

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
	response = &PostgresqlGetCollectorPolicyResponse{}
	err = resp.Parse(response)
	if err != nil {
		return
	}
	return response, nil
}

type PostgresqlGetCollectorPolicyRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例id
}

type PostgresqlGetCollectorPolicyRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
	RegionId  *string `json:"regionId,omitempty"` // 资源池id
}

type PostgresqlGetCollectorPolicyResponseReturnObj struct {
	InstName           string `json:"instName"`           // 实例名称
	ProdInstId         string `json:"prodInstId"`         // 实例id
	SqlCollectorStatus string `json:"sqlCollectorStatus"` // SQL审计状态，取值：enable | disabled
	LogInterval        int32  `json:"logInterval"`        // 定时收集SQL日志的间隔.单位：分钟
}

type PostgresqlGetCollectorPolicyResponse struct {
	StatusCode int32                                          `json:"statusCode"`
	Error      string                                         `json:"error"`     // 错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    string                                         `json:"message"`   // 描述信息
	ReturnObj  *PostgresqlGetCollectorPolicyResponseReturnObj `json:"returnObj"` // 返回对象
}
