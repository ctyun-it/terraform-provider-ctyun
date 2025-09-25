package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketReplicationApi
/* 跨域复制用于在不同ZOS数据中心（地域）之间自动异步复制存储桶中的对象。它可以近实时地将源存储桶中对象的创建、更新和删除等操作复制到目标存储桶，从而满足了跨地域容灾和用户数据复制的需求。通过调用该接口可创建跨域复制规则。
 */type ZosPutBucketReplicationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketReplicationApi(client *core.CtyunClient) *ZosPutBucketReplicationApi {
	return &ZosPutBucketReplicationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-replication",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketReplicationApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketReplicationRequest) (*ZosPutBucketReplicationResponse, error) {
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
	var resp ZosPutBucketReplicationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketReplicationRequest struct {
	Bucket         string   `json:"bucket,omitempty"`         /*  桶名  */
	RegionID       string   `json:"regionID,omitempty"`       /*  区域 ID  */
	TargetBucket   string   `json:"targetBucket,omitempty"`   /*  目标桶名  */
	TargetRegionID string   `json:"targetRegionID,omitempty"` /*  目标区域 ID  */
	Prefixes       []string `json:"prefixes"`                 /*  桶内对象前缀，默认为空数组，即同步所有。空数组：[]  */
	Plot           *bool    `json:"plot"`                     /*  同步策略, 同步时是否允许删除。默认为 false  */
	History        *bool    `json:"history"`                  /*  是否同步历史数据, 默认为 true  */
}

type ZosPutBucketReplicationResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
