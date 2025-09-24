package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketLifecycleConfApi
/* 为指定桶设置生命周期配置。
 */type ZosPutBucketLifecycleConfApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketLifecycleConfApi(client *core.CtyunClient) *ZosPutBucketLifecycleConfApi {
	return &ZosPutBucketLifecycleConfApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-lifecycle-conf",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketLifecycleConfApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketLifecycleConfRequest) (*ZosPutBucketLifecycleConfResponse, error) {
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
	var resp ZosPutBucketLifecycleConfResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketLifecycleConfRequest struct {
	Bucket                 string                                                  `json:"bucket,omitempty"`       /*  桶名  */
	RegionID               string                                                  `json:"regionID,omitempty"`     /*  区域 ID  */
	LifecycleConfiguration *ZosPutBucketLifecycleConfLifecycleConfigurationRequest `json:"lifecycleConfiguration"` /*  生命周期配置  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRequest struct {
	Rules []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesRequest `json:"rules"` /*  生命周期规则  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesRequest struct {
	ID                             string                                                                                     `json:"ID,omitempty"`                   /*  规则ID  */
	Expiration                     *ZosPutBucketLifecycleConfLifecycleConfigurationRulesExpirationRequest                     `json:"expiration"`                     /*  用日期或天数指定对象的过期时间。若未传此参数，则参数 transitions、noncurrentVersionExpiration、noncurrentVersionTransitions 至少指定其中一个  */
	Status                         string                                                                                     `json:"status,omitempty"`               /*  规则是否启用，值为 Enabled 或 Disabled  */
	NoncurrentVersionExpiration    *ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionExpirationRequest    `json:"noncurrentVersionExpiration"`    /*  标识历史版本的过期规则。若未传此参数，则参数 expiration、transitions、noncurrentVersionTransitions 至少指定其中一个  */
	NoncurrentVersionTransitions   []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionTransitionsRequest `json:"noncurrentVersionTransitions"`   /*  标识历史版本的转存储规则。若未传此参数，则参数 expiration、transitions、noncurrentVersionExpiration 至少指定其中一个  */
	Filter                         *ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterRequest                         `json:"filter"`                         /*  过滤规则。可为空{}，若非空则至少应包含prefix、tag或and中的一个。若未传此参数，则需要指定 prefix 参数  */
	Prefix                         string                                                                                     `json:"prefix,omitempty"`               /*  识别规则所适用的一个或多个对象的前缀，若未传此参数，需要指定 filter 参数。若此参数与 filter 参数中的 prefix 同时存在，则 filter 中的 prefix 参数生效  */
	AbortIncompleteMultipartUpload *ZosPutBucketLifecycleConfLifecycleConfigurationRulesAbortIncompleteMultipartUploadRequest `json:"abortIncompleteMultipartUpload"` /*  指定自不完整的多部分上传开始后，在自动永久删除上传的所有部分之前将等待的天数  */
	Transitions                    []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesTransitionsRequest                  `json:"transitions"`                    /*  指定桶内对象何时过渡到指定的存储类别。若未传此参数，则参数 expiration、noncurrentVersionExpiration、noncurrentVersionTransitions 至少指定其中一个  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesExpirationRequest struct {
	Date                      string `json:"date,omitempty"`            /*  ISO-8601 格式的日期字符串，精确到天。表示对象在什么日期被移动或删除。且与参数 expiredObjectDeleteMarker 以及 days 不能共存，必须为 UTC 午夜0时   */
	ExpiredObjectDeleteMarker *bool  `json:"expiredObjectDeleteMarker"` /*  指定是否自动移除过期删除标记。如果设置为 true，删除标记将过期；如果设置为 false，则策略不执行任何操作。且与参数 date 以及 days 不能共存  */
	Days                      int64  `json:"days,omitempty"`            /*  表示受该规则约束的对象的寿命，以天为单位。该值必须是一个非零的正整数。且与参数 date 以及 expiredObjectDeleteMarker 不能共存  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionExpirationRequest struct {
	NoncurrentDays int64 `json:"noncurrentDays,omitempty"` /*  指定对象在 OSS 可以执行关联操作之前处于非当前状态的天数。该值必须是一个非零的正整数  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionTransitionsRequest struct {
	NoncurrentDays int64  `json:"noncurrentDays,omitempty"` /*  指定对象在 ZOS 可以执行关联操作之前处于非当前状态的天数  */
	StorageClass   string `json:"storageClass,omitempty"`   /*  用于存储对象的存储类。仅限于 GLACIER，STANDARD_IA  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterRequest struct {
	And    *ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndRequest `json:"and"`              /*  这在生命周期规则过滤器中用于将逻辑 AND 应用于两个或多个谓词。生命周期规则将应用于与 And 运算符中配置的所有谓词匹配的任何对象。若要同时使用 prefix, tag 参数，请使用此参数。  */
	Prefix string                                                                `json:"prefix,omitempty"` /*  标识规则适用的一个或多个对象的前缀。若要与 tag 参数一起使用，请使用 and 参数  */
	Tag    *ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterTagRequest `json:"tag"`              /*  这个标签必须存在于对象的标签集中，以便规则的应用。若要与 prefix 参数一起使用，请使用 and 参数  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesAbortIncompleteMultipartUploadRequest struct {
	DaysAfterInitiation int64 `json:"daysAfterInitiation,omitempty"` /*  指定 OSS 中止未完成分段上传的天数。  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesTransitionsRequest struct {
	Date         string `json:"date,omitempty"`         /*  指示对象何时转换到指定的存储类，与 days 不能共存。日期值为 ISO 8601 格式，必须为 UTC 午夜0时  */
	Days         int64  `json:"days,omitempty"`         /*  指示对象在创建后转换到指定存储类的天数，与 date 不能共存。该值必须是正整数  */
	StorageClass string `json:"storageClass,omitempty"` /*  该对象过渡到的存储类，可选值为 GLACIER，STANDARD_IA  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndRequest struct {
	Prefix string                                                                      `json:"prefix,omitempty"` /*  标识规则适用的一个或多个对象的前缀  */
	Tags   []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndTagsRequest `json:"tags"`             /*  所有这些标签都必须存在于对象的标签集中，才能应用规则  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterTagRequest struct {
	Key   string `json:"key,omitempty"`   /*  对象名称  */
	Value string `json:"value,omitempty"` /*  标签值  */
}

type ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndTagsRequest struct {
	Key   string `json:"key,omitempty"`   /*  对象名称  */
	Value string `json:"value,omitempty"` /*  标签值  */
}

type ZosPutBucketLifecycleConfResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
