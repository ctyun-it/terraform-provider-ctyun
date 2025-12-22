package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PostgresqlAddSecurityGroupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPostgresqlAddSecurityGroupApi(client *ctyunsdk.CtyunClient) *PostgresqlAddSecurityGroupApi {
	return &PostgresqlAddSecurityGroupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/securityGroup/add",
		},
	}
}

func (this *PostgresqlAddSecurityGroupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PostgresqlAddSecurityGroupRequest, header *PostgresqlAddSecurityGroupRequestHeader) (response *PostgresqlAddSecurityGroupResponse, err error) {
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

	if req.SecurityGroupId == "" {
		err = errors.New("missing required field: SecurityGroupId")
		return
	}
	if req.InstanceId == "" {
		err = errors.New("missing required field: InstanceName(实例名称)")
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response = &PostgresqlAddSecurityGroupResponse{}
	err = resp.Parse(response)
	if err != nil {
		return
	}
	return response, nil
}

type PostgresqlAddSecurityGroupRequest struct {
	SecurityGroupId string `json:"securityGroupId"` // 绑定安全组id
	InstanceId      string `json:"instanceId"`      // 实例id
}

type PostgresqlAddSecurityGroupRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
}

type PostgresqlAddSecurityGroupResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码，参考下方状态码
	Error      string `json:"error"`      // 错误码
	Message    string `json:"message"`    // 描述信息
}
