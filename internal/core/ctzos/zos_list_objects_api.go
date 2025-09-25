package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListObjectsApi
/* 查询桶内的所有对象。
 */type ZosListObjectsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListObjectsApi(client *core.CtyunClient) *ZosListObjectsApi {
	return &ZosListObjectsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-objects",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListObjectsApi) Do(ctx context.Context, credential core.Credential, req *ZosListObjectsRequest) (*ZosListObjectsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	if req.Delimiter != "" {
		ctReq.AddParam("delimiter", req.Delimiter)
	}
	if req.Marker != "" {
		ctReq.AddParam("marker", req.Marker)
	}
	if req.MaxKeys != 0 {
		ctReq.AddParam("maxKeys", strconv.FormatInt(int64(req.MaxKeys), 10))
	}
	if req.Prefix != "" {
		ctReq.AddParam("prefix", req.Prefix)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListObjectsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListObjectsRequest struct {
	Bucket    string /*  桶名  */
	RegionID  string /*  区域 ID  */
	Delimiter string /*  定界符是您用来对键进行分组的字符  */
	Marker    string /*  指示从哪里开始列出。ZOS会在这个指定的对象之后开始列出。标记可以是桶中的任何对象  */
	MaxKeys   int64  /*  一次返回keys的最大数目（默认值和上限为1000）  */
	Prefix    string /*  返回的key的前缀  */
}

type ZosListObjectsResponse struct {
	StatusCode  int64                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                           `json:"message,omitempty"`     /*  状态描述  */
	Description string                           `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListObjectsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                           `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                           `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListObjectsReturnObjResponse struct {
	IsTruncated    *bool                                            `json:"isTruncated"`            /*  指示返回的分段列表是否被截断。值为 true 表示列表已被截断。  */
	Marker         string                                           `json:"marker,omitempty"`       /*  指示存储桶列表中的开始位置。如果标记是随请求一起发送的，则标记包含在响应中。  */
	Contents       []*ZosListObjectsReturnObjContentsResponse       `json:"contents"`               /*  对象的容器  */
	NextMarker     string                                           `json:"nextMarker,omitempty"`   /*  下一个 Marker  */
	Name           string                                           `json:"name,omitempty"`         /*  桶名  */
	Prefix         string                                           `json:"prefix,omitempty"`       /*  返回的key的前缀  */
	Delimiter      string                                           `json:"delimiter,omitempty"`    /*  定界符是您用来对键进行分组的字符  */
	MaxKeys        int64                                            `json:"maxKeys,omitempty"`      /*  一次返回keys的最大数目  */
	CommonPrefixes []*ZosListObjectsReturnObjCommonPrefixesResponse `json:"commonPrefixes"`         /*  在计算返回数时，所有键（最多 1,000 个）汇总在一个公共前缀中计为单个返回。仅当指定分隔符时，响应才能包含 commonPrefixes  */
	EncodingType   string                                           `json:"encodingType,omitempty"` /*  对响应中的对象键进行编码并指定要使用的编码方法，目前允许的值只有 url  */
}

type ZosListObjectsReturnObjContentsResponse struct {
	ETag         string                                        `json:"ETag,omitempty"`         /*  ETag  */
	Key          string                                        `json:"key,omitempty"`          /*  对象名  */
	LastModified string                                        `json:"lastModified,omitempty"` /*  最后更改时间, ISO-8601 格式的日期字符串  */
	Owner        *ZosListObjectsReturnObjContentsOwnerResponse `json:"owner"`                  /*  所有者  */
	Size         int64                                         `json:"size,omitempty"`         /*  大小  */
	StorageClass string                                        `json:"storageClass,omitempty"` /*  存储类，可能的值有：STANDARD（标准存储）、STANDARD_IA（低频存储）、GLACIER（归档存储）  */
	RawType      string                                        `json:"type,omitempty"`         /*  对象类型，普通对象为 `Normal`，软链接对象为 `Symlink`  */
}

type ZosListObjectsReturnObjCommonPrefixesResponse struct {
	Prefix string `json:"prefix,omitempty"` /*  前缀  */
}

type ZosListObjectsReturnObjContentsOwnerResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名  */
	ID          string `json:"ID,omitempty"`          /*  用户 ID  */
}
