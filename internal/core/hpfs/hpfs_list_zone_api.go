package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// HpfsListZoneApi
/* 查询一个地域下的所有支持并行文件的可用区及该可用区所支持的文件系统类型
 */type HpfsListZoneApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewHpfsListZoneApi(client *core.CtyunClient) *HpfsListZoneApi {
	return &HpfsListZoneApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/hpfs/list-zone",
			ContentType:  "application/json",
		},
	}
}

func (a *HpfsListZoneApi) Do(ctx context.Context, credential core.Credential, req *HpfsListZoneRequest) (*HpfsListZoneResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
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
	var resp HpfsListZoneResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type HpfsListZoneRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页包含的元素个数范围(1-50)，默认值为10  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  列表的分页页码，默认值为1  */
}

type HpfsListZoneResponse struct {
	StatusCode  int32                          `json:"statusCode"`  /*  返回状态码(800 为成功，900 为失败)  */
	Message     string                         `json:"message"`     /*  响应描述  */
	Description string                         `json:"description"` /*  响应描述  */
	ReturnObj   *HpfsListZoneReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                         `json:"errorCode"`   /*  业务细分码，为 product.module.code 三段式码  */
	Error       string                         `json:"error"`       /*  业务细分码，为product.module.code三段式大驼峰码  */
}

type HpfsListZoneReturnObjResponse struct {
	ZoneList     []*HpfsListZoneReturnObjZoneListResponse `json:"zoneList"`     /*  查询的可用区列表  */
	TotalCount   int32                                    `json:"totalCount"`   /*  当前资源池下可用区总数  */
	CurrentCount int32                                    `json:"currentCount"` /*  当前页码的元素个数  */
	PageSize     int32                                    `json:"pageSize"`     /*  每页个数  */
	PageNo       int32                                    `json:"pageNo"`       /*  当前页数  */
}

type HpfsListZoneReturnObjZoneListResponse struct {
	AzName        string   `json:"azName"`        /*  可用区名称，其他需要可用区参数的接口需要依赖该名称结果  */
	AzDisplayName string   `json:"azDisplayName"` /*  可用区展示名  */
	StorageTypes  []string `json:"storageTypes"`  /*  可用区支持的存储类型  */
}
