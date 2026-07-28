package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbDescribeIpWhitelistApi 查询白名单列表信息
type MongodbDescribeIpWhitelistApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDescribeIpWhitelistApi(client *ctyunsdk.CtyunClient) *MongodbDescribeIpWhitelistApi {
	return &MongodbDescribeIpWhitelistApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v2/openApi/queryWhiteList",
		},
	}
}

func (this *MongodbDescribeIpWhitelistApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDescribeIpWhitelistRequest, headers *MongodbDescribeIpWhitelistRequestHeaders) (response *MongodbDescribeIpWhitelistResponse, err error) {
	builder := this.WithCredential(&credential)

	if headers.RegionID == "" {
		err = errors.New("regionId is required")
		return nil, err
	}
	builder.AddHeader("regionId", headers.RegionID)

	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	if req.ProdInstId == "" {
		err = errors.New("prodInstId is required")
		return nil, err
	}
	builder.AddParam("prodInstId", req.ProdInstId)

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

type MongodbDescribeIpWhitelistRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID
}

type MongodbDescribeIpWhitelistRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDescribeIpWhitelistResponse struct {
	StatusCode int32                                   `json:"statusCode"` // 接口状态码
	Message    *string                                 `json:"message"`    // 返回消息
	Error      string                                  `json:"error"`      // 错误码
	ReturnObj  *MongodbDescribeIpWhitelistResponseData `json:"returnObj"`  // 返回对象
}

type MongodbDescribeIpWhitelistResponseData struct {
	WhitelistGroup []MongodbWhitelistGroup `json:"list"` // 白名单分组列表
}

type MongodbWhitelistGroup struct {
	GroupName     string `json:"groupName"`     // 白名单分组名称
	WhiteListType int32  `json:"whiteListType"` // 白名单分组名称
	Id            int32  `json:"id"`            // 白名单分组名称
	IpList        string `json:"ipList"`        // IP列表
	IpType        string `json:"ipType"`        // IP列表
}
