package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcPrefixlistQueryApi
/* 查看 prefixlist 列表信息
 */type CtvpcPrefixlistQueryApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistQueryApi(client *core.CtyunClient) *CtvpcPrefixlistQueryApi {
	return &CtvpcPrefixlistQueryApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/prefixlist/query",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistQueryApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistQueryRequest) (*CtvpcPrefixlistQueryResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PrefixListID != nil {
		ctReq.AddParam("prefixListID", *req.PrefixListID)
	}
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
	var resp CtvpcPrefixlistQueryResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistQueryRequest struct {
	RegionID     string  /*  资源池ID  */
	PrefixListID *string /*  prefixlistID  */
	QueryContent *string /*  模糊查询关键字  */
	PageNumber   int32   /*  列表的页码，默认值为 1  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10  */
}

type CtvpcPrefixlistQueryResponse struct {
	StatusCode  int32                                  `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcPrefixlistQueryReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcPrefixlistQueryReturnObjResponse struct {
	PrefixList   []*CtvpcPrefixlistQueryReturnObjPrefixListResponse `json:"prefixList"`   /*  prefixList  */
	TotalCount   int32                                              `json:"totalCount"`   /*  列表条目数。  */
	CurrentCount int32                                              `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                              `json:"totalPage"`    /*  分页查询时总页数。  */
}

type CtvpcPrefixlistQueryReturnObjPrefixListResponse struct {
	PrefixListID    *string                                                           `json:"prefixListID,omitempty"` /*  prefixlist id  */
	Name            *string                                                           `json:"name,omitempty"`         /*  prefixlist 名称  */
	Limit           int32                                                             `json:"limit"`                  /*  前缀列表支持的最大条目容量  */
	AddressType     *string                                                           `json:"addressType,omitempty"`  /*  地址类型，4：ipv4，6：ipv6  */
	Description     *string                                                           `json:"description,omitempty"`  /*  描述  */
	CreatedAt       *string                                                           `json:"createdAt,omitempty"`    /*  创建时间  */
	UpdatedAt       *string                                                           `json:"updatedAt,omitempty"`    /*  更新时间  */
	PrefixListRules []*CtvpcPrefixlistQueryReturnObjPrefixListPrefixListRulesResponse `json:"prefixListRules"`        /*  前缀列表规则  */
}

type CtvpcPrefixlistQueryReturnObjPrefixListPrefixListRulesResponse struct {
	PrefixListRuleID *string `json:"prefixListRuleID,omitempty"` /*  prefixlistRuleID  */
	Cidr             *string `json:"cidr,omitempty"`             /*  前缀列表条目,cidr  */
	Description      *string `json:"description,omitempty"`      /*  描述  */
}
