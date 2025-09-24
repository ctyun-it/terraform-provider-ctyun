package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListMountpointApi
/* 根据文件系统sfsUID，查询指定文件系统挂载点列表
 */type SfsSfsListMountpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListMountpointApi(client *core.CtyunClient) *SfsSfsListMountpointApi {
	return &SfsSfsListMountpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-mountpoint-sfs",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListMountpointApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListMountpointRequest) (*SfsSfsListMountpointResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsUID", req.SfsUID)
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
	var resp SfsSfsListMountpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListMountpointRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsUID   string `json:"sfsUID,omitempty"`   /*  弹性文件功能系统唯一 ID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的数量  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  当前页码  */
}

type SfsSfsListMountpointResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                 `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListMountpointReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码.参考[结果码]  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsSfsListMountpointReturnObjResponse struct {
	MountPointList string `json:"mountPointList"` /*  文件系统挂载点ID列表  */
	CurrentCount   int32  `json:"currentCount"`   /*  当前查询到的文件系统挂载点ID个数  */
	TotalCount     int32  `json:"totalCount"`     /*  资源池下文件系统挂载点ID总数  */
	PageSize       int32  `json:"pageSize"`       /*  每页个数  */
	PageNo         int32  `json:"pageNo"`         /*  当前页数  */
}
