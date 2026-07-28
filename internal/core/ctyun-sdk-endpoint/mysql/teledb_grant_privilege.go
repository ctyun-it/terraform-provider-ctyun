package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbGrantPrivilegeApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbGrantPrivilegeApi(client *ctyunsdk.CtyunClient) *TeledbGrantPrivilegeApi {
	return &TeledbGrantPrivilegeApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/account/privilege-grant",
		},
	}
}

func (this *TeledbGrantPrivilegeApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbGrantPrivilegeRequest, header *TeledbGrantPrivilegeRequestHeader) (CGrantPrivilegeResp *TeledbGrantPrivilegeResponse, err error) {
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
	CGrantPrivilegeResp = &TeledbGrantPrivilegeResponse{}
	err = resp.Parse(CGrantPrivilegeResp)
	if err != nil {
		return
	}
	return CGrantPrivilegeResp, nil
}

type TeledbGrantPrivilegeRequest struct {
	OuterProdInstId       string              `json:"outerProdInstId"` // 实例ID，必填
	AccountName           string              `json:"accountName"`     // 账户名称
	SchemaPrivilegeVOList []SchemaPrivilegeVO `json:"schemaPrivilegeVOList"`
}
type TeledbGrantPrivilegeRequestHeader struct {
	ProjectID string `json:"project-id"`
	InstID    string `json:"instId"`   // 实例ID，必填
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbGrantPrivilegeResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
