package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewIPv6BandwidthListApi
/* 查看 IPv6 带宽列表。
 */type CtvpcNewIPv6BandwidthListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewIPv6BandwidthListApi(client *core.CtyunClient) *CtvpcNewIPv6BandwidthListApi {
	return &CtvpcNewIPv6BandwidthListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ipv6_bandwidth/new-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewIPv6BandwidthListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewIPv6BandwidthListRequest) (*CtvpcNewIPv6BandwidthListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.BandwidthID != nil {
		ctReq.AddParam("bandwidthID", *req.BandwidthID)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcNewIPv6BandwidthListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewIPv6BandwidthListRequest struct {
	RegionID     string  /*  资源池 ID  */
	QueryContent *string /*  【模糊查询】 IPv6 带宽实例名称 / 带宽 ID  */
	BandwidthID  *string /*  IPv6 带宽 ID  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcNewIPv6BandwidthListResponse struct {
	StatusCode  int32                                         `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                       `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                       `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                       `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcNewIPv6BandwidthListReturnObjResponse `json:"returnObj"`             /*  返回查询的共享带宽实例的详细信息。  */
	Error       *string                                       `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewIPv6BandwidthListReturnObjResponse struct {
	TotalCount   int32 `json:"totalCount"`   /*  总数  */
	CurrentCount int32 `json:"currentCount"` /*  分页查询时每页的行数  */
	TotalPage    int32 `json:"totalPage"`    /*  分页查询时总页数  */
}
