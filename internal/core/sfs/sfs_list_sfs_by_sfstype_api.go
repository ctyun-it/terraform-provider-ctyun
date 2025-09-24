package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// SfsListSfsBySfstypeApi
/* 根据资源池 ID 查询指定存储类型的⽂件系统列表
 */type SfsListSfsBySfstypeApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewSfsListSfsBySfstypeApi(client *core.CtyunClient) *SfsListSfsBySfstypeApi {
	return &SfsListSfsBySfstypeApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/sfs/list-sfs-by-sfstype",
			ContentType:  "application/json",
		},
	}
}

func (a *SfsListSfsBySfstypeApi) Do(ctx context.Context, credential core.Credential, req *SfsListSfsBySfstypeRequest) (*SfsListSfsBySfstypeResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("sfsType", req.SfsType)
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
	var resp SfsListSfsBySfstypeResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type SfsListSfsBySfstypeRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	SfsType  string `json:"sfsType,omitempty"`  /*  文件系统类型(capacity/performance)  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码  */
}

type SfsListSfsBySfstypeResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                                `json:"message"`     /*  响应描述，一般为英文描述  */
	Description string                                `json:"description"` /*  响应描述，一般为中文描述  */
	ReturnObj   *SfsListSfsBySfstypeReturnObjResponse `json:"returnObj"`   /*  参考[returnObj]  */
	ErrorCode   string                                `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                                `json:"error"`       /*  业务细分码，为product.module.code三段式码大驼峰形式  */
}

type SfsListSfsBySfstypeReturnObjResponse struct {
	TotalCount   int32 `json:"totalCount"`   /*  资源池下用户弹性文件总数  */
	CurrentCount int32 `json:"currentCount"` /*  当前页码下查询回来的用户弹性文件数  */
	Total        int32 `json:"total"`        /*  资源池下用户弹性文件总数  */
	PageSize     int32 `json:"pageSize"`     /*  每页包含的元素个数。默认为1  */
	PageNo       int32 `json:"pageNo"`       /*  当前页码。默认为10  */
}
