package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlResetRootPasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlResetRootPasswordApi(client *ctyunsdk.CtyunClient) *PgsqlResetRootPasswordApi {
	return &PgsqlResetRootPasswordApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/inst-user/reset-root-password",
		},
	}
}

func (this *PgsqlResetRootPasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlResetRootPasswordRequest, header *PgsqlResetRootPasswordRequestHeader) (ResetRootPasswordResp *PgsqlResetRootPasswordResponse, err error) {
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
	ResetRootPasswordResp = &PgsqlResetRootPasswordResponse{}
	err = resp.Parse(ResetRootPasswordResp)
	if err != nil {
		return
	}
	return ResetRootPasswordResp, nil
}

type PgsqlResetRootPasswordRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	Password   string `json:"password"`   // 账户密码
}
type PgsqlResetRootPasswordRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlResetRootPasswordResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
