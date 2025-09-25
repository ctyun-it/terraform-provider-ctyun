package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCreateMultipartUploadApi
/* 初始化分段上传事件。请注意 OpenAPI 本身不具备分段上传能力，初始化分段上传后需根据返回信息借助 sdk 等工具进行分段上传。
 */type ZosCreateMultipartUploadApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCreateMultipartUploadApi(client *core.CtyunClient) *ZosCreateMultipartUploadApi {
	return &ZosCreateMultipartUploadApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/create-multipart-upload",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCreateMultipartUploadApi) Do(ctx context.Context, credential core.Credential, req *ZosCreateMultipartUploadRequest) (*ZosCreateMultipartUploadResponse, error) {
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
	var resp ZosCreateMultipartUploadResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCreateMultipartUploadRequest struct {
	Bucket       string `json:"bucket,omitempty"`       /*  桶名  */
	RegionID     string `json:"regionID,omitempty"`     /*  区域 ID  */
	Key          string `json:"key,omitempty"`          /*  对象名  */
	ACL          string `json:"ACL,omitempty"`          /*  ACL，可选的有 private，public-read，public-read-write，authenticated-read，默认private  */
	StorageClass string `json:"storageClass,omitempty"` /*  存储类，可选的有 STANDARD，STANDARD_IA，GLACIER，默认STANDARD  */
}

type ZosCreateMultipartUploadResponse struct {
	ReturnObj   *ZosCreateMultipartUploadReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	StatusCode  int64                                      `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                     `json:"message,omitempty"`     /*  状态描述  */
	Description string                                     `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                                     `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                     `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosCreateMultipartUploadReturnObjResponse struct {
	Bucket   string `json:"bucket,omitempty"`   /*  桶名  */
	Key      string `json:"key,omitempty"`      /*  对象名  */
	UploadID string `json:"uploadID,omitempty"` /*  分段上传ID  */
}
