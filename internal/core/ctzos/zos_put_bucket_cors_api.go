package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutBucketCorsApi
/* 设置指定桶的跨域访问策略。
 */type ZosPutBucketCorsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutBucketCorsApi(client *core.CtyunClient) *ZosPutBucketCorsApi {
	return &ZosPutBucketCorsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-bucket-cors",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutBucketCorsApi) Do(ctx context.Context, credential core.Credential, req *ZosPutBucketCorsRequest) (*ZosPutBucketCorsResponse, error) {
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
	var resp ZosPutBucketCorsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutBucketCorsRequest struct {
	Bucket            string                                    `json:"bucket,omitempty"`   /*  桶名  */
	RegionID          string                                    `json:"regionID,omitempty"` /*  区域 ID  */
	CORSConfiguration *ZosPutBucketCorsCORSConfigurationRequest `json:"CORSConfiguration"`  /*  跨域访问策略  */
}

type ZosPutBucketCorsCORSConfigurationRequest struct {
	CORSRules []*ZosPutBucketCorsCORSConfigurationCORSRulesRequest `json:"CORSRules"` /*  规则  */
}

type ZosPutBucketCorsCORSConfigurationCORSRulesRequest struct {
	AllowedHeaders []string `json:"allowedHeaders"`          /*  Access-Control-Request-Headers 标头中指定的标头  */
	AllowedMethods []string `json:"allowedMethods"`          /*  您允许源执行的 HTTP 方法。有效值为 GET 、 PUT 、 HEAD 、 POST 和 DELETE  */
	AllowedOrigins []string `json:"allowedOrigins"`          /*  您希望用户能够从中访问存储桶的一个或多个来源  */
	ExposeHeaders  []string `json:"exposeHeaders"`           /*  您希望用户能够从他们的应用程序（例如，从 JavaScript XMLHttpRequest 对象）访问的响应中的一个或多个标头。  */
	MaxAgeSeconds  int64    `json:"maxAgeSeconds,omitempty"` /*  浏览器缓存指定资源的预检响应的时间（以秒为单位）  */
}

type ZosPutBucketCorsResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
