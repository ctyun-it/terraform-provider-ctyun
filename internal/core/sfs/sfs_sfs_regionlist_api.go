package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsRegionlistApi
/* 给定文件系统类型，查询支持该文件系统类型的资源池
 */type SfsSfsRegionlistApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsRegionlistApi(client *core.CtyunClient) *SfsSfsRegionlistApi {
	return &SfsSfsRegionlistApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/regionlist",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsRegionlistApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsRegionlistRequest) (*SfsSfsRegionlistResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.FileSystemType != "" {
		ctReq.AddParam("fileSystemType", req.FileSystemType)
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsRegionlistResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsRegionlistRequest struct {
	FileSystemType string `json:"fileSystemType,omitempty"` /*  文件系统类型。取值范围：performance-性能型、capacity-标准型、hdd_e-标准型专属。不传表示查询所有类型  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  分页查询时每页包含的地域数  */
	PageNo         int32  `json:"pageNo,omitempty"`         /*  列表的分页页码  */
}

type SfsSfsRegionlistResponse struct {
	StatusCode  int32  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   string `json:"returnObj"`   /*  returnObj  */
	ErrorCode   string `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码。  */
}
