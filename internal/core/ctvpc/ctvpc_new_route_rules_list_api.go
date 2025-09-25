package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewRouteRulesListApi
/* 查询路由表规则列表
 */type CtvpcNewRouteRulesListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewRouteRulesListApi(client *core.CtyunClient) *CtvpcNewRouteRulesListApi {
	return &CtvpcNewRouteRulesListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/route-table/new-list-rules",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewRouteRulesListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewRouteRulesListRequest) (*CtvpcNewRouteRulesListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("routeTableID", req.RouteTableID)
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
	var resp CtvpcNewRouteRulesListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewRouteRulesListRequest struct {
	RegionID     string /*  区域id  */
	RouteTableID string /*  路由表 id  */
	PageNumber   int32  /*  列表的页码，默认值为 1。  */
	PageNo       int32  /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcNewRouteRulesListResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewRouteRulesListReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewRouteRulesListReturnObjResponse struct {
	RouteRules   []*CtvpcNewRouteRulesListReturnObjRouteRulesResponse `json:"routeRules"`   /*  路由规则  */
	TotalCount   int32                                                `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                `json:"totalPage"`    /*  总页数  */
}

type CtvpcNewRouteRulesListReturnObjRouteRulesResponse struct {
	NextHopID   *string `json:"nextHopID,omitempty"`   /*  下一跳设备 id  */
	NextHopType *string `json:"nextHopType,omitempty"` /*  vpcpeering / havip / bm / vm / natgw/ igw6 / dc / ticc / vpngw / enic  */
	Destination *string `json:"destination,omitempty"` /*  无类别域间路由  */
	IpVersion   int32   `json:"ipVersion"`             /*  4 表示 ipv4, 6 表示 ipv6  */
	Description *string `json:"description,omitempty"` /*  规则描述  */
	RouteRuleID *string `json:"routeRuleID,omitempty"` /*  路由规则 id  */
}
