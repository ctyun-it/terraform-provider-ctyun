package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutObjectTaggingApi
/* 对指定对象设置标签，通过标签对文件进行分类。
 */type ZosPutObjectTaggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutObjectTaggingApi(client *core.CtyunClient) *ZosPutObjectTaggingApi {
	return &ZosPutObjectTaggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-object-tagging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutObjectTaggingApi) Do(ctx context.Context, credential core.Credential, req *ZosPutObjectTaggingRequest) (*ZosPutObjectTaggingResponse, error) {
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
	var resp ZosPutObjectTaggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutObjectTaggingRequest struct {
	Bucket    string                             `json:"bucket,omitempty"`    /*  桶名  */
	Key       string                             `json:"key,omitempty"`       /*  对象名  */
	VersionID string                             `json:"versionID,omitempty"` /*  版本ID，在开启多版本时可使用  */
	RegionID  string                             `json:"regionID,omitempty"`  /*  区域 ID  */
	Tagging   *ZosPutObjectTaggingTaggingRequest `json:"tagging"`             /*  标签集  */
}

type ZosPutObjectTaggingTaggingRequest struct {
	TagSet []*ZosPutObjectTaggingTaggingTagSetRequest `json:"tagSet"` /*  标签集  */
}

type ZosPutObjectTaggingTaggingTagSetRequest struct {
	Key   string `json:"key,omitempty"`   /*  标签键  */
	Value string `json:"value,omitempty"` /*  标签值  */
}

type ZosPutObjectTaggingResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
