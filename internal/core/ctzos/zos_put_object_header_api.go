package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosPutObjectHeaderApi
/* 设置对象的 HTTP 头。
 */type ZosPutObjectHeaderApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosPutObjectHeaderApi(client *core.CtyunClient) *ZosPutObjectHeaderApi {
	return &ZosPutObjectHeaderApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/put-object-header",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosPutObjectHeaderApi) Do(ctx context.Context, credential core.Credential, req *ZosPutObjectHeaderRequest) (*ZosPutObjectHeaderResponse, error) {
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
	var resp ZosPutObjectHeaderResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosPutObjectHeaderRequest struct {
	Bucket   string                            `json:"bucket,omitempty"`   /*  桶名  */
	Key      string                            `json:"key,omitempty"`      /*  对象名  */
	RegionID string                            `json:"regionID,omitempty"` /*  区域 ID  */
	Headers  *ZosPutObjectHeaderHeadersRequest `json:"headers"`            /*  HTTP头，仅限于 CacheControl, ContentDisposition, ContentEncoding, ContentLanguage, ContentType, Expires 六种，至少填写其中一种  */
}

type ZosPutObjectHeaderHeadersRequest struct {
	CacheControl       string `json:"CacheControl,omitempty"`       /*  控制缓存行为的指令  */
	ContentDisposition string `json:"ContentDisposition,omitempty"` /*  指定响应中的文件名和操作行为  */
	ContentEncoding    string `json:"ContentEncoding,omitempty"`    /*  指定响应内容的编码方式  */
	ContentLanguage    string `json:"ContentLanguage,omitempty"`    /*  指定响应内容的语言  */
	ContentType        string `json:"ContentType,omitempty"`        /*  指定响应内容的媒体类型  */
	Expires            string `json:"Expires,omitempty"`            /*  指定响应过期的日期时间  */
}

type ZosPutObjectHeaderResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
