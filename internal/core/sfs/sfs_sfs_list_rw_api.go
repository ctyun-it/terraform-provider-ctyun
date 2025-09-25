package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsSfsListRwApi
/* 查看文件系统只读/读写信息
 */type SfsSfsListRwApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsSfsListRwApi(client *core.CtyunClient) *SfsSfsListRwApi {
	return &SfsSfsListRwApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-rw",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsSfsListRwApi) Do(ctx context.Context, credential core.Credential, req *SfsSfsListRwRequest) (*SfsSfsListRwResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.UID != "" {
		ctReq.AddParam("UID", req.UID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	ctReq.AddParam("pageSize", req.PageSize)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp SfsSfsListRwResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsSfsListRwRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	UID      string `json:"UID,omitempty"`      /*  文件系统ID  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize string `json:"pageSize,omitempty"` /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type SfsSfsListRwResponse struct {
	StatusCode  int32                          `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)   */
	Message     string                         `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                         `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsSfsListRwReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                         `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Error       string                         `json:"error"`       /*  业务细分码，为Product.Module.Code三段式码大驼峰形式  */
}

type SfsSfsListRwReturnObjResponse struct {
	Total int32 `json:"total"` /*  文件系统读写权限信息总数  */
}
