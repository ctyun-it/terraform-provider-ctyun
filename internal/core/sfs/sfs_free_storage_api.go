package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsFreeStorageApi
/* 根据资源池 ID 和文件系统name，查询文件系统剩余容量
 */type SfsFreeStorageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsFreeStorageApi(client *core.CtyunClient) *SfsFreeStorageApi {
	return &SfsFreeStorageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/free-storage",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsFreeStorageApi) Do(ctx context.Context, credential core.Credential, req *SfsFreeStorageRequest) (*SfsFreeStorageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("name", req.Name)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsFreeStorageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsFreeStorageRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	Name     string `json:"name,omitempty"`     /*  文件系统名称  */
}

type SfsFreeStorageResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800为成功，900为失败)  */
	Message     string `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string `json:"description"` /*  响应描述，一般为中文描述  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}
