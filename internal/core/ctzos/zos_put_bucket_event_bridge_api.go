package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketEventBridgeApi
/* 为指定的桶设置事件总线功能。
 */type ZosPutBucketEventBridgeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketEventBridgeApi(client *core.CtyunClient) *ZosPutBucketEventBridgeApi {
	return &ZosPutBucketEventBridgeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-event-bridge",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketEventBridgeApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketEventBridgeRequest) (*ZosPutBucketEventBridgeResponse, error) {
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
	var resp ZosPutBucketEventBridgeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketEventBridgeRequest struct {
	Bucket                   string `json:"bucket,omitempty"`         /*  桶名。  */
	RegionID                 string `json:"regionID,omitempty"`       /*  区域 ID  */
	BucketEventBridgeEnabled bool   `json:"bucketEventBridgeEnabled"` /*  桶的事件总线配置，true为启用桶的事件总线，false为关闭桶的事件总线  */
}

type ZosPutBucketEventBridgeResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
