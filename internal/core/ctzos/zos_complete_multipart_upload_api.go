package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCompleteMultipartUploadApi
/* 完成分段上传。
 */type ZosCompleteMultipartUploadApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCompleteMultipartUploadApi(client *core.CtyunClient) *ZosCompleteMultipartUploadApi {
	return &ZosCompleteMultipartUploadApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/complete-multipart-upload",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCompleteMultipartUploadApi) Do(ctx context.Context, credential core.Credential, req *ZosCompleteMultipartUploadRequest) (*ZosCompleteMultipartUploadResponse, error) {
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
	var resp ZosCompleteMultipartUploadResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCompleteMultipartUploadRequest struct {
	Bucket          string                                            `json:"bucket,omitempty"`   /*  桶名  */
	Key             string                                            `json:"key,omitempty"`      /*  对象名称  */
	RegionID        string                                            `json:"regionID,omitempty"` /*  区域 ID  */
	MultipartUpload *ZosCompleteMultipartUploadMultipartUploadRequest `json:"multipartUpload"`    /*  分段上传相关信息  */
	UploadID        string                                            `json:"uploadID,omitempty"` /*  分段上传的ID  */
}

type ZosCompleteMultipartUploadMultipartUploadRequest struct {
	Parts []*ZosCompleteMultipartUploadMultipartUploadPartsRequest `json:"parts"` /*  分段信息  */
}

type ZosCompleteMultipartUploadMultipartUploadPartsRequest struct {
	PartNumber int64  `json:"partNumber,omitempty"` /*  分段编号  */
	ETag       string `json:"ETag,omitempty"`       /*  ETag  */
}

type ZosCompleteMultipartUploadResponse struct {
	ReturnObj   *ZosCompleteMultipartUploadReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	StatusCode  int64                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                       `json:"message,omitempty"`     /*  状态描述  */
	Description string                                       `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosCompleteMultipartUploadReturnObjResponse struct {
	ETag      string `json:"ETag,omitempty"`      /*  ETag  */
	Key       string `json:"key,omitempty"`       /*  对象名  */
	Bucket    string `json:"bucket,omitempty"`    /*  桶  */
	VersionID string `json:"versionID,omitempty"` /*  版本ID，在开启多版本时可使用  */
	Location  string `json:"location,omitempty"`  /*  上传的对象位置  */
}
