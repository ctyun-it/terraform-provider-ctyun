package mysql

import (
	"context"
	"errors"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbResetRootPasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbResetRootPasswordApi(client *ctyunsdk.CtyunClient) *TeledbResetRootPasswordApi {
	return &TeledbResetRootPasswordApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/instance/reset-root-password",
		},
	}
}

func (this *TeledbResetRootPasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbResetRootPasswordRequest, header *TeledbResetRootPasswordRequestHeader) (ResetPasswordResp *TeledbResetRootPasswordResponse, err error) {
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
	ResetPasswordResp = &TeledbResetRootPasswordResponse{}
	err = resp.Parse(ResetPasswordResp)
	if err != nil {
		return
	}
	return ResetPasswordResp, nil
}

type TeledbResetRootPasswordRequest struct {
	OuterProdInstId string `json:"outerProdInstId"` // 实例ID，必填
	RootPassword    string `json:"rootPassword"`    // 账户密码（安全考虑需要用base64加密后传输）
}
type TeledbResetRootPasswordRequestHeader struct {
	ProjectID string `json:"projectID"`
	InstID    string `json:"instId"`    // 实例ID，必填
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbResetRootPasswordResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
