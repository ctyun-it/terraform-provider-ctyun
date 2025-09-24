package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcDhcpoptionsetsqueryApi
/* 查询dhcpoptionsets
 */type CtvpcDhcpoptionsetsqueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcDhcpoptionsetsqueryApi(client *core.CtyunClient) *CtvpcDhcpoptionsetsqueryApi {
	return &CtvpcDhcpoptionsetsqueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/dhcpoptionsets/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcDhcpoptionsetsqueryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcDhcpoptionsetsqueryRequest) (*CtvpcDhcpoptionsetsqueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
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
	var resp CtvpcDhcpoptionsetsqueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcDhcpoptionsetsqueryRequest struct {
	RegionID     string  /*  资源池 ID  */
	QueryContent *string /*  模糊查询  */
	PageNumber   int32   /*  列表的页码，默认值为 1。  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcDhcpoptionsetsqueryResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcDhcpoptionsetsqueryReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcDhcpoptionsetsqueryReturnObjResponse struct {
	Results      []*CtvpcDhcpoptionsetsqueryReturnObjResultsResponse `json:"results"`      /*  dhcpoptionsets组  */
	TotalCount   int32                                               `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                               `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                               `json:"totalPage"`    /*  总页数  */
}

type CtvpcDhcpoptionsetsqueryReturnObjResultsResponse struct {
	DhcpOptionSetsID *string   `json:"dhcpOptionSetsID,omitempty"` /*  dhcpoptionsets  ID  */
	Name             *string   `json:"name,omitempty"`             /*  名字  */
	Description      *string   `json:"description,omitempty"`      /*  描述  */
	DomainName       []*string `json:"domainName"`                 /*  域名  */
	DnsList          []*string `json:"dnsList"`                    /*  ip 列表  */
	VpcList          []*string `json:"vpcList"`                    /*  vpc 列表  */
	CreatedAt        *string   `json:"createdAt,omitempty"`        /*  创建时间  */
	UpdatedAt        *string   `json:"updatedAt,omitempty"`        /*  更新时间  */
}
