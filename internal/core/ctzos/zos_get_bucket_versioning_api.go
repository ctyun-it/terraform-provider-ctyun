package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketVersioningApi
/* 查询桶是否开启了版本控制。
 */type ZosGetBucketVersioningApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketVersioningApi(client *core.CtyunClient) *ZosGetBucketVersioningApi {
	return &ZosGetBucketVersioningApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-versioning",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketVersioningApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketVersioningRequest) (*ZosGetBucketVersioningResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketVersioningResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketVersioningRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketVersioningResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回码<br>取值范围：800 成功  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
	ReturnObj   struct {
		Status string `json:"status"` // 桶的版本状态，值是 Enabled 或 Suspended
	} `json:"returnObj"`
}
