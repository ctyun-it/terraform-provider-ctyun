package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketLoggingApi
/* 查询指定桶的日志转存配置。
 */type ZosGetBucketLoggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketLoggingApi(client *core.CtyunClient) *ZosGetBucketLoggingApi {
	return &ZosGetBucketLoggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-logging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketLoggingApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketLoggingRequest) (*ZosGetBucketLoggingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketLoggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketLoggingRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketLoggingResponse struct {
	StatusCode  int64                                 `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string                                `json:"message,omitempty"`     /*  状态描述  */
	Description string                                `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketLoggingReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketLoggingReturnObjResponse struct {
	LoggingEnabled *ZosGetBucketLoggingReturnObjLoggingEnabledResponse `json:"loggingEnabled"` /*  描述日志的存储位置以及分配给桶的所有日志对象键的前缀。  */
}

type ZosGetBucketLoggingReturnObjLoggingEnabledResponse struct {
	TargetPrefix string                                                            `json:"targetPrefix,omitempty"` /*  所有日志对象键的前缀。如果一个桶中需存储来自多个桶的日志文件，可以使用前缀来区分哪些日志文件来自哪个桶。  */
	TargetBucket string                                                            `json:"targetBucket,omitempty"` /*  指定希望 OSS 存储服务器访问日志的桶。可以是本桶，也可以配置多个桶，将多个桶的日志存储到同一个目标桶。在这种情况下，可以为每个源桶选择一个不同的targetPrefix，以便交付的日志文件可以区分  */
	TargetGrants []*ZosGetBucketLoggingReturnObjLoggingEnabledTargetGrantsResponse `json:"targetGrants"`           /*  授权信息的容器。对对象所有权使用桶所有者强制设置桶不支持目标授予  */
}

type ZosGetBucketLoggingReturnObjLoggingEnabledTargetGrantsResponse struct {
	Grantee    *ZosGetBucketLoggingReturnObjLoggingEnabledTargetGrantsGranteeResponse `json:"grantee"`              /*  被授予权限的人的容器  */
	Permission string                                                                 `json:"permission,omitempty"` /*  分配给桶的被授权者的日志记录权限  */
}

type ZosGetBucketLoggingReturnObjLoggingEnabledTargetGrantsGranteeResponse struct {
	RawType      string `json:"type,omitempty"`         /*  被授权者类型， CanonicalUser，AmazonCustomerByEmail 二者之一  */
	EmailAddress string `json:"emailAddress,omitempty"` /*  被授权者的邮箱  */
	DisplayName  string `json:"displayName,omitempty"`  /*  被授权者的显示名  */
	ID           string `json:"ID,omitempty"`           /*  被授权者的 ID  */
}
