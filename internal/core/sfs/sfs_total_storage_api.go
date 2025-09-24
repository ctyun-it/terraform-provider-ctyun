package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// SfsTotalStorageApi
/* 根据资源池 ID 和 文件系统name，查询文件系统总容量
 */type SfsTotalStorageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsTotalStorageApi(client *core.CtyunClient) *SfsTotalStorageApi {
	return &SfsTotalStorageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/total-storage",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsTotalStorageApi) Do(ctx context.Context, credential core.Credential, req *SfsTotalStorageRequest) (*SfsTotalStorageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("name", req.Name)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsTotalStorageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsTotalStorageRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	Name     string `json:"name,omitempty"`     /*  文件系统名称  */
}

type SfsTotalStorageResponse struct {
	StatusCode  int32                             `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                            `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                            `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsTotalStorageReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                            `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                            `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsTotalStorageReturnObjResponse struct {
	Fuid      string `json:"fuid"`      /*  弹性文件功能系统唯一 ID  */
	Name      string `json:"name"`      /*  文件系统名称  */
	RegionID  string `json:"regionID"`  /*  资源池 ID  */
	TotalSize string `json:"totalSize"` /*  文件系统总容量  */
}
