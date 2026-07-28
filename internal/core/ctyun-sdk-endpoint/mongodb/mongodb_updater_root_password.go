package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbUpdatePasswordApi 手动备份实例
type MongodbUpdatePasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdatePasswordApi(client *ctyunsdk.CtyunClient) *MongodbUpdatePasswordApi {
	return &MongodbUpdatePasswordApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/updatePassword",
		},
	}
}

func (this *MongodbUpdatePasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdatePasswordRequest, headers *MongodbUpdatePasswordRequestHeaders) (response *MongodbUpdatePasswordResponse, err error) {
	builder := this.WithCredential(&credential)
	_, err = builder.WriteJson(req)
	if err != nil {
		return nil, err
	}

	if headers.RegionID == "" {
		err = errors.New("regionId is required")
		return nil, err
	}
	builder.AddHeader("regionId", headers.RegionID)

	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	resp, err := this.client.RequestToEndpoint(ctx, EndpointNameMongodb, builder)
	if err != nil {
		return nil, err
	}

	err = resp.Parse(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type MongodbUpdatePasswordRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID
	Password   string `json:"nPassword"`  // nPassword
}

type MongodbUpdatePasswordRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbUpdatePasswordResponse struct {
	StatusCode int32                                  `json:"statusCode"` // 接口状态码
	Message    *string                                `json:"message"`    // 返回消息
	ReturnObj  MongodbUpdatePasswordResponseReturnObj `json:"returnObj"`  // 返回对象
}
type MongodbUpdatePasswordResponseReturnObj struct {
}
