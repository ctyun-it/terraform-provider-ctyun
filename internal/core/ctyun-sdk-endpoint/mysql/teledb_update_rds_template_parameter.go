package mysql

import (
	"context"
	"encoding/json"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateRdsTemplateParameterApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateRdsTemplateParameterApi(client *ctyunsdk.CtyunClient) *TeledbUpdateRdsTemplateParameterApi {
	return &TeledbUpdateRdsTemplateParameterApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/parameter/rds-parameter",
		},
	}
}

func (this *TeledbUpdateRdsTemplateParameterApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateRdsTemplateParameterRequest, header *TeledbUpdateRdsTemplateParameterRequestHeader) (UpdateParameterTemplateResp *TeledbUpdateRdsTemplateParameterResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)
	builder.AddHeader("inst-id", header.InstID)

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	UpdateParameterTemplateResp = &TeledbUpdateRdsTemplateParameterResponse{}
	err = resp.Parse(UpdateParameterTemplateResp)
	if err != nil {
		return
	}
	return UpdateParameterTemplateResp, nil
}

func mapToString(value map[string]string) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type TeledbUpdateRdsTemplateParameterRequest struct {
	OuterProdInstId string             `json:"outerProdInstId"`
	ID              *int64             `json:"id,omitempty"`
	Parameters      *map[string]string `json:"parameters,omitempty"`
}

type TeledbUpdateRdsTemplateParameterRequestHeader struct {
	ProjectID *string `json:"Project-Id,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
	InstID    string  `json:"inst-id"`
}
type TeledbUpdateRdsTemplateParameterResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
