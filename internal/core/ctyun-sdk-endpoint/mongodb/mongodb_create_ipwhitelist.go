package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbCreateIpWhitelistApi 创建白名单分组
type MongodbCreateIpWhitelistApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbCreateIpWhitelistApi(client *ctyunsdk.CtyunClient) *MongodbCreateIpWhitelistApi {
	return &MongodbCreateIpWhitelistApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/addWhiteList",
		},
	}
}

func (this *MongodbCreateIpWhitelistApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbCreateIpWhitelistRequest, headers *MongodbCreateIpWhitelistRequestHeaders) (response *MongodbCreateIpWhitelistResponse, err error) {
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

type MongodbCreateIpWhitelistRequest struct {
	ProdInstId    string `json:"prodInstId"`    // 实例ID
	GroupName     string `json:"groupName"`     // 白名单分组名
	IpType        string `json:"ipType"`        // 白名单分组名
	WhiteListType string `json:"whiteListType"` // 白名单分组名
	IpList        string `json:"ipList"`        // IP列表
}

type MongodbCreateIpWhitelistRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbCreateIpWhitelistResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
