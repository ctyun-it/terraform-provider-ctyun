package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsTotalSizeSfsApi
/* 根据资源池 ID 和 文件系统sfsUID，查询文件系统总容量
 */type SfsSfsTotalSizeSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsTotalSizeSfsApi(client *core.CtyunClient) *SfsSfsTotalSizeSfsApi {
	return &SfsSfsTotalSizeSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/total-storage-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsTotalSizeSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsTotalSizeSfsRequest) (*SfsSfsTotalSizeSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsTotalSizeSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsTotalSizeSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
}

type SfsSfsTotalSizeSfsResponse struct {
	StatusCode  int32                                `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                               `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                               `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsTotalSizeSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                               `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                               `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsTotalSizeSfsReturnObjResponse struct {
	SfsUID    string `json:"sfsUID"`    /*  弹性文件功能系统唯一 ID  */
	SfsName   string `json:"sfsName"`   /*  文件系统名称  */
	RegionID  string `json:"regionID"`  /*  资源池 ID  */
	TotalSize string `json:"totalSize"` /*  文件系统总容量  */
}
