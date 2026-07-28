package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CfListConcurrencyConfigsApi
/* 查询函数的并发配额列表 */
type CfListConcurrencyConfigsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListConcurrencyConfigsApi(client *core.CtyunClient) *CfListConcurrencyConfigsApi {
	return &CfListConcurrencyConfigsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/resources/functions/reservedConcurrency",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListConcurrencyConfigsApi) Do(ctx context.Context, credential core.Credential, req *CfListConcurrencyConfigsRequest) (*CfListConcurrencyConfigsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.FunctionName != nil && *req.FunctionName != "" {
		ctReq.AddParam("functionName", *req.FunctionName)
	}
	if req.PageIndex != nil && *req.PageIndex != "" {
		ctReq.AddParam("pageIndex", *req.PageIndex)
	}
	if req.PageSize != nil && *req.PageSize != "" {
		ctReq.AddParam("pageSize", *req.PageSize)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListConcurrencyConfigsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListConcurrencyConfigsRequest struct {
	RegionId     string  `json:"regionId"`               /*  资源池ID，请参考<a target="_blank" href="https://www.ctyun.cn/document/10006234/10985593">资源池列表</a>  */
	FunctionName *string `json:"functionName,omitempty"` /*  函数名称，支持模糊匹配，若不指定则列出所有函数  */
	PageIndex    *string `json:"pageIndex,omitempty"`    /*  页码，取值 >= 1，默认值为 1  */
	PageSize     *string `json:"pageSize,omitempty"`     /*  每页大小，取值范围[1, 100]，默认值为 50  */
}

type CfListConcurrencyConfigsResponse struct {
	StatusCode *int32                                     `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Error      *string                                    `json:"error"`      /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                    `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListConcurrencyConfigsReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListConcurrencyConfigsReturnObjResponse struct {
	Data       []*CfListConcurrencyConfigsReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListConcurrencyConfigsReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListConcurrencyConfigsReturnObjDataResponse struct {
	FunctionName        *string `json:"functionName"`        /*  函数名称  */
	ReservedConcurrency *int32  `json:"reservedConcurrency"` /*  最大并发实例数，默认配额时取值范围 [0, 100]  */
}

type CfListConcurrencyConfigsReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}
