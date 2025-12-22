package mongodb

import (
	"context"
	"errors"
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
)

// MongodbUpdateAccountPasswordApi 更新数据库账号密码
type MongodbUpdateAccountPasswordApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdateAccountPasswordApi(client *ctyunsdk.CtyunClient) *MongodbUpdateAccountPasswordApi {
	return &MongodbUpdateAccountPasswordApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/updateAccountPassword",
		},
	}
}

func (this *MongodbUpdateAccountPasswordApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdateAccountPasswordRequest, headers *MongodbUpdateAccountPasswordRequestHeaders) (response *MongodbUpdateAccountPasswordResponse, err error) {
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
	if headers.ProdInstId == "" {
		err = errors.New("prodInstId is required")
		return nil, err
	}
	builder.AddHeader("prodInstId", headers.ProdInstId)
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

type MongodbUpdateAccountPasswordRequest struct {
	AccountName     string `json:"accountName"` // 账号名称
	AccountPassword string `json:"nPassword"`   // 新密码
	Database        string `json:"authDb"`      // 新密码
}

type MongodbUpdateAccountPasswordRequestHeaders struct {
	ProjectID  *string `json:"projectId,omitempty"` // 项目ID
	RegionID   string  `json:"regionId"`            // 资源池regionId
	ProdInstId string  `json:"prodInstId"`          // 实例ID
}

type MongodbUpdateAccountPasswordResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}
