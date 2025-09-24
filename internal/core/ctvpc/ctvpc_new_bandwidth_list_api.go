package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewBandwidthListApi
/* 调用此接口可查询指定区域下共享带宽实例列表。
 */type CtvpcNewBandwidthListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewBandwidthListApi(client *core.CtyunClient) *CtvpcNewBandwidthListApi {
	return &CtvpcNewBandwidthListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/bandwidth/new-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewBandwidthListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewBandwidthListRequest) (*CtvpcNewBandwidthListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
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
	var resp CtvpcNewBandwidthListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewBandwidthListRequest struct {
	RegionID     string  /*  共享带宽所在的区域id。  */
	QueryContent *string /*  【模糊查询】 共享带宽实例名称 / 带宽 ID  */
	ProjectID    *string /*  企业项目 ID，默认为'0'  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcNewBandwidthListResponse struct {
	StatusCode  int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewBandwidthListReturnObjResponse `json:"returnObj"`             /*  返回查询的共享带宽实例的详细信息。  */
	Error       *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewBandwidthListReturnObjResponse struct {
	TotalCount   int32                                               `json:"totalCount"`   /*  列表条目数。  */
	CurrentCount int32                                               `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                               `json:"totalPage"`    /*  分页查询时总页数。  */
	Bandwidths   []*CtvpcNewBandwidthListReturnObjBandwidthsResponse `json:"bandwidths"`   /*  共享带宽列表  */
}

type CtvpcNewBandwidthListReturnObjBandwidthsResponse struct {
	Id        *string                                                 `json:"id,omitempty"`        /*  共享带宽id。  */
	Status    *string                                                 `json:"status,omitempty"`    /*  共享带宽状态: ACTIVE / EXPIRED / FREEZING  */
	Bandwidth int32                                                   `json:"bandwidth"`           /*  共享带宽的带宽峰值， 单位：Mbps。  */
	Name      *string                                                 `json:"name,omitempty"`      /*  共享带宽名称。  */
	CreatedAt *string                                                 `json:"createdAt,omitempty"` /*  创建时间  */
	ExpireAt  *string                                                 `json:"expireAt,omitempty"`  /*  过期时间  */
	Eips      []*CtvpcNewBandwidthListReturnObjBandwidthsEipsResponse `json:"eips"`                /*  绑定的弹性 IP 列表  */
}

type CtvpcNewBandwidthListReturnObjBandwidthsEipsResponse struct {
	Ip    *string `json:"ip,omitempty"`    /*  弹性 IP 的 IP  */
	EipID *string `json:"eipID,omitempty"` /*  弹性 IP 的 ID  */
}
