package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListIPv6Api
/* 调用此接口可查询 IPv6 列表。
 */type CtvpcListIPv6Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListIPv6Api(client *core.CtyunClient) *CtvpcListIPv6Api {
	return &CtvpcListIPv6Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ipv6/ipv6-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListIPv6Api) Do(ctx context.Context, credential core.Credential, req *CtvpcListIPv6Request) (*CtvpcListIPv6Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.SubnetID != nil {
		ctReq.AddParam("subnetID", *req.SubnetID)
	}
	if req.IpAddress != nil {
		ctReq.AddParam("ipAddress", *req.IpAddress)
	}
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
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
	var resp CtvpcListIPv6Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListIPv6Request struct {
	RegionID  string  /*  资源池ID  */
	VpcID     *string /*  vpc id  */
	SubnetID  *string /*  子网id  */
	IpAddress *string /*  ipv6地址  */
	Page      int32   /*  分页参数  */
	PageNo    int32   /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize  int32   /*  每页数据量大小，取值1-50  */
}

type CtvpcListIPv6Response struct {
	StatusCode   int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListIPv6ReturnObjResponse `json:"returnObj"`             /*  业务数据  */
	TotalCount   int32                             `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                             `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                             `json:"totalPage"`             /*  总页数  */
	Error        *string                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListIPv6ReturnObjResponse struct{}
