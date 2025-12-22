package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbUpdateIpWhitelistApi 更新白名单分组
type MongodbUpdateIpWhitelistApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdateIpWhitelistApi(client *ctyunsdk.CtyunClient) *MongodbUpdateIpWhitelistApi {
	return &MongodbUpdateIpWhitelistApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/updateWhiteList",
		},
	}
}

func (this *MongodbUpdateIpWhitelistApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdateIpWhitelistRequest, headers *MongodbUpdateIpWhitelistRequestHeaders) (response *MongodbUpdateIpWhitelistResponse, err error) {
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

type MongodbUpdateIpWhitelistRequest struct {
	ProdInstId    string `json:"prodInstId"`    // 实例ID
	WhiteListId   string `json:"id"`            // 白名单分组序号
	IpType        string `json:"ipType"`        // ip类型，取值ipv4或者ipv6
	WhiteListType string `json:"whiteListType"` // 修改后的允许访问的白名单列表，JSON数组字符串，示例："[\"10.138.16.9\"]"
	GroupName     string `json:"groupName"`     // 白名单分组名称
	IpList        string `json:"ipList"`        // 修改后的允许访问的白名单列表，JSON数组字符串，示例："[\"10.138.16.9\"]"
}

type MongodbUpdateIpWhitelistRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbUpdateIpWhitelistResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
