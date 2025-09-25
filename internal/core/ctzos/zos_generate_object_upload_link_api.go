package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGenerateObjectUploadLinkApi
/* 生成一个临时的预签名的Url，没有权限访问集群的用户可以通过该Url上传文件到集群(通过http的post方式)。
 */type ZosGenerateObjectUploadLinkApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGenerateObjectUploadLinkApi(client *core.CtyunClient) *ZosGenerateObjectUploadLinkApi {
	return &ZosGenerateObjectUploadLinkApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/oss/generate-object-upload-link",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGenerateObjectUploadLinkApi) Do(ctx context.Context, credential core.Credential, req *ZosGenerateObjectUploadLinkRequest) (*ZosGenerateObjectUploadLinkResponse, error) {
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
	var resp ZosGenerateObjectUploadLinkResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGenerateObjectUploadLinkRequest struct {
	RegionID  string `json:"regionID,omitempty"`  /*  区域 ID  */
	Bucket    string `json:"bucket,omitempty"`    /*  桶名  */
	Key       string `json:"key,omitempty"`       /*  对象名  */
	ExpiresIn int64  `json:"expiresIn,omitempty"` /*  url 过期时间，默认 3600  */
}

type ZosGenerateObjectUploadLinkResponse struct {
	ReturnObj   *ZosGenerateObjectUploadLinkReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	StatusCode  int64                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                                        `json:"message,omitempty"`     /*  状态描述  */
	Description string                                        `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosGenerateObjectUploadLinkReturnObjResponse struct {
	Url    string                                              `json:"url,omitempty"` /*  URL链接  */
	Fields *ZosGenerateObjectUploadLinkReturnObjFieldsResponse `json:"fields"`        /*  字段，直接将此参数作为 POST 请求的 form 参数即可，使用示例如：`curl -v -i -X POST -H "Content-Type: multipart/form-data" -F "signature=xxx" -F "AWSAccessKeyId=xxx" -F "key=xxx" -F "policy=xxx" -F "file=@filepath"  https://xxx`  */
}

type ZosGenerateObjectUploadLinkReturnObjFieldsResponse struct {
	Policy         string `json:"policy,omitempty"`         /*  策略  */
	AWSAccessKeyId string `json:"AWSAccessKeyId,omitempty"` /*  ak  */
	Key            string `json:"key,omitempty"`            /*  对象名  */
	Signature      string `json:"signature,omitempty"`      /*  签名  */
}
