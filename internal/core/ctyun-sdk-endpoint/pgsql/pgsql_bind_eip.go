package pgsql

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type PgsqlBindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewPgsqlBindEipApi(client *ctyunsdk.CtyunClient) *PgsqlBindEipApi {
	return &PgsqlBindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/bind",
		},
	}
}

func (this *PgsqlBindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *PgsqlBindEipRequest, header *PgsqlBindEipRequestHeader) (createResponse *PgsqlBindEipResponse, err error) {
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
	response := PgsqlBindEipResponse{}
	err = resp.Parse(&response)
	if err != nil {
		return
	}
	return &response, nil
}

type PgsqlBindEipRequest struct {
	EipID  string `json:"eipId"`  // 弹性id
	Eip    string `json:"eip"`    // 弹性ip
	InstID string `json:"instId"` // 实例id
}

type PgsqlBindEipRequestHeader struct {
	ProjectId *string `json:"projectId,omitempty"`
}

type PgsqlBindEipResponse struct {
	StatusCode int32   `json:"statusCode"`
	Error      *string `json:"error,omitempty"`   //错误码。当接口失败时才返回具体错误编码，成功不返回或者为空
	Message    *string `json:"message,omitempty"` //描述信息
}
