package mongodb

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// MongodbCreateParamTemplateApi 创建参数组
type MongodbCreateParamTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbCreateParamTemplateApi(client *ctyunsdk.CtyunClient) *MongodbCreateParamTemplateApi {
	return &MongodbCreateParamTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS2/v2/openApi/createParameterGroup",
		},
	}
}

func (this *MongodbCreateParamTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbCreateParamTemplateRequest, headers *MongodbCreateParamTemplateRequestHeaders) (response *MongodbCreateParamTemplateResponse, err error) {
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

type MongodbCreateParamTemplateRequest struct {
	TemplateName  string `json:"parameterGroupName"` // 参数组名称
	TemplateDesc  string `json:"description"`        // 参数组描述
	EngineVersion string `json:"dbVersion"`          // 引擎版本
	NodeType      string `json:"nodeType"`           // 引擎类型
}

type MongodbCreateParamTemplateRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbCreateParamTemplateResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
	ReturnObj  *string `json:"returnObj"`  // 返回对象(参数组ID)
}

// MongodbDeleteParamTemplateApi 删除参数组
type MongodbDeleteParamTemplateApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDeleteParamTemplateApi(client *ctyunsdk.CtyunClient) *MongodbDeleteParamTemplateApi {
	return &MongodbDeleteParamTemplateApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS/v2/parameter/deleteParameterGroup",
		},
	}
}

func (this *MongodbDeleteParamTemplateApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDeleteParamTemplateRequest, headers *MongodbDeleteParamTemplateRequestHeaders) (response *MongodbDeleteParamTemplateResponse, err error) {
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

type MongodbDeleteParamTemplateRequest struct {
	ParameterGroupName string `json:"name"`           // 参数组ID
	Description        string `json:"description"`    // 参数组ID
	Engine             string `json:"engine_version"` // 参数组ID
}

type MongodbDeleteParamTemplateRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDeleteParamTemplateResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}

// MongodbUpdateParamTemplateDescApi 修改参数组描述
type MongodbUpdateParamTemplateDescApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbUpdateParamTemplateDescApi(client *ctyunsdk.CtyunClient) *MongodbUpdateParamTemplateDescApi {
	return &MongodbUpdateParamTemplateDescApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/DDS/v2/parameter/modifyParameterGroupDesc",
		},
	}
}

func (this *MongodbUpdateParamTemplateDescApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbUpdateParamTemplateDescRequest, headers *MongodbUpdateParamTemplateDescRequestHeaders) (response *MongodbUpdateParamTemplateDescResponse, err error) {
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

type MongodbUpdateParamTemplateDescRequest struct {
	TemplateId   string `json:"templateId"`   // 参数组ID
	TemplateDesc string `json:"templateDesc"` // 参数组描述
}

type MongodbUpdateParamTemplateDescRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbUpdateParamTemplateDescResponse struct {
	StatusCode int32   `json:"statusCode"` // 接口状态码
	Message    *string `json:"message"`    // 返回消息
	Error      string  `json:"error"`      // 错误码
}

// MongodbDescribeParamTemplatesApi 查询所有参数组
type MongodbDescribeParamTemplatesApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewMongodbDescribeParamTemplatesApi(client *ctyunsdk.CtyunClient) *MongodbDescribeParamTemplatesApi {
	return &MongodbDescribeParamTemplatesApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/DDS/v2/parameter/describeParameterGroups",
		},
	}
}

func (this *MongodbDescribeParamTemplatesApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *MongodbDescribeParamTemplatesRequest, headers *MongodbDescribeParamTemplatesRequestHeaders) (response *MongodbDescribeParamTemplatesResponse, err error) {
	builder := this.WithCredential(&credential)

	if headers.RegionID == "" {
		err = errors.New("regionId is required")
		return nil, err
	}
	builder.AddHeader("regionId", headers.RegionID)

	if headers.ProjectID != nil {
		builder.AddHeader("project-id", *headers.ProjectID)
	}

	if req.EngineType != nil {
		builder.AddParam("engineType", *req.EngineType)
	}
	if req.EngineVersion != nil {
		builder.AddParam("engineVersion", *req.EngineVersion)
	}
	if req.TemplateType != nil {
		builder.AddParam("templateType", *req.TemplateType)
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

type MongodbDescribeParamTemplatesRequest struct {
	EngineType    *string `json:"engineType,omitempty"`    // 引擎类型
	EngineVersion *string `json:"engineVersion,omitempty"` // 引擎版本
	TemplateType  *string `json:"templateType,omitempty"`  // 模板类型
	PageNow       int32   `json:"pageNow"`                 // 当前页码
	PageSize      int32   `json:"pageSize"`                // 每页记录数
}

type MongodbDescribeParamTemplatesRequestHeaders struct {
	ProjectID *string `json:"projectId,omitempty"` // 项目ID
	RegionID  string  `json:"regionId"`            // 资源池regionId
}

type MongodbDescribeParamTemplatesResponse struct {
	StatusCode int32                                   `json:"statusCode"` // 接口状态码
	Message    *string                                 `json:"message"`    // 返回消息
	Error      string                                  `json:"error"`      // 错误码
	ReturnObj  *MongodbDescribeParamTemplatesReturnObj `json:"returnObj"`  // 返回对象
}

type MongodbDescribeParamTemplatesReturnObj struct {
	Total    int32                      `json:"total"`    // 总记录数
	Pages    int32                      `json:"pages"`    // 总页数
	PageSize int32                      `json:"pageSize"` // 每页记录数
	PageNow  int32                      `json:"pageNow"`  // 当前页码
	List     []MongodbParamTemplateInfo `json:"list"`     // 参数组信息列表
}

type MongodbParamTemplateInfo struct {
	TemplateId       string  `json:"templateId"`                 // 参数组ID
	TemplateName     string  `json:"templateName"`               // 参数组名称
	TemplateDesc     string  `json:"templateDesc"`               // 参数组描述
	EngineType       string  `json:"engineType"`                 // 引擎类型
	EngineVersion    string  `json:"engineVersion"`              // 引擎版本
	TemplateType     string  `json:"templateType"`               // 模板类型
	SourceTemplateId *string `json:"sourceTemplateId,omitempty"` // 源参数组ID
	CreatedTime      string  `json:"createdTime"`                // 创建时间
	UpdatedTime      string  `json:"updatedTime"`                // 更新时间
}
