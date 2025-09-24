package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsCreateDuplicateApi
/* 为文件系统创建跨域复制
 */type SfsSfsCreateDuplicateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsCreateDuplicateApi(client *core.CtyunClient) *SfsSfsCreateDuplicateApi {
	return &SfsSfsCreateDuplicateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/create-duplicate",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsCreateDuplicateApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsCreateDuplicateRequest) (*SfsSfsCreateDuplicateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsCreateDuplicateRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsCreateDuplicateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsCreateDuplicateRequest struct {
	SrcSfsUID   string `json:"srcSfsUID,omitempty"`   /*  源文件系统唯一ID  */
	DstSfsUID   string `json:"dstSfsUID,omitempty"`   /*  目标文件系统唯一ID  */
	SrcRegionID string `json:"srcRegionID,omitempty"` /*  源文件系统所在资源池ID  */
	DstRegionID string `json:"dstRegionID,omitempty"` /*  目的文件系统所在资源池ID  */
}

type SfsSfsCreateDuplicateResponse struct {
	ReturnObj   string `json:"returnObj"`   /*  返回参数  */
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*   业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
