package mongodb

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbUnbindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUnbindEipApi(client *ctyunsdk.CtyunClient) *MongodbUnbindEipApi {
	return &MongodbUnbindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/unbind",
		},
	}
}

func (this *MongodbUnbindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUnbindEipRequest, header *MongodbUnbindEipRequestHeader) (unbindResponse *MongodbUnbindEipResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return
	}
	if header.ProjectID != nil {
		builder.AddHeader("project-id", *header.ProjectID)
	}
	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return
	}
	unbindResponse = &MongodbUnbindEipResponse{}
	err = resp.Parse(unbindResponse)
	if err != nil {
		return
	}
	return unbindResponse, nil
}

type MongodbUnbindEipRequest struct {
	EipID  string `json:"eipId"`
	Eip    string `json:"eip"`
	InstID string `json:"instId"`
}

type MongodbUnbindEipRequestHeader struct {
	ProjectID *string `json:"project_id,omitempty"`
}

type MongodbUnbindEipResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}
