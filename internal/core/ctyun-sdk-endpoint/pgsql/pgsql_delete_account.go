package pgsql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlDeleteAccountApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlDeleteAccountApi(client *ctyunsdk.CtyunClient) *PgsqlDeleteAccountApi {
	return &PgsqlDeleteAccountApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "PG/v1/inst-user/drop",
		},
	}
}

func (this *PgsqlDeleteAccountApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlDeleteAccountRequest, header *PgsqlDeleteAccountRequestHeader) (DeleteAccountResp *PgsqlDeleteAccountResponse, err error) {
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
	if req.Username == "" {
		err = errors.New("user_name is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	DeleteAccountResp = &PgsqlDeleteAccountResponse{}
	err = resp.Parse(DeleteAccountResp)
	if err != nil {
		return
	}
	return DeleteAccountResp, nil
}

type PgsqlDeleteAccountRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID，必填
	Username   string `json:"username"`   // 账户名称
}
type PgsqlDeleteAccountRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlDeleteAccountResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
