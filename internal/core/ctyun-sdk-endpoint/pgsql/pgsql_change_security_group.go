package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PostgresqlChangeSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPostgresqlChangeSecurityGroupApi(client *ctyunsdk.CtyunClient) *PostgresqlChangeSecurityGroupApi {
	return &PostgresqlChangeSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup/change",
		},
	}
}

func (this *PostgresqlChangeSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PostgresqlChangeSecurityGroupRequest, header *PostgresqlChangeSecurityGroupRequestHeader) (response *PostgresqlChangeSecurityGroupResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}

	// 添加请求头参数
	if header.ProjectId == nil {
		err = errors.New("missing required field: ProjectId")
		return
	}
	builder.AddHeader("project-id", *header.ProjectId)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response = &PostgresqlChangeSecurityGroupResponse{}
	err = resp.Parse(response)
	if err != nil {
		return
	}
	return response, nil
}

type PostgresqlChangeSecurityGroupRequest struct {
	SecurityGroupId    string `json:"securityGroupId"`    // 原安全组id
	InstanceId         string `json:"instanceId"`         // 实例id
	NewSecurityGroupId string `json:"newSecurityGroupId"` // 新安全组id
}

type PostgresqlChangeSecurityGroupRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
}

type PostgresqlChangeSecurityGroupResponse struct {
	StatusCode int32    `json:"statusCode"` // 接口状态码，参考下方状态码
	Error      string   `json:"error"`      // 错误码
	Message    string   `json:"message"`    // 描述信息
	ReturnObj  struct{} `json:"returnObj"`  // 返回对象
}
