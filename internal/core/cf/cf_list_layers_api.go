package cf

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CfListLayersApi
/* 查询层列表 */
type CfListLayersApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCfListLayersApi(client *core.CtyunClient) *CfListLayersApi {
	return &CfListLayersApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/openapi/v1/layers",
			ContentType:  "application/json",
		},
	}
}

func (a *CfListLayersApi) Do(ctx context.Context, credential core.Credential, req *CfListLayersRequest) (*CfListLayersResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
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
	if req.Scope != nil && *req.Scope != "" {
		ctReq.AddParam("scope", *req.Scope)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CfListLayersResponse
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CfListLayersRequest struct {
	RegionId  string  `json:"regionId"`            /*  资源池id  */
	PageIndex *int32  `json:"pageIndex,omitempty"` /*  页码，默认1  */
	PageSize  *int32  `json:"pageSize,omitempty"`  /*  分页大小，默认10  */
	Search    *string `json:"search,omitempty"`    /*  模糊查询的关键字，默认为空  */
	OrderBy   *string `json:"orderBy,omitempty"`   /*  排序字段名，需要排序的字段名，默认按创建时间  */
	Order     *string `json:"order,omitempty"`     /*  ASC升序，DESC降序，默认降序  */
	Scope     *string `json:"scope,omitempty"`     /*  all 代表返回所有层（包括公共层），custom 代表仅返回自定义层，默认返回所有层  */
}

type CfListLayersResponse struct {
	StatusCode *int32                         `json:"statusCode"` /*  状态码,0表示成功，非0表示不成功  */
	Error      *string                        `json:"error"`      /*  错误码  */
	Message    *string                        `json:"message"`    /*  信息  */
	ReturnObj  *CfListLayersReturnObjResponse `json:"returnObj"`  /*  返回实体  */
}

type CfListLayersReturnObjResponse struct {
	LayerName         *string                            `json:"layerName"`         /*  层名  */
	Acl               *int32                             `json:"acl"`               /*  0 表示私有层，1 表示官方层  */
	Region            *string                            `json:"region"`            /*  资源池编号  */
	Version           *int32                             `json:"version"`           /*  版本  */
	Description       *string                            `json:"description"`       /*  版本描述信息  */
	CompatibleRuntime []*string                          `json:"compatibleRuntime"` /*  版本运行时环境列表  */
	Ctrn              *string                            `json:"ctrn"`              /*  版本唯一标识  */
	Code              *CfListLayersReturnObjCodeResponse `json:"code"`              /*  版本代码配置  */
	Codesize          *int32                             `json:"codesize"`          /*  代码大小  */
	CodeChecksum      *string                            `json:"codeChecksum"`      /*  代码校验码  */
	CreateTime        *string                            `json:"createTime"`        /*  版本创建时间  */
	TenantId          *int32                             `json:"tenantId"`          /*  租户 ID  */
}

type CfListLayersReturnObjCodeResponse struct {
	OssBucketName *string `json:"ossBucketName"` /*  oss的bucket  */
	OssObjectName *string `json:"ossObjectName"` /*  oss的name  */
}
