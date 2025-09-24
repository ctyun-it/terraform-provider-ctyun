package mysql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type TeledbBindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewTeledbBindEipApi(client *ctyunsdk.CtyunClient) *TeledbBindEipApi {
	return &TeledbBindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/bind",
		},
	}
}

func (this *TeledbBindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *TeledbBindEipRequest, header *TeledbBindEipRequestHeader) (bindResponse *TeledbBindEipResponse, err error) {
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
	bindResponse = &TeledbBindEipResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

type TeledbBindEipRequest struct {
	EipID  string `json:"eipId"`
	Eip    string `json:"eip"`
	InstID string `json:"instId"`
}

type TeledbBindEipRequestHeader struct {
	ProjectID *string `json:"project_id"`
}

type TeledbBindEipResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}
