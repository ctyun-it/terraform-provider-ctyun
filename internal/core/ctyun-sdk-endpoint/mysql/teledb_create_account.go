package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateAccountApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCreateAccountApi(client *ctyunsdk.CtyunClient) *TeledbCreateAccountApi {
	return &TeledbCreateAccountApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/account",
		},
	}
}

func (this *TeledbCreateAccountApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateAccountRequest, header *TeledbCreateAccountRequestHeader) (CreateAccountResp *TeledbCreateAccountResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}
	if req.OuterProdInstId == "" || header.InstID == "" {
		err = errors.New("instId 为空")
		return
	}
	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("inst-id", header.InstID)
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CreateAccountResp = &TeledbCreateAccountResponse{}
	err = resp.Parse(CreateAccountResp)
	if err != nil {
		return
	}
	return CreateAccountResp, nil
}

type SchemaPrivilegeVO struct {
	GrantSchema  string `json:"grantSchema"`
	ReadOnly     *bool  `json:"readOnly,omitempty"`
	DDLPrivilege *bool  `json:"ddlPrivilege,omitempty"`
	DMLPrivilege *bool  `json:"dmlPrivilege,omitempty"`
	ReadAndWrite *bool  `json:"readAndWrite,omitempty"`
}

type TeledbCreateAccountRequest struct {
	OuterProdInstId       string              `json:"outerProdInstId"` // 实例ID，必填
	AccountName           string              `json:"accountName"`     // 账户名称
	AccountPassword       string              `json:"accountPassword"` // 账户密码（安全考虑需要用base64加密后传输）
	SchemaPrivilegeVOList []SchemaPrivilegeVO `json:"schemaPrivilegeVOList"`
}
type TeledbCreateAccountRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCreateAccountResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
