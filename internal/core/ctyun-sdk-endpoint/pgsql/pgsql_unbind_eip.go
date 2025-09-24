package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlUnBindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlUnBindEipApi(client *ctyunsdk.CtyunClient) *PgsqlUnBindEipApi {
	return &PgsqlUnBindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/unbind",
		},
	}
}

func (this *PgsqlUnBindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlUnBindEipRequest, header *PgsqlUnBindEipRequestHeader) (createResponse *PgsqlUnBindEipResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectId != nil {
		builder.AddHeader("Project-Id", *header.ProjectId)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNamePgSql, builder)
	if err != nil {
		return
	}
	response := PgsqlUnBindEipResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

type PgsqlUnBindEipRequest struct {
	EipID  string `json:"eipId"`  // 弹性id
	Eip    string `json:"eip"`    // 弹性ip
	InstID string `json:"instId"` // 实例id
}

type PgsqlUnBindEipRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
}

type PgsqlUnBindEipResponse struct {
	StatusCode int32   `json:"statusCode"`
	Error      *string `json:"error,omitempty"`   //错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    *string `json:"message,omitempty"` //描述信息
}
