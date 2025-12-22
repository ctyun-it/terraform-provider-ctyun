package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUpdateParameterTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUpdateParameterTemplateApi(client *ctyunsdk.CtyunClient) *TeledbUpdateParameterTemplateApi {
	return &TeledbUpdateParameterTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPut,
			UrlPath: "/RDS2/v1/open-api/parameter/modify-parameter",
		},
	}
}

func (this *TeledbUpdateParameterTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUpdateParameterTemplateRequest, header *TeledbUpdateParameterTemplateRequestHeader) (UpdateParameterTemplateResp *TeledbUpdateParameterTemplateResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}

	builder.AddHeader("regionId", header.RegionID)

	builder.AddParam("id", fmt.Sprintf("%d", req.ID))

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	UpdateParameterTemplateResp = &TeledbUpdateParameterTemplateResponse{}
	err = resp.Parse(UpdateParameterTemplateResp)
	if err != nil {
		return
	}
	return UpdateParameterTemplateResp, nil
}

func valueToString(value []ParameterObj) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type ParameterObj struct {
	OldValue       string `json:"oldValue"`
	ParameterName  string `json:"parameterName"`
	ParameterValue string `json:"parameterValue"`
	Restart        string `json:"restart,omitempty"`
}

type TeledbUpdateParameterTemplateRequest struct {
	ID    int64          `json:"id"`
	Value []ParameterObj `json:"value"`
}

type TeledbUpdateParameterTemplateRequestHeader struct {
	ProjectID *string `json:"Project-Id,omitempty"`
	RegionID  string  `json:"regionId"` // 资源池ID，必填
}
type TeledbUpdateParameterTemplateResponse struct {
	StatusCode int32   `json:"statusCode"`      // 接口状态码
	Error      *string `json:"error,omitempty"` // 错误码，失败时返回，成功时为空
	Message    string  `json:"message"`         // 描述信息
}
