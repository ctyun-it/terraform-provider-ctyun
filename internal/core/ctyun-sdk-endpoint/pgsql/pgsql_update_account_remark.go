package pgsql

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUpdateAccountRemarkApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUpdateAccountRemarkApi(client *ctyunsdk.CtyunClient) *PgsqlUpdateAccountRemarkApi {
	return &PgsqlUpdateAccountRemarkApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/PG/v1/inst-user/description",
		},
	}
}

func (this *PgsqlUpdateAccountRemarkApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUpdateAccountRemarkRequest, header *PgsqlUpdateAccountRemarkRequestHeader) (UpdateAccountRemarkResp *PgsqlUpdateAccountRemarkResponse, err error) {
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
	builder.AddHeader("regionId", header.RegionID)
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	UpdateAccountRemarkResp = &PgsqlUpdateAccountRemarkResponse{}
	err = resp.Parse(UpdateAccountRemarkResp)
	if err != nil {
		return
	}
	return UpdateAccountRemarkResp, nil
}

type PgsqlUpdateAccountRemarkRequest struct {
	ProdInstId  string  `json:"prodInstId"` // 实例ID，必填
	Username    string  `json:"username"`   // 账户名称
	Description *string `json:"description"`
}
type PgsqlUpdateAccountRemarkRequestHeader struct {
	ProjectID *string `json:"Project-Id	,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type PgsqlUpdateAccountRemarkResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
