package mongodb

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// MongodbCreateBackupApi 手动备份实例
type MongodbCreateBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbCreateBackupApi(client *ctyunsdk.CtyunClient) *MongodbCreateBackupApi {
	return &MongodbCreateBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/createBackup",
		},
	}
}

func (this *MongodbCreateBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbCreateBackupRequest, headers *MongodbCreateBackupRequestHeaders) (response *MongodbCreateBackupResponse, err error) {
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

type MongodbCreateBackupRequest struct {
	ProdInstId   string  `json:"prodInstId"`             // 实例ID
	BackupMethod *string `json:"backupMethod,omitempty"` // 备份方式
	BackupType   *string `json:"backupType,omitempty"`   // 备份类型
	BackupName   *string `json:"backupName,omitempty"`   // 备份名称
	Description  *string `json:"description,omitempty"`  // 备份描述
}

type MongodbCreateBackupRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbCreateBackupResponse struct {
	StatusCode int32                                `json:"statusCode"` // 接口状态码
	Message    *string                              `json:"message"`    // 返回消息
	Error      string                               `json:"error"`      // 错误码
	ReturnObj  MongodbCreateBackupResponseReturnObj `json:"returnObj"`  // 返回对象
}
type MongodbCreateBackupResponseReturnObj struct {
	Code string `json:"code"` // 返回对象
}

// MongodbDeleteBackupApi 删除实例的单个备份
type MongodbDeleteBackupApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDeleteBackupApi(client *ctyunsdk.CtyunClient) *MongodbDeleteBackupApi {
	return &MongodbDeleteBackupApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/deleteBackupFileJob",
		},
	}
}

func (this *MongodbDeleteBackupApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDeleteBackupRequest, headers *MongodbDeleteBackupRequestHeaders) (response *MongodbDeleteBackupResponse, err error) {
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

type MongodbDeleteBackupRequest struct {
	ProdInstId string `json:"prodInstId"` // 实例ID
	BackupId   string `json:"backupId"`   // 备份ID
}

type MongodbDeleteBackupRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDeleteBackupResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}

// MongodbDescribeBackupsApi 查询实例备份列表
type MongodbDescribeBackupsApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDescribeBackupsApi(client *ctyunsdk.CtyunClient) *MongodbDescribeBackupsApi {
	return &MongodbDescribeBackupsApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS2/v2/openApi/describeBackups",
		},
	}
}

func (this *MongodbDescribeBackupsApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDescribeBackupsRequest, headers *MongodbDescribeBackupsRequestHeaders) (response *MongodbDescribeBackupsResponse, err error) {
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

	//if req.BackupId != nil {
	//	builder.AddParam("backupId", *req.BackupId)
	//}
	if req.BackupType != nil {
		builder.AddParam("backupType", *req.BackupType)
	}
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

type MongodbDescribeBackupsRequest struct {
	ProdInstId string  `json:"prodInstId"`           // 实例ID
	BackupName *string `json:"backupName,omitempty"` // 备份ID
	BackupType *string `json:"backupType,omitempty"` // 备份ID
	PageNow    int32   `json:"pageNow"`              // 当前页码
	PageSize   int32   `json:"pageSize"`             // 每页记录数
}

type MongodbDescribeBackupsRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDescribeBackupsResponse struct {
	StatusCode int32                            `json:"statusCode"` // 接口状态码
	Message    *string                          `json:"message"`    // 返回消息
	Error      string                           `json:"error"`      // 错误码
	ReturnObj  *MongodbDescribeBackupsReturnObj `json:"returnObj"`  // 返回对象
}

type MongodbDescribeBackupsReturnObj struct {
	Total    int32               `json:"total"`    // 总记录数
	Pages    int32               `json:"pages"`    // 总页数
	PageSize int32               `json:"pageSize"` // 每页记录数
	PageNow  int32               `json:"pageNow"`  // 当前页码
	List     []MongodbBackupInfo `json:"list"`     // 备份信息列表
}

type MongodbBackupInfo struct {
	BackupId          int32   `json:"id"`                // 备份ID
	BackupName        string  `json:"backupName"`        // 备份名称
	BackupMethod      string  `json:"backupMethod"`      // 备份方式
	BackupType        string  `json:"backupType"`        // 备份类型
	BackupStatus      string  `json:"backupStatus"`      // 备份状态
	BackupStartTime   string  `json:"backupStartTime"`   // 备份开始时间
	BackupEndTime     string  `json:"backupEndTime"`     // 备份结束时间
	BackupSize        int64   `json:"backupSize"`        // 备份大小
	Description       *string `json:"description"`       // 备份描述
	BackupTriggerType string  `json:"backupTriggerType"` // 备份触发类型
}
