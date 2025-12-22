package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlLockAccountApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlLockAccountApi(client *ctyunsdk.CtyunClient) *PgsqlLockAccountApi {
	return &PgsqlLockAccountApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/PG/v1/inst-user/lock",
		},
	}
}

func (this *PgsqlLockAccountApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlLockAccountRequest, header *PgsqlLockAccountRequestHeader) (CGrantPrivilegeResp *PgsqlLockAccountResponse, err error) {
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
	CGrantPrivilegeResp = &PgsqlLockAccountResponse{}
	err = resp.Parse(CGrantPrivilegeResp)
	if err != nil {
		return
	}
	return CGrantPrivilegeResp, nil
}

type PgsqlLockAccountRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	Username   string `json:"username"`
}
type PgsqlLockAccountRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlLockAccountResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
