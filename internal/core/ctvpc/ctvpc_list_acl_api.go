package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListAclApi
/* 查看 Acl 列表信息
 */type CtvpcListAclApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListAclApi(client *core.CtyunClient) *CtvpcListAclApi {
	return &CtvpcListAclApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/acl/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListAclApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListAclRequest) (*CtvpcListAclResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.AclID != nil {
		ctReq.AddParam("aclID", *req.AclID)
	}
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	if req.Name != nil {
		ctReq.AddParam("name", *req.Name)
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
	var resp CtvpcListAclResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListAclRequest struct {
	RegionID   string  /*  资源池ID  */
	AclID      *string /*  aclID  */
	ProjectID  *string /*  企业项目 ID，默认为"0"  */
	Name       *string /*  acl Name  */
	PageNumber int32   /*  列表的页码，默认值为 1  */
	PageNo     int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10  */
}

type CtvpcListAclResponse struct {
	StatusCode   int32                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    *CtvpcListAclReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	TotalCount   int32                          `json:"totalCount"`            /*  列表条目数。  */
	CurrentCount int32                          `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                          `json:"totalPage"`             /*  分页查询时总页数。  */
}

type CtvpcListAclReturnObjResponse struct {
	Acls []*CtvpcListAclReturnObjAclsResponse `json:"acls"` /*  acls  */
}

type CtvpcListAclReturnObjAclsResponse struct {
	AclID *string `json:"aclID,omitempty"` /*  acl id  */
	Name  *string `json:"name,omitempty"`  /*  acl 名称  */
}
