package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlGrantPrivilegeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlGrantPrivilegeApi(client *ctyunsdk.CtyunClient) *PgsqlGrantPrivilegeApi {
	return &PgsqlGrantPrivilegeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/inst-user/grant-privilege",
		},
	}
}

func (this *PgsqlGrantPrivilegeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlGrantPrivilegeRequest, header *PgsqlGrantPrivilegeRequestHeader) (CGrantPrivilegeResp *PgsqlGrantPrivilegeResponse, err error) {
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

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	CGrantPrivilegeResp = &PgsqlGrantPrivilegeResponse{}
	err = resp.Parse(CGrantPrivilegeResp)
	if err != nil {
		return
	}
	return CGrantPrivilegeResp, nil
}

type PgsqlGrantPrivilegeRequest struct {
	ProdInstId    string  `json:"prodInstId"` // 实例ID，必填
	DbName        string  `json:"dbName"`     // 账户名称
	SchemaName    *string `json:"schemaName,omitempty"`
	Username      string  `json:"username"`
	UserPrivilege string  `json:"userPrivilege"`
}
type PgsqlGrantPrivilegeRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlGrantPrivilegeResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
