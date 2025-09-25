package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcQueryResourcesByLabelApi
/* 获取绑定标签的资源列表
 */type CtvpcQueryResourcesByLabelApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryResourcesByLabelApi(client *core.CtyunClient) *CtvpcQueryResourcesByLabelApi {
	return &CtvpcQueryResourcesByLabelApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/labels/query_resources_by_label",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryResourcesByLabelApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryResourcesByLabelRequest) (*CtvpcQueryResourcesByLabelResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.LabelID != nil {
		ctReq.AddParam("labelID", *req.LabelID)
	}
	if req.LabelKey != nil {
		ctReq.AddParam("labelKey", *req.LabelKey)
	}
	if req.LabelValue != nil {
		ctReq.AddParam("labelValue", *req.LabelValue)
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcQueryResourcesByLabelResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryResourcesByLabelRequest struct {
	RegionID   string  /*  区域ID  */
	LabelID    *string /*  标签ID，label的三个参数至少选填一个  */
	LabelKey   *string /*  标签 key  */
	LabelValue *string /*  标签 取值  */
	PageNumber int32   /*  列表的页码，默认值为 1  */
	PageSize   int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10  */
}

type CtvpcQueryResourcesByLabelResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcQueryResourcesByLabelReturnObjResponse `json:"returnObj"`             /*  返回结果  */
}

type CtvpcQueryResourcesByLabelReturnObjResponse struct {
	Results      []*CtvpcQueryResourcesByLabelReturnObjResultsResponse `json:"results"`      /*  绑定的标签列表  */
	TotalCount   int32                                                 `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                 `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                 `json:"totalPage"`    /*  总页数  */
}

type CtvpcQueryResourcesByLabelReturnObjResultsResponse struct {
	ResourceID   *string `json:"resourceID,omitempty"`   /*  资源ID  */
	ResourceType *string `json:"resourceType,omitempty"` /*  资源类型，vpc / subnet / acl / security_group / route_table / havip / port  / multicast_domain / vpc_peer / vpce_endpoint / vpce_endpoint_service / ipv6_gateway / elb / private_nat / nat / listener，network表示eip、bandwidth、ipv6_bandwidth类型  */
}
