package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListRouteTableApi
/* 查询路由表列表
 */type CtvpcListRouteTableApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListRouteTableApi(client *core.CtyunClient) *CtvpcListRouteTableApi {
	return &CtvpcListRouteTableApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/route-table/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListRouteTableApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListRouteTableRequest) (*CtvpcListRouteTableResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.RouteTableID != nil {
		ctReq.AddParam("routeTableID", *req.RouteTableID)
	}
	if req.RawType != 0 {
		ctReq.AddParam("type", strconv.FormatInt(int64(req.RawType), 10))
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
	var resp CtvpcListRouteTableResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListRouteTableRequest struct {
	RegionID     string  /*  区域id  */
	VpcID        *string /*  关联的vpcID  */
	QueryContent *string /*  对路由表名字 / 路由表描述 / 路由表 id 进行模糊查询  */
	RouteTableID *string /*  路由表 id  */
	RawType      int32   /*  路由表类型:0-子网路由表；2-网关路由表  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListRouteTableResponse struct {
	StatusCode   int32                                   `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                 `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                 `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                 `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	TotalCount   int32                                   `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                                   `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                                   `json:"totalPage"`             /*  总页数  */
	ReturnObj    []*CtvpcListRouteTableReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error        *string                                 `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListRouteTableReturnObjResponse struct {
	Name            *string `json:"name,omitempty"`        /*  路由表名字  */
	Description     *string `json:"description,omitempty"` /*  路由表描述  */
	VpcID           *string `json:"vpcID,omitempty"`       /*  虚拟私有云 id  */
	Id              *string `json:"id,omitempty"`          /*  路由 id  */
	Freezing        *bool   `json:"freezing"`              /*  是否冻结  */
	RouteRulesCount int32   `json:"routeRulesCount"`       /*  路由表中的路由数  */
	CreatedAt       *string `json:"createdAt,omitempty"`   /*  创建时间  */
	UpdatedAt       *string `json:"updatedAt,omitempty"`   /*  更新时间  */
	RawType         int32   `json:"type"`                  /*  路由表类型:0-子网路由表，2-网关路由表  */
	Origin          *string `json:"origin,omitempty"`      /*  路由表来源：default-系统默认; user-用户创建  */
}
