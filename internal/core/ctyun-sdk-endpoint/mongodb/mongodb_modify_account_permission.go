package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbModifyAccountPermissionApi 修改数据库账号权限
type MongodbModifyAccountPermissionApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbModifyAccountPermissionApi(client *ctyunsdk.CtyunClient) *MongodbModifyAccountPermissionApi {
	return &MongodbModifyAccountPermissionApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/modifyAccountPermission",
		},
	}
}

func (this *MongodbModifyAccountPermissionApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbModifyAccountPermissionRequest, headers *MongodbModifyAccountPermissionRequestHeaders) (response *MongodbModifyAccountPermissionResponse, err error) {
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

	//var response MongodbModifyAccountPermissionResponse
	err = resp.Parse(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type MongodbModifyAccountPermissionRequest struct {
	AccountName  string                `json:"user"`            // 账号名称
	DatabaseName *string               `json:"db,omitempty"`    // 数据库名称
	Roles        *[]MongodbAccountRole `json:"roles,omitempty"` // 权限
}

type MongodbModifyAccountPermissionRequestHeaders struct {
	ProjectID  *string `json:"projectId,omitempty"` // 项目ID
	RegionID   string  `json:"regionId"`            // 资源池regionId
	ProdInstId string  `json:"prodInstId"`          // 实例ID
}

type MongodbModifyAccountPermissionResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
