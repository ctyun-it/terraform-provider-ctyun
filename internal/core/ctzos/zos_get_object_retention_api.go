package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetObjectRetentionApi
/* 获取对象的保留期限设置。
 */type ZosGetObjectRetentionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetObjectRetentionApi(client *core.CtyunClient) *ZosGetObjectRetentionApi {
	return &ZosGetObjectRetentionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/get-object-retention",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetObjectRetentionApi) Do(ctx context.Context, credential core.Credential, req *ZosGetObjectRetentionRequest) (*ZosGetObjectRetentionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("bucket", req.Bucket)
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("key", req.Key)
	ctReq.AddParam("versionID", req.VersionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetObjectRetentionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetObjectRetentionRequest struct {
	Bucket    string /*  桶名  */
	RegionID  string /*  区域 ID  */
	Key       string /*  对象名称  */
	VersionID string /*  版本ID，在开启多版本时可使用  */
}

type ZosGetObjectRetentionResponse struct {
	StatusCode  int64                                   `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                  `json:"message,omitempty"`     /*  状态描述  */
	Description string                                  `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ReturnObj   *ZosGetObjectRetentionReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	ErrorCode   string                                  `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                  `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGetObjectRetentionReturnObjResponse struct {
	Retention *ZosGetObjectRetentionReturnObjRetentionResponse `json:"retention"` /*  保留配置  */
}

type ZosGetObjectRetentionReturnObjRetentionResponse struct {
	Mode            string `json:"mode,omitempty"`            /*  保留模式，COMPLIANCE 或 GOVERNANCE  */
	RetainUntilDate int64  `json:"retainUntilDate,omitempty"` /*  保留截止日期  */
}
