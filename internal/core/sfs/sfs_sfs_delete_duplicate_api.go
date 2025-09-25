package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsDeleteDuplicateApi
/* 删除文件系统跨域复制配置
 */type SfsSfsDeleteDuplicateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsDeleteDuplicateApi(client *core.CtyunClient) *SfsSfsDeleteDuplicateApi {
	return &SfsSfsDeleteDuplicateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/sfs/delete-duplicate",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsDeleteDuplicateApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsDeleteDuplicateRequest) (*SfsSfsDeleteDuplicateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*SfsSfsDeleteDuplicateRequest
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
	var resp SfsSfsDeleteDuplicateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsDeleteDuplicateRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  源/目的弹性文件功能系统唯一 ID  */
}

type SfsSfsDeleteDuplicateResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Description string                                  `json:"description"` /*  响应描述，一般为英文描述  */
	ReturnObj   *SfsSfsDeleteDuplicateReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                  `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码，参考[结果码]  */
	Error       string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsDeleteDuplicateReturnObjResponse struct {
	RegionID string `json:"regionID"` /*  资源池ID  */
}
