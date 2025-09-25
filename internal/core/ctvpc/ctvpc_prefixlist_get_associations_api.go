package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcPrefixlistGetAssociationsApi
/* 查看 prefixlist 绑定资源
 */type CtvpcPrefixlistGetAssociationsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcPrefixlistGetAssociationsApi(client *core.CtyunClient) *CtvpcPrefixlistGetAssociationsApi {
	return &CtvpcPrefixlistGetAssociationsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/prefixlist/get_associations",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcPrefixlistGetAssociationsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcPrefixlistGetAssociationsRequest) (*CtvpcPrefixlistGetAssociationsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("prefixListID", req.PrefixListID)
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
	var resp CtvpcPrefixlistGetAssociationsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcPrefixlistGetAssociationsRequest struct {
	RegionID     string /*  区域 id  */
	PrefixListID string /*  前缀列表的ID  */
	PageNumber   int32  /*  列表的页码，默认值为 1。  */
	PageNo       int32  /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize     int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcPrefixlistGetAssociationsResponse struct {
	StatusCode   int32                                              `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                            `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                            `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                            `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcPrefixlistGetAssociationsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	TotalCount   int32                                              `json:"totalCount"`            /*  总条数  */
	TotalPage    int32                                              `json:"totalPage"`             /*  总页数  */
	CurrentCount int32                                              `json:"currentCount"`          /*  当前页总数  */
}

type CtvpcPrefixlistGetAssociationsReturnObjResponse struct {
	ResourceID   *string `json:"resourceID,omitempty"`   /*  资源ID  */
	ResourceName *string `json:"resourceName,omitempty"` /*  资源名称  */
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型  */
}
