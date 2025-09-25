package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosDeleteBucketTaggingApi
/* 删除指定桶上的标签。
 */type ZosDeleteBucketTaggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosDeleteBucketTaggingApi(client *core.CtyunClient) *ZosDeleteBucketTaggingApi {
	return &ZosDeleteBucketTaggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/delete-bucket-tagging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosDeleteBucketTaggingApi) Do(ctx context.Context, credential core.Credential, req *ZosDeleteBucketTaggingRequest) (*ZosDeleteBucketTaggingResponse, error) {
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
	var resp ZosDeleteBucketTaggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosDeleteBucketTaggingRequest struct {
	Bucket   string `json:"bucket,omitempty"`   /*  桶名  */
	RegionID string `json:"regionID,omitempty"` /*  区域 ID  */
}

type ZosDeleteBucketTaggingResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
