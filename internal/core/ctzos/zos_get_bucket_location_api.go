package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketLocationApi
/* 查询桶所在的区域。
 */type ZosGetBucketLocationApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketLocationApi(client *core.CtyunClient) *ZosGetBucketLocationApi {
	return &ZosGetBucketLocationApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-location",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketLocationApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketLocationRequest) (*ZosGetBucketLocationResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketLocationResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketLocationRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketLocationResponse struct {
	StatusCode  int64                                  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                 `json:"message,omitempty"`     /*  状态描述  */
	Description string                                 `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetBucketLocationReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                 `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                 `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketLocationReturnObjResponse struct {
	LocationConstraint string `json:"locationConstraint,omitempty"` /*  所在区域  */
}
