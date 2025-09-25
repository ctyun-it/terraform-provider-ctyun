package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketLifecycleConfApi
/* 查询指定桶的生命周期配置。
 */type ZosGetBucketLifecycleConfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketLifecycleConfApi(client *core.CtyunClient) *ZosGetBucketLifecycleConfApi {
	return &ZosGetBucketLifecycleConfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-lifecycle-conf",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketLifecycleConfApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketLifecycleConfRequest) (*ZosGetBucketLifecycleConfResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketLifecycleConfResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketLifecycleConfRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketLifecycleConfResponse struct {
	StatusCode  int64                                       `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                      `json:"message,omitempty"`     /*  状态描述  */
	Description string                                      `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketLifecycleConfReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                      `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                      `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketLifecycleConfReturnObjResponse struct {
	Rules []*ZosGetBucketLifecycleConfReturnObjRulesResponse `json:"rules"` /*  规则  */
}

type ZosGetBucketLifecycleConfReturnObjRulesResponse struct {
	ID                             string                                                                         `json:"ID,omitempty"`                   /*  规则ID  */
	Expiration                     *ZosGetBucketLifecycleConfReturnObjRulesExpirationResponse                     `json:"expiration"`                     /*  用日期或天数指定对象的过期时间  */
	Status                         string                                                                         `json:"status,omitempty"`               /*  规则是否启用，值为 Enabled 或 Disabled  */
	NoncurrentVersionExpiration    *ZosGetBucketLifecycleConfReturnObjRulesNoncurrentVersionExpirationResponse    `json:"noncurrentVersionExpiration"`    /*  标识历史版本的过期规则  */
	NoncurrentVersionTransitions   []*ZosGetBucketLifecycleConfReturnObjRulesNoncurrentVersionTransitionsResponse `json:"noncurrentVersionTransitions"`   /*  标识历史版本的转存储规则  */
	Filter                         *ZosGetBucketLifecycleConfReturnObjRulesFilterResponse                         `json:"filter"`                         /*  过滤应用规则的对象  */
	Prefix                         string                                                                         `json:"prefix,omitempty"`               /*  识别规则所适用的一个或多个对象的前缀  */
	AbortIncompleteMultipartUpload *ZosGetBucketLifecycleConfReturnObjRulesAbortIncompleteMultipartUploadResponse `json:"abortIncompleteMultipartUpload"` /*  指定自不完整的多部分上传开始后，在自动永久删除上传的所有部分之前将等待的天数  */
	Transitions                    []*ZosGetBucketLifecycleConfReturnObjRulesTransitionsResponse                  `json:"transitions"`                    /*  指定桶内对象何时过渡到指定的存储类别。  */
}

type ZosGetBucketLifecycleConfReturnObjRulesExpirationResponse struct {
	Date                      string `json:"date,omitempty"`            /*  ISO-8601 格式的日期字符串，精确到天。表示对象在什么日期被移动或删除。  */
	ExpiredObjectDeleteMarker *bool  `json:"expiredObjectDeleteMarker"` /*  指定是否自动移除过期删除标记  */
	Days                      int64  `json:"days,omitempty"`            /*  表示受该规则约束的对象的寿命，以天为单位。该值必须是一个非零的正整数  */
}

type ZosGetBucketLifecycleConfReturnObjRulesNoncurrentVersionExpirationResponse struct {
	NoncurrentDays int64 `json:"noncurrentDays,omitempty"` /*  指定对象在 OSS 可以执行关联操作之前处于非当前状态的天数。该值必须是一个非零的正整数  */
}

type ZosGetBucketLifecycleConfReturnObjRulesNoncurrentVersionTransitionsResponse struct {
	NoncurrentDays int64  `json:"noncurrentDays,omitempty"` /*  指定对象在 ZOS 可以执行关联操作之前处于非当前状态的天数  */
	StorageClass   string `json:"storageClass,omitempty"`   /*  用于存储对象的存储类  */
}

type ZosGetBucketLifecycleConfReturnObjRulesFilterResponse struct {
	And    *ZosGetBucketLifecycleConfReturnObjRulesFilterAndResponse `json:"and"`              /*  这在生命周期规则过滤器中用于将逻辑 AND 应用于两个或多个谓词。生命周期规则将应用于与 And 运算符中配置的所有谓词匹配的任何对象  */
	Prefix string                                                    `json:"prefix,omitempty"` /*  标识规则适用的一个或多个对象的前缀  */
	Tag    *ZosGetBucketLifecycleConfReturnObjRulesFilterTagResponse `json:"tag"`              /*  这个标签必须存在于对象的标签集中，以便规则的应用  */
}

type ZosGetBucketLifecycleConfReturnObjRulesAbortIncompleteMultipartUploadResponse struct {
	DaysAfterInitiation int64 `json:"daysAfterInitiation,omitempty"` /*  指定 OSS 中止未完成分段上传的天数。  */
}

type ZosGetBucketLifecycleConfReturnObjRulesTransitionsResponse struct {
	Date         string `json:"date,omitempty"`         /*  指示对象何时转换到指定的存储类。日期值为 ISO 8601 格式，精确到天。  */
	Days         int64  `json:"days,omitempty"`         /*  指示对象在创建后转换到指定存储类的天数。该值必须是正整数  */
	StorageClass string `json:"storageClass,omitempty"` /*  该对象过渡到的存储类  */
}

type ZosGetBucketLifecycleConfReturnObjRulesFilterAndResponse struct {
	Prefix string                                                          `json:"prefix,omitempty"` /*  标识规则适用的一个或多个对象的前缀  */
	Tags   []*ZosGetBucketLifecycleConfReturnObjRulesFilterAndTagsResponse `json:"tags"`             /*  所有这些标签都必须存在于对象的标签集中，才能应用规则  */
}

type ZosGetBucketLifecycleConfReturnObjRulesFilterTagResponse struct {
	Key   string `json:"key,omitempty"`   /*  标签名称  */
	Value string `json:"value,omitempty"` /*  标签值  */
}

type ZosGetBucketLifecycleConfReturnObjRulesFilterAndTagsResponse struct {
	Key   string `json:"key,omitempty"`   /*  标签名称  */
	Value string `json:"value,omitempty"` /*  标签值  */
}
