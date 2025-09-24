package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListEndpointServiceTransitIPApi
/* 终端节点服务中转IP查询接口
 */type CtvpcListEndpointServiceTransitIPApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListEndpointServiceTransitIPApi(client *core.CtyunClient) *CtvpcListEndpointServiceTransitIPApi {
	return &CtvpcListEndpointServiceTransitIPApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/list-endpoint-service-transit-ip",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListEndpointServiceTransitIPApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListEndpointServiceTransitIPRequest) (*CtvpcListEndpointServiceTransitIPResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	ctReq.AddParam("endpointServiceID", req.EndpointServiceID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListEndpointServiceTransitIPResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListEndpointServiceTransitIPRequest struct {
	RegionID          string /*  资源池ID  */
	Page              int32  /*  分页参数, 默认1  */
	PageNo            int32  /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize          int32  /*  每页数据量大小， 默认10  */
	EndpointServiceID string /*  终端节点服务id  */
}

type CtvpcListEndpointServiceTransitIPResponse struct {
	StatusCode   int32                                               `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                             `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                             `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                             `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	TotalCount   int32                                               `json:"totalCount"`            /*  总条数  */
	TotalPage    int32                                               `json:"totalPage"`             /*  总页数  */
	CurrentCount int32                                               `json:"currentCount"`          /*  总页数  */
	Error        *string                                             `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    *CtvpcListEndpointServiceTransitIPReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcListEndpointServiceTransitIPReturnObjResponse struct {
	TransitIPs   []*CtvpcListEndpointServiceTransitIPReturnObjTransitIPsResponse `json:"transitIPs"`             /*  接口业务数据  */
	TotalCount   int32                                                           `json:"totalCount"`             /*  总条数  */
	TotalPage    int32                                                           `json:"totalPage,omitempty"`    /*  总页数  */
	CurrentCount int32                                                           `json:"currentCount,omitempty"` /*  当前条目  */
}

type CtvpcListEndpointServiceTransitIPReturnObjTransitIPsResponse struct {
	SubnetID          *string `json:"subnetID,omitempty"`          /*  子网 ID  */
	EndpointServiceID *string `json:"endpointServiceID,omitempty"` /*  终端服务节点 ID  */
	TransitIP         *string `json:"transitIP,omitempty"`         /*  中转地址  */
	CreatedAt         *string `json:"createdAt,omitempty"`         /*  创建时间  */
	UpdatedAt         *string `json:"updatedAt,omitempty"`         /*  更新时间  */
}
