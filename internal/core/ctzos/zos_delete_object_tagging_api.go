package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosDeleteObjectTaggingApi
/* 删除指定对象的标签。
 */type ZosDeleteObjectTaggingApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosDeleteObjectTaggingApi(client *core.CtyunClient) *ZosDeleteObjectTaggingApi {
	return &ZosDeleteObjectTaggingApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/delete-object-tagging",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosDeleteObjectTaggingApi) Do(ctx context.Context, credential core.Credential, req *ZosDeleteObjectTaggingRequest) (*ZosDeleteObjectTaggingResponse, error) {
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
	var resp ZosDeleteObjectTaggingResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosDeleteObjectTaggingRequest struct {
	Bucket    string `json:"bucket,omitempty"`    /*  桶名  */
	Key       string `json:"key,omitempty"`       /*  对象名  */
	VersionID string `json:"versionID,omitempty"` /*  版本ID，在开启多版本时可使用  */
	RegionID  string `json:"regionID,omitempty"`  /*  区域 ID  */
}

type ZosDeleteObjectTaggingResponse struct {
	StatusCode  int64  `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
