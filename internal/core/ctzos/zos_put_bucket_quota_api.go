package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketQuotaApi
/* 对桶配额进行修改。
 */type ZosPutBucketQuotaApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketQuotaApi(client *core.CtyunClient) *ZosPutBucketQuotaApi {
	return &ZosPutBucketQuotaApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-quota",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketQuotaApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketQuotaRequest) (*ZosPutBucketQuotaResponse, error) {
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
	var resp ZosPutBucketQuotaResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketQuotaRequest struct {
	Bucket     string `json:"bucket,omitempty"`     /*  桶名  */
	RegionID   string `json:"regionID,omitempty"`   /*  区域 ID  */
	Enabled    *bool  `json:"enabled"`              /*  是否开启配额限制，默认值为false，值为true时maxSizeKb和maxObjects至少一个大于等于0；不传该字段时将关闭配额限制  */
	MaxSizeKb  int64  `json:"maxSizeKb,omitempty"`  /*  最大的size容量(单位KB)，传入小于0或不传值均为无限制  */
	MaxObjects int64  `json:"maxObjects,omitempty"` /*  最大的objects数量，传入小于0或不传值均为无限制  */
}

type ZosPutBucketQuotaResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
