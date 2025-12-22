package mysql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCopyParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbCopyParameterTemplateApi(client *ctyunsdk.CtyunClient) *TeledbCopyParameterTemplateApi {
	return &TeledbCopyParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/parameter/clone-parameter-group",
		},
	}
}

func (this *TeledbCopyParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCopyParameterTemplateRequest, header *TeledbCopyParameterTemplateRequestHeader) (CopyParameterTemplateResp *TeledbCopyParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != "" {
		builder.AddHeader("project-id", header.ProjectID)
	}

	if header.RegionID == "" {
		err = fmt.Errorf("region_id is required")
		return
	}
	builder.AddHeader("regionId", header.RegionID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	CopyParameterTemplateResp = &TeledbCopyParameterTemplateResponse{}
	err = resp.Parse(CopyParameterTemplateResp)
	if err != nil {
		return
	}
	return CopyParameterTemplateResp, nil
}

type TeledbCopyParameterTemplateRequest struct {
	SourceParameterGroupId int64   `json:"sourceParameterGroupId"`       // 实例ID，必填
	ParameterGroupName     string  `json:"parameterGroupName"`           // 账户名称
	ParameterGroupDesc     *string `json:"parameterGroupDesc,omitempty"` // 账户密码（安全考虑需要用base64加密后传输）
}
type TeledbCopyParameterTemplateRequestHeader struct {
	ProjectID string `json:"projectID"`
	RegionID  string `json:"region_id"` // 资源池ID，必填
}
type TeledbCopyParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
