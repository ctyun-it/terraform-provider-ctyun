package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsSfsFreeSizeSfsApi
/* 根据资源池 ID和文件系统ID，查询文件系统剩余容量
 */type SfsSfsFreeSizeSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsFreeSizeSfsApi(client *core.CtyunClient) *SfsSfsFreeSizeSfsApi {
	return &SfsSfsFreeSizeSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/free-storage-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsFreeSizeSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsFreeSizeSfsRequest) (*SfsSfsFreeSizeSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsFreeSizeSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsFreeSizeSfsRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
}

type SfsSfsFreeSizeSfsResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     string                              `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                              `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsFreeSizeSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsFreeSizeSfsReturnObjResponse struct {
	SfsUID   string `json:"sfsUID"`   /*  弹性文件功能系统唯一 ID  */
	SfsName  string `json:"sfsName"`  /*  文件系统名称  */
	RegionID string `json:"regionID"` /*  资源池 ID  */
	FreeSize string `json:"freeSize"` /*  文件系统剩余容量  */
}
