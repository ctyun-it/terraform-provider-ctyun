package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListRegionApi
/* 查询并行文件支持的地域
 */type HpfsListRegionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListRegionApi(client *core.CtyunClient) *HpfsListRegionApi {
	return &HpfsListRegionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-region",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListRegionApi) Do(ctx context.Context, credential core.Credential, req *HpfsListRegionRequest) (*HpfsListRegionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
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
	var resp HpfsListRegionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListRegionRequest struct {
	PageSize int32 `json:"pageSize,omitempty"` /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo   int32 `json:"pageNo,omitempty"`   /*  列表的分页页码，默认值为1  */
}

type HpfsListRegionResponse struct {
	StatusCode  int32                            `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                           `json:"message"`     /*  响应描述  */
	Description string                           `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListRegionReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                           `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                           `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListRegionReturnObjResponse struct {
	RegionList   []*HpfsListRegionReturnObjRegionListResponse `json:"regionList"`   /*  查询的地域详情列表  */
	TotalCount   int32                                        `json:"totalCount"`   /*  支持并行文件的地域总数  */
	CurrentCount int32                                        `json:"currentCount"` /*  当前页码的元素个数  */
	PageSize     int32                                        `json:"pageSize"`     /*  每页个数  */
	PageNo       int32                                        `json:"pageNo"`       /*  当前页数  */
}

type HpfsListRegionReturnObjRegionListResponse struct {
	RegionID   string `json:"regionID"`   /*  资源池ID  */
	RegionName string `json:"regionName"` /*  资源池名字  */
}
