package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsSetReadWriteSfsApi
/* 通过资源池ID和文件系统ID，设置文件系统读写
 */type SfsSfsSetReadWriteSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsSetReadWriteSfsApi(client *core.CtyunClient) *SfsSfsSetReadWriteSfsApi {
	return &SfsSfsSetReadWriteSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/set-rw",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsSetReadWriteSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsSetReadWriteSfsRequest) (*SfsSfsSetReadWriteSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsSetReadWriteSfsRequest
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
	var resp SfsSfsSetReadWriteSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsSetReadWriteSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
}

type SfsSfsSetReadWriteSfsResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                  `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                  `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsSetReadWriteSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                  `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsSetReadWriteSfsReturnObjResponse struct {
	OperationID string `json:"operationID"` /*  设置文件系统读写的操作ID  */
}
