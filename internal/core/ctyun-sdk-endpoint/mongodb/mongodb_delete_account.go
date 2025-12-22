package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbDeleteAccountApi 删除数据库账号
type MongodbDeleteAccountApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDeleteAccountApi(client *ctyunsdk.CtyunClient) *MongodbDeleteAccountApi {
	return &MongodbDeleteAccountApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/deleteAccount",
		},
	}
}

func (this *MongodbDeleteAccountApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDeleteAccountRequest, headers *MongodbDeleteAccountRequestHeaders) (response *MongodbDeleteAccountResponse, err error) {
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
	if headers.ProdInstId == "" {
		err = errors.New("prodInstId is required")
		return nil, err
	}
	builder.AddHeader("prodInstId", headers.ProdInstId)

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

type MongodbDeleteAccountRequest struct {
	AccountName  string `json:"user"` // 账号名称
	DatabaseName string `json:"db"`   // 账号名称
}

type MongodbDeleteAccountRequestHeaders struct {
	ProdInstId string  `json:"prodInstId"`          // 实例ID
	ProjectID  *string `json:"projectId,omitempty"` // 项目ID
	RegionID   string  `json:"regionId"`            // 资源池regionId
}

type MongodbDeleteAccountResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
