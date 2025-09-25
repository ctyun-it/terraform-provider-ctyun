package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCreateBucketApi
/* 在指定账号下创建一个新桶。
 */type ZosCreateBucketApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCreateBucketApi(client *core.CtyunClient) *ZosCreateBucketApi {
	return &ZosCreateBucketApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/create-bucket",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCreateBucketApi) Do(ctx context.Context, credential core.Credential, req *ZosCreateBucketRequest) (*ZosCreateBucketResponse, error) {
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
	var resp ZosCreateBucketResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCreateBucketRequest struct {
	RegionID        string                                 `json:"regionID,omitempty"`    /*  区域 ID  */
	ACL             string                                 `json:"ACL,omitempty"`         /*  桶权限，可选值为 `private`、`public-read`、`public-read-write`，分别表示私有、公共读、公共读写，默认为 `private`  */
	Bucket          string                                 `json:"bucket,omitempty"`      /*  桶名称，不可为空。长度 3-63 个字符内（含）字符只能有大小写字母、数字以及英文句号 （.） 和中划线 （-）。禁止两个英文句号（.）或英文句号（.）中划线（-）相邻。禁止英文句号（.）和中划线（-）作为开头或结尾。禁止使用 ip 地址作为桶名  */
	ProjectID       string                                 `json:"projectID,omitempty"`   /*  企业项目ID，默认将使用 default 企业项目，其 ID 为 "0"  */
	CmkUUID         string                                 `json:"cmkUUID,omitempty"`     /*  密钥管理服务中创建的密钥 ID，使用此参数时，isEncrypted 必须为 true。当 isEncrypted 为 true 但未指定此参数时，会自动创建密钥  */
	IsEncrypted     *bool                                  `json:"isEncrypted"`           /*  加密状态  */
	StorageType     string                                 `json:"storageType,omitempty"` /*  存储类型，可选的值为 `STANDARD`、`STANDARD_IA`、`GLACIER`，分别表示标准、低频、归档，默认为 `STANDARD`  */
	AZPolicy        string                                 `json:"AZPolicy,omitempty"`    /*  az 策略，可选值为 `single-az`、`multi-az`，分别表示单 az、多 az，默认为 `single-az`  */
	Labels          []*ZosCreateBucketLabelsRequest        `json:"labels"`                /*  桶标签  */
	OtherBucketInfo *ZosCreateBucketOtherBucketInfoRequest `json:"otherBucketInfo"`       /*  其他创建桶信息，如启用对象锁定 {"ObjectLockEnabledForBucket": true}  */
}

type ZosCreateBucketLabelsRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签名  */
	Value string `json:"value,omitempty"` /*  值  */
}

type ZosCreateBucketOtherBucketInfoRequest struct {
	ObjectLockEnabledForBucket *bool `json:"ObjectLockEnabledForBucket"` /*  是否开启对象锁定  */
}

type ZosCreateBucketResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
