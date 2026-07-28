package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbRevokePrivilegeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbRevokePrivilegeApi(client *ctyunsdk.CtyunClient) *TeledbRevokePrivilegeApi {
	return &TeledbRevokePrivilegeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/account/operator-privilege-revoke",
		},
	}
}

func (this *TeledbRevokePrivilegeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbRevokePrivilegeRequest, header *TeledbRevokePrivilegeRequestHeader) (CRevokePrivilegeResp *TeledbRevokePrivilegeResponse, err error) {
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
	CRevokePrivilegeResp = &TeledbRevokePrivilegeResponse{}
	err = resp.Parse(CRevokePrivilegeResp)
	if err != nil {
		return
	}
	return CRevokePrivilegeResp, nil
}

type InvokePrivilegeVo struct {
	GrantSchema  string `json:"grantSchema"`
	DDLPrivilege *bool  `json:"ddlPrivilege,omitempty"`
	DMLPrivilege *bool  `json:"dmlPrivilege,omitempty"`
}
type TeledbRevokePrivilegeRequest struct {
	OuterProdInstId       string              `json:"outerProdInstId"` // 实例ID，必填
	AccountName           string              `json:"accountName"`     // 账户名称
	SchemaPrivilegeVOList []InvokePrivilegeVo `json:"schemaPrivilegeVOList"`
}
type TeledbRevokePrivilegeRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbRevokePrivilegeResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
