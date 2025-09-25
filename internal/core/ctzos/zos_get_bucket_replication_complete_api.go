package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetBucketReplicationCompleteApi
/* 获取某个跨域复制规则是否完成历史数据的同步。
 */type ZosGetBucketReplicationCompleteApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetBucketReplicationCompleteApi(client *core.CtyunClient) *ZosGetBucketReplicationCompleteApi {
	return &ZosGetBucketReplicationCompleteApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-bucket-replication-complete",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetBucketReplicationCompleteApi) Do(ctx context.Context, credential core.Credential, req *ZosGetBucketReplicationCompleteRequest) (*ZosGetBucketReplicationCompleteResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("fuid", req.Fuid)
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetBucketReplicationCompleteResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetBucketReplicationCompleteRequest struct {
	Bucket   string /*  桶名  */
	Fuid     string /*  同步规则ID  */
	RegionID string /*  区域ID  */
}

type ZosGetBucketReplicationCompleteResponse struct {
	StatusCode  int64                                               `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                              `json:"message,omitempty"`     /*  状态描述  */
	Description string                                              `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   []*ZosGetBucketReplicationCompleteReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                              `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                              `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetBucketReplicationCompleteReturnObjResponse struct {
	Fuid     string `json:"fuid,omitempty"`     /*  同步规则ID  */
	Complete int64  `json:"complete,omitempty"` /*  同步是否完成，1 完成，0 未完成  */
}
