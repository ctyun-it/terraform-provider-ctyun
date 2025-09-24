package mongodb

import (
	"context"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

type MongodbBindEipApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbBindEipApi(client *ctyunsdk.CtyunClient) *MongodbBindEipApi {
	return &MongodbBindEipApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/teledb-dcp/v2/openapi/dcp-order-info/eips/bind",
		},
	}
}

func (this *MongodbBindEipApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbBindEipRequest, header *MongodbBindEipRequestHeader) (bindResponse *MongodbBindEipResponse, err error) {
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
	bindResponse = &MongodbBindEipResponse{}
	err = resp.Parse(bindResponse)
	if err != nil {
		return
	}
	return bindResponse, nil
}

type MongodbBindEipRequest struct {
	EipID  string `json:"eipId"`
	Eip    string `json:"eip"`
	InstID string `json:"instId"`
	HostIp string `json:"hostIp"` //主机ip
}

type MongodbBindEipRequestHeader struct {
	ProjectID *string `json:"project-id,omitempty"`
}

type MongodbBindEipResponse struct {
	StatusCode int32  `json:"statusCode"` // 接口状态码
	Error      string `json:"error"`      // 错误码，失败时返回，成功时为空
	Message    string `json:"message"`    // 描述信息
}
