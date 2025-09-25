package mysql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbUnbindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbUnbindEipApi(client *ctyunsdk.CtyunClient) *TeledbUnbindEipApi {
	return &TeledbUnbindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/unbind",
		},
	}
}

func (this *TeledbUnbindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbUnbindEipRequest, header *TeledbUnbindEipRequestHeader) (unbindResponse *TeledbUnbindEipResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameCtdas, builder)
	if err != nil {
		return
	}
	unbindResponse = &TeledbUnbindEipResponse{}
	err = resp.Parse(unbindResponse)
	if err != nil {
		return
	}
	return unbindResponse, nil
}

type TeledbUnbindEipRequest struct {
	EipID  string `json:"eipId"`
	Eip    string `json:"eip"`
	InstID string `json:"instId"`
}

type TeledbUnbindEipRequestHeader struct {
	ProjectID *string `json:"project_id,omitempty"`
}

type TeledbUnbindEipResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}
