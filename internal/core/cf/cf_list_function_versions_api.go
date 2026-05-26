package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListFunctionVersionsApi
/* 分页查询版本列表 */
type CfListFunctionVersionsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListFunctionVersionsApi(client *core.CtyunClient) *CfListFunctionVersionsApi {
	return &CfListFunctionVersionsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/functions/{functionName}/versions",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListFunctionVersionsApi) Do(ctx context.Context, credential core.Credential, req *CfListFunctionVersionsRequest) (*CfListFunctionVersionsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("functionName", req.FunctionName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PageIndex != nil && *req.PageIndex != 0 {
		ctReq.AddParam("pageIndex", strconv.FormatInt(int64(*req.PageIndex), 10))
	}
	if req.PageSize != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	if req.Search != nil && *req.Search != "" {
		ctReq.AddParam("search", *req.Search)
	}
	if req.OrderBy != nil && *req.OrderBy != "" {
		ctReq.AddParam("orderBy", *req.OrderBy)
	}
	if req.Order != nil && *req.Order != "" {
		ctReq.AddParam("order", *req.Order)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListFunctionVersionsResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListFunctionVersionsRequest struct {
	FunctionName string  `json:"functionName"`        /*  函数名称  */
	RegionId     string  `json:"regionId"`            /*  资源池id  */
	PageIndex    *int32  `json:"pageIndex,omitempty"` /*  分页页码，默认为1  */
	PageSize     *int32  `json:"pageSize,omitempty"`  /*  分页大小，默认为10  */
	Search       *string `json:"search,omitempty"`    /*  模糊查询的关键字，默认为空  */
	OrderBy      *string `json:"orderBy,omitempty"`   /*  排序字段名，默认按版本号show_version_id，还支持创建时间created_at  */
	Order        *string `json:"order,omitempty"`     /*  ASC升序，DESC降序，默认降序  */
}

type CfListFunctionVersionsResponse struct {
	StatusCode *int32                                   `json:"statusCode"` /*  状态码。0表示成功，其他值表示失败  */
	Code       *string                                  `json:"code"`       /*  错误码。CF_0表示成功，其他值表示失败  */
	Message    *string                                  `json:"message"`    /*  错误描述信息  */
	ReturnObj  *CfListFunctionVersionsReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListFunctionVersionsReturnObjResponse struct {
	Data       []*CfListFunctionVersionsReturnObjDataResponse     `json:"data"`       /*  分页数据  */
	Pagination *CfListFunctionVersionsReturnObjPaginationResponse `json:"pagination"` /*  分页信息  */
}

type CfListFunctionVersionsReturnObjDataResponse struct {
	Description       *string                                                       `json:"description"`       /*  描述  */
	VersionId         *string                                                       `json:"versionId"`         /*  版本 ID  */
	CreateTime        *string                                                       `json:"createTime"`        /*  创建时间  */
	UpdateTime        *string                                                       `json:"updateTime"`        /*  更新时间  */
	BindingAlias      []*string                                                     `json:"bindingAlias"`      /*  版本关联的别名列表  */
	FunctionsSnapshot *CfListFunctionVersionsReturnObjDataFunctionsSnapshotResponse `json:"functionsSnapshot"` /*  版本对应的函数快照信息  */
}

type CfListFunctionVersionsReturnObjPaginationResponse struct {
	PageIndex *int32 `json:"pageIndex"` /*  页码  */
	PageSize  *int32 `json:"pageSize"`  /*  每页大小  */
	Total     *int32 `json:"total"`     /*  总记录数  */
}

type CfListFunctionVersionsReturnObjDataFunctionsSnapshotResponse struct {
	Functions *CfListFunctionVersionsReturnObjDataFunctionsSnapshotFunctionsResponse `json:"functions"` /*  函数信息,请见函数配置  */
}

type CfListFunctionVersionsReturnObjDataFunctionsSnapshotFunctionsResponse struct{}
