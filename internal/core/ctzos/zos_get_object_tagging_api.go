package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetObjectTaggingApi
/* 查询指导对象的标签。
 */type ZosGetObjectTaggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetObjectTaggingApi(client *core.CtyunClient) *ZosGetObjectTaggingApi {
	return &ZosGetObjectTaggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-object-tagging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetObjectTaggingApi) Do(ctx context.Context, credential core.Credential, req *ZosGetObjectTaggingRequest) (*ZosGetObjectTaggingResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("key", req.Key)
	if req.VersionID != "" {
		ctReq.AddParam("versionID", req.VersionID)
	}
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetObjectTaggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetObjectTaggingRequest struct {
	Bucket    string /*  桶名  */
	Key       string /*  对象名  */
	VersionID string /*  版本ID，在开启多版本时可使用  */
	RegionID  string /*  区域 ID  */
}

type ZosGetObjectTaggingResponse struct {
	StatusCode  int64                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                `json:"message,omitempty"`     /*  状态描述  */
	Description string                                `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetObjectTaggingReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetObjectTaggingReturnObjResponse struct {
	TagSet []*ZosGetObjectTaggingReturnObjTagSetResponse `json:"tagSet"` /*  标签集  */
}

type ZosGetObjectTaggingReturnObjTagSetResponse struct {
	Key   string `json:"key,omitempty"`   /*  标签键  */
	Value string `json:"value,omitempty"` /*  标签值  */
}
