package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsListRegionSfsApi
/* 查询弹性文件系统支持的地域
 */type SfsListRegionSfsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListRegionSfsApi(client *core.CtyunClient) *SfsListRegionSfsApi {
	return &SfsListRegionSfsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-region-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListRegionSfsApi) Do(ctx context.Context, credential core.Credential, req *SfsListRegionSfsRequest) (*SfsListRegionSfsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsListRegionSfsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListRegionSfsRequest struct{}

type SfsListRegionSfsResponse struct {
	StatusCode  int32                              `json:"statusCode"`  /*   返回状态码(800 为成功，900 为失败)  */
	Message     string                             `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                             `json:"description"` /*  响应描述，一般为英文描述  */
	ReturnObj   *SfsListRegionSfsReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                             `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                             `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsListRegionSfsReturnObjResponse struct {
	RegionList []*SfsListRegionSfsReturnObjRegionListResponse `json:"regionList"` /*  查询的地域详情列表  */
}

type SfsListRegionSfsReturnObjRegionListResponse struct {
	UUID       string `json:"UUID"`       /*  资源池uuid     */
	RegionName string `json:"regionName"` /*  资源池名字  */
}
