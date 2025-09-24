package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketReplicationRegionApi
/* 获取可复制到的目标桶所在的地域。
 */type ZosGetBucketReplicationRegionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketReplicationRegionApi(client *core.CtyunClient) *ZosGetBucketReplicationRegionApi {
	return &ZosGetBucketReplicationRegionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-replication-region",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketReplicationRegionApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketReplicationRegionRequest) (*ZosGetBucketReplicationRegionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketReplicationRegionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketReplicationRegionRequest struct {
	Bucket   string /*  桶名  */
	RegionID string /*  区域 ID  */
}

type ZosGetBucketReplicationRegionResponse struct {
	StatusCode  int64                                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                            `json:"message,omitempty"`     /*  状态描述  */
	Description string                                            `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   []*ZosGetBucketReplicationRegionReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                            `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                            `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketReplicationRegionReturnObjResponse struct {
	RegionID   string   `json:"regionID,omitempty"`   /*  区域ID  */
	RegionIDv2 string   `json:"regionIDv2,omitempty"` /*  新版区域ID  */
	RegionName string   `json:"regionName,omitempty"` /*  区域名  */
	Num        int64    `json:"num,omitempty"`        /*  数量  */
	Buckets    []string `json:"buckets"`              /*  桶名的数组  */
}
