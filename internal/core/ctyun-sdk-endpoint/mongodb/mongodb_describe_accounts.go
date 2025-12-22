package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
)

// MongodbDescribeAccountsApi 查询实例的数据库账号信息
type MongodbDescribeAccountsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDescribeAccountsApi(client *ctyunsdk.CtyunClient) *MongodbDescribeAccountsApi {
	return &MongodbDescribeAccountsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v2/openApi/describeAccounts",
		},
	}
}

func (this *MongodbDescribeAccountsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDescribeAccountsRequest, headers *MongodbDescribeAccountsRequestHeaders) (response *MongodbDescribeAccountsResponse, err error) {
	builder := this.WithCredential(&credential)

	// 添加必需的headers
	if headers.RegionID == "" {
		err = errors.New("regionId is required")
		return nil, err
	}
	builder.AddHeader("regionId", headers.RegionID)

	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	// 添加必需的查询参数
	if req.ProdInstId == "" {
		err = errors.New("prodInstId is required")
		return nil, err
	}
	builder.AddParam("prodInstId", req.ProdInstId)

	// 添加分页参数
	if req.PageNow != 0 {
		builder.AddParam("pageNow", strconv.FormatInt(int64(req.PageNow), 10))
	}
	if req.PageSize != 0 {
		builder.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
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

type MongodbDescribeAccountsRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID
	PageNow    int32  `json:"pageNow"`    // 当前页码
	PageSize   int32  `json:"pageSize"`   // 每页记录数
}

type MongodbDescribeAccountsRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDescribeAccountsResponse struct {
	StatusCode int32                             `json:"statusCode"` // 接口状态码
	Message    *string                           `json:"message"`    // 返回消息
	Error      string                            `json:"error"`      // 错误码
	ReturnObj  *MongodbDescribeAccountsReturnObj `json:"returnObj"`  // 返回对象
}

type MongodbDescribeAccountsReturnObj struct {
	Total    int32                `json:"total"`    // 总记录数
	Pages    int32                `json:"pages"`    // 总页数
	PageSize int32                `json:"pageSize"` // 每页记录数
	List     []MongodbAccountInfo `json:"list"`     // 账号信息列表
	PageNow  int32                `json:"pageNow"`  // 当前页码
}

type MongodbAccountRole struct {
	Role string `json:"role"` // 角色
	DB   string `json:"db"`   // 数据库
}

type MongodbAccountInfo struct {
	Roles []MongodbAccountRole `json:"roles"` // 角色列表
	User  string               `json:"user"`  // 用户名
	DB    string               `json:"db"`    // 数据库
}
