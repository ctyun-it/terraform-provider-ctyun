package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsListSfsByOndemandApi
/* 查询指定付费类型的⽂件系统列表，付费类型有包年包月、按量付费两种。
 */type SfsListSfsByOndemandApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListSfsByOndemandApi(client *core.CtyunClient) *SfsListSfsByOndemandApi {
	return &SfsListSfsByOndemandApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-sfs-by-ondemand",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListSfsByOndemandApi) Do(ctx context.Context, credential core.Credential, req *SfsListSfsByOndemandRequest) (*SfsListSfsByOndemandResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("onDemand", strconv.FormatBool(req.OnDemand))
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
	var resp SfsListSfsByOndemandResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListSfsByOndemandRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID  */
	OnDemand bool   `json:"onDemand"`           /*  是否按需订购  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数，默认10  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码，默认1  */
}

type SfsListSfsByOndemandResponse struct {
	StatusCode  int32                                  `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                 `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                 `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsListSfsByOndemandReturnObjResponse `json:"returnObj"`   /*  参考returnObj  */
	ErrorCode   string                                 `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                 `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsListSfsByOndemandReturnObjResponse struct {
	Total    int32 `json:"total"`    /*  文件系统读写权限信息总数  */
	PageSize int32 `json:"pageSize"` /*  每页个数  */
	PageNo   int32 `json:"pageNo"`   /*  当前页数  */
}
