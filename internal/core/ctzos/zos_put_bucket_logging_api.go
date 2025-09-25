package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketLoggingApi
/* 为指定桶设置日志转存配置。
 */type ZosPutBucketLoggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketLoggingApi(client *core.CtyunClient) *ZosPutBucketLoggingApi {
	return &ZosPutBucketLoggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-logging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketLoggingApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketLoggingRequest) (*ZosPutBucketLoggingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosPutBucketLoggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketLoggingRequest struct {
	Bucket              string                                         `json:"bucket,omitempty"`    /*  桶名  */
	RegionID            string                                         `json:"regionID,omitempty"`  /*  区域 ID  */
	BucketLoggingStatus *ZosPutBucketLoggingBucketLoggingStatusRequest `json:"bucketLoggingStatus"` /*  日志转存配置  */
}

type ZosPutBucketLoggingBucketLoggingStatusRequest struct {
	LoggingEnabled *ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledRequest `json:"loggingEnabled"` /*  日志转存配置  */
}

type ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledRequest struct {
	TargetPrefix string                                                                     `json:"targetPrefix"`           /*  所有日志对象键的前缀。如果你在一个桶中存储来自多个桶的日志文件，你可以使用前缀来区分哪些日志文件来自哪个桶。  */
	TargetBucket string                                                                     `json:"targetBucket,omitempty"` /*  指定希望 ZOS 存储服务器访问日志的桶。你可以让你的日志传递到你拥有的任何桶，包括被记录的同一个桶。你也可以配置多个桶，将它们的日志传递到同一个目标桶。在这种情况下，你应该为每个源桶选择一个不同的TargetPrefix，以便交付的日志文件可以通过密钥区分  */
	TargetGrants []*ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsRequest `json:"targetGrants"`           /*  授权信息  */
}

type ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsRequest struct {
	Grantee    *ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsGranteeRequest `json:"grantee"`              /*  被授权者信息  */
	Permission string                                                                          `json:"permission,omitempty"` /*  分配给桶的被授权者的日志记录权限。 支持 FULL_CONTROL，READ，WRITE  */
}

type ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsGranteeRequest struct {
	RawType      string `json:"type,omitempty"`         /*  被授权者类型， CanonicalUser，AmazonCustomerByEmail 两者之一。 type 为 CanonicalUser 时，必填 ID；为 AmazonCustomerByEmail，必填 emailAddress  */
	EmailAddress string `json:"emailAddress,omitempty"` /*  被授权者的邮箱  */
	DisplayName  string `json:"displayName,omitempty"`  /*  被授权者的显示名  */
	ID           string `json:"ID,omitempty"`           /*  被授权者的 ID  */
}

type ZosPutBucketLoggingResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
