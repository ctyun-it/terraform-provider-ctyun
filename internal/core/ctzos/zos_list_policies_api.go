package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// ZosListPoliciesApi
/* 获取策略列表。
 */type ZosListPoliciesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosListPoliciesApi(client *core.CtyunClient) *ZosListPoliciesApi {
	return &ZosListPoliciesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/oss/list-policies",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosListPoliciesApi) Do(ctx context.Context, credential core.Credential, req *ZosListPoliciesRequest) (*ZosListPoliciesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.Keyword != "" {
		ctReq.AddParam("keyword", req.Keyword)
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosListPoliciesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosListPoliciesRequest struct {
	RegionID string /*  区域 ID  */
	Keyword  string /*  模糊查询策略名，不区分大小写  */
	PageSize int64  /*  单页数量，取值范围 1~50，默认值为10  */
	Page     int64  /*  页码, 若与参数 pageNo 同时存在，以 pageNo 为准，默认值为1  */
	PageNo   int64  /*  页码，默认值为1  */
}

type ZosListPoliciesResponse struct {
	StatusCode  int64                             `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为处理中/失败）  */
	Message     string                            `json:"message,omitempty"`     /*  状态描述  */
	ReturnObj   *ZosListPoliciesReturnObjResponse `json:"returnObj"`             /*  响应对象  */
	Description string                            `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string                            `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为 product.module.code 三段式码  */
	Error       string                            `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}

type ZosListPoliciesReturnObjResponse struct {
	TotalCount   int64                                     `json:"totalCount,omitempty"`   /*  总数  */
	CurrentCount int64                                     `json:"currentCount,omitempty"` /*  当前页记录数  */
	Result       []*ZosListPoliciesReturnObjResultResponse `json:"result"`                 /*  查询结果列表  */
}

type ZosListPoliciesReturnObjResultResponse struct {
	Fuser_last_updated string `json:"fuser_last_updated,omitempty"` /*  最近更新时间  */
	Policy_document    string `json:"policy_document,omitempty"`    /*  策略  */
	Note               string `json:"note,omitempty"`               /*  策略备注  */
	Created_time       string `json:"created_time,omitempty"`       /*  策略创建时间  */
	Policy_name        string `json:"policy_name,omitempty"`        /*  策略名  */
}
