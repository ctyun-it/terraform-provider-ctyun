package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosListObjectVersionsApi
/* 查询对象的版本信息。
 */type ZosListObjectVersionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListObjectVersionsApi(client *core.CtyunClient) *ZosListObjectVersionsApi {
	return &ZosListObjectVersionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-object-versions",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListObjectVersionsApi) Do(ctx context.Context, credential core.Credential, req *ZosListObjectVersionsRequest) (*ZosListObjectVersionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	if req.KeyMarker != "" {
		ctReq.AddParam("keyMarker", req.KeyMarker)
	}
	if req.Prefix != "" {
		ctReq.AddParam("prefix", req.Prefix)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListObjectVersionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListObjectVersionsRequest struct {
	Bucket    string /*  桶名  */
	RegionID  string /*  区域 ID  */
	KeyMarker string /*  指示从哪里开始列出，ZOS 会在这个指定的键之后开始列出，即列出名称大于keyMarker的对象。标记可以是桶中的任何键。默认值为空字符串 ""。  */
	Prefix    string /*  返回的key的前缀，注意 key 与 prefix 完全相同的对象不会返回。默认值为空字符串 ""。  */
}

type ZosListObjectVersionsResponse struct {
	StatusCode  int64                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  状态描述  */
	Description string                                  `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosListObjectVersionsReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListObjectVersionsReturnObjResponse struct {
	Keys          []*ZosListObjectVersionsReturnObjKeysResponse          `json:"keys"`                    /*  具有 prefix 前缀的对象列表。  */
	NextKeyMarker string                                                 `json:"nextKeyMarker,omitempty"` /*  下一个 keyMarker。若没有则为 null。  */
	KeyMarker     string                                                 `json:"keyMarker,omitempty"`     /*  标记截断响应中返回的最后一个 key。  */
	DeleteMarkers []*ZosListObjectVersionsReturnObjDeleteMarkersResponse `json:"deleteMarkers"`           /*  作为删除标记的对象的容器。  */
}

type ZosListObjectVersionsReturnObjKeysResponse struct {
	Key      string `json:"key,omitempty"`      /*  对象名  */
	Versions string `json:"versions,omitempty"` /*  版本信息  */
}

type ZosListObjectVersionsReturnObjDeleteMarkersResponse struct {
	LastModified string                                                    `json:"lastModified,omitempty"` /*  最后修改时间，为 ISO 8601 格式  */
	VersionID    string                                                    `json:"versionID,omitempty"`    /*  版本ID，在开启多版本时可使用  */
	Key          string                                                    `json:"key,omitempty"`          /*  对象名  */
	Owner        *ZosListObjectVersionsReturnObjDeleteMarkersOwnerResponse `json:"owner"`                  /*  创建删除标记的账户  */
	IsLatest     *bool                                                     `json:"isLatest"`               /*  是否为最新版本  */
}

type ZosListObjectVersionsReturnObjDeleteMarkersOwnerResponse struct {
	DisplayName string `json:"displayName,omitempty"` /*  展示名  */
	ID          string `json:"ID,omitempty"`          /*  用户 ID  */
}
