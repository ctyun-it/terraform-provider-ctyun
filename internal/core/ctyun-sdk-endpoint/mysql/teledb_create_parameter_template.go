package mysql

import (
	"context"
	"fmt"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbCreateParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	template core.CtyunRequestTemplate
	client   *ctyunsdk.CtyunClient
}

func NewTeledbCreateParameterTemplateApi(client *ctyunsdk.CtyunClient) *TeledbCreateParameterTemplateApi {
	return &TeledbCreateParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/RDS2/v1/open-api/parameter/parameter-group",
		},
	}
}

func (this *TeledbCreateParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbCreateParameterTemplateRequest, header *TeledbCreateParameterTemplateRequestHeader) (CreateParameterTemplateResp *TeledbCreateParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
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
	CreateParameterTemplateResp = &TeledbCreateParameterTemplateResponse{}
	err = resp.Parse(CreateParameterTemplateResp)
	if err != nil {
		return
	}
	return CreateParameterTemplateResp, nil
}

type TeledbCreateParameterTemplateRequest struct {
	ParameterGroupName string  `json:"parameterGroupName"`    // 实例ID，必填
	Engine             string  `json:"engine"`                // 账户名称
	Description        *string `json:"description,omitempty"` // 账户密码（安全考虑需要用base64加密后传输）
}
type TeledbCreateParameterTemplateRequestHeader struct {
	ProjectID *string `json:"Project-Id,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type TeledbCreateParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
