package mysql

import (
	"context"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbDeleteParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbDeleteParameterTemplateApi(client *ctyunsdk.CtyunClient) *TeledbDeleteParameterTemplateApi {
	return &TeledbDeleteParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodDelete,
			UrlPath: "/RDS2/v1/open-api/parameter/parameter-group",
		},
	}
}

func (this *TeledbDeleteParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbDeleteParameterTemplateRequest, header *TeledbDeleteParameterTemplateRequestHeader) (DeleteParameterTemplateResp *TeledbDeleteParameterTemplateResponse, err error) {
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

	builder.AddParam("id", fmt.Sprintf("%d", req.ID))
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	DeleteParameterTemplateResp = &TeledbDeleteParameterTemplateResponse{}
	err = resp.Parse(DeleteParameterTemplateResp)
	if err != nil {
		return
	}
	return DeleteParameterTemplateResp, nil
}

type TeledbDeleteParameterTemplateRequest struct {
	ID int64 `json:"id"`
}
type TeledbDeleteParameterTemplateRequestHeader struct {
	ProjectID string `json:"Project-Id"`
	RegionID  string `json:"regionId"` // 资源池ID，必填
}
type TeledbDeleteParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
