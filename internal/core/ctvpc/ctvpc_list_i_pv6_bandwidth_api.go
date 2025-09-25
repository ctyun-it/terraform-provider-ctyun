package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListIPv6BandwidthApi
/* 查看 IPv6 带宽列表。
 */type CtvpcListIPv6BandwidthApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListIPv6BandwidthApi(client *core.CtyunClient) *CtvpcListIPv6BandwidthApi {
	return &CtvpcListIPv6BandwidthApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ipv6_bandwidth/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListIPv6BandwidthApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListIPv6BandwidthRequest) (*CtvpcListIPv6BandwidthResponse, error) {
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
	var resp CtvpcListIPv6BandwidthResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListIPv6BandwidthRequest struct {
	RegionID     string  /*  资源池 ID  */
	QueryContent *string /*  【模糊查询】 IPv6 带宽实例名称 / 带宽 ID  */
	BandwidthID  *string /*  IPv6 带宽 ID  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListIPv6BandwidthResponse struct {
	StatusCode   int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListIPv6BandwidthReturnObjResponse `json:"returnObj"`             /*  返回查询的共享带宽实例的详细信息。  */
	TotalCount   int32                                      `json:"totalCount"`            /*  总数  */
	CurrentCount int32                                      `json:"currentCount"`          /*  分页查询时每页的行数  */
	TotalPage    int32                                      `json:"totalPage"`             /*  分页查询时总页数  */
	Error        *string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListIPv6BandwidthReturnObjResponse struct {
	Id            *string `json:"id,omitempty"`            /*  IPv6 带宽 ID  */
	Status        *string `json:"status,omitempty"`        /*  IPv6 带宽状态: ACTIVE（正常） / EXPIRED（过期） / FREEZING（冻结） /CREATEING（创建中）  */
	Name          *string `json:"name,omitempty"`          /*  IPv6 带宽名字  */
	Bandwidth     int32   `json:"bandwidth"`               /*  IPv6 带宽峰值 mbps  */
	ResourceSpec  *string `json:"resourceSpec,omitempty"`  /*  独享 / 共享  */
	PaymentType   *string `json:"paymentType,omitempty"`   /*  计费类型  */
	CreatedTime   *string `json:"createdTime,omitempty"`   /*  IPv6 带宽创建时间  */
	ExpiredTime   *string `json:"expiredTime,omitempty"`   /*  IPv6 带宽过期时间  */
	IpAddress     *string `json:"ipAddress,omitempty"`     /*  IP 地址  */
	Ipv6GatewayID *string `json:"ipv6GatewayID,omitempty"` /*  IPv6 网关 ID  */
}
