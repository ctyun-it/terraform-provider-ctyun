package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbDeleteIpWhitelistApi 删除白名单分组
type MongodbDeleteIpWhitelistApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDeleteIpWhitelistApi(client *ctyunsdk.CtyunClient) *MongodbDeleteIpWhitelistApi {
	return &MongodbDeleteIpWhitelistApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/deleteWhiteList",
		},
	}
}

func (this *MongodbDeleteIpWhitelistApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDeleteIpWhitelistRequest, headers *MongodbDeleteIpWhitelistRequestHeaders) (response *MongodbDeleteIpWhitelistResponse, err error) {
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

type MongodbDeleteIpWhitelistRequest struct {
	ProdInstId  string `json:"prodInstId"`  // 实例ID
	WhiteListId string `json:"whiteListId"` // 白名单分组名称
}

type MongodbDeleteIpWhitelistRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDeleteIpWhitelistResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
