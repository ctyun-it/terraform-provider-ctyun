package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListEndpointApi
/* 查看终端节点列表
 */type CtvpcListEndpointApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListEndpointApi(client *core.CtyunClient) *CtvpcListEndpointApi {
	return &CtvpcListEndpointApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpce/list-endpoint",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListEndpointApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListEndpointRequest) (*CtvpcListEndpointResponse, error) {
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
	if req.EndpointName != nil {
		ctReq.AddParam("endpointName", *req.EndpointName)
	}
	if req.QueryContent != nil {
		ctReq.AddParam("queryContent", *req.QueryContent)
	}
	if req.EndpointID != nil {
		ctReq.AddParam("endpointID", *req.EndpointID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListEndpointResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListEndpointRequest struct {
	RegionID     string  /*  资源池ID  */
	Page         int32   /*  分页参数, 默认 1  */
	PageNo       int32   /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	PageSize     int32   /*  每页数据量大小, 默认 10  */
	EndpointName *string /*  终端节点名，该字段为精确匹配, 与 endpointID 是或的条件进行查询  */
	QueryContent *string /*  支持对终端节点名进行模糊匹配, 与 endpointName 和 endpointID 是且条件进行查询  */
	EndpointID   *string /*  终端节点 ID  */
}

type CtvpcListEndpointResponse struct {
	StatusCode   int32                                 `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                               `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                               `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                               `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	TotalCount   int32                                 `json:"totalCount"`            /*  总条数  */
	TotalPage    int32                                 `json:"totalPage"`             /*  总页数  */
	CurrentCount int32                                 `json:"currentCount"`          /*  总页数  */
	ReturnObj    []*CtvpcListEndpointReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error        *string                               `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListEndpointReturnObjResponse struct {
	Endpoints []*CtvpcListEndpointReturnObjEndpointsResponse `json:"endpoints"` /*  终端节点列表数据  */
}

type CtvpcListEndpointReturnObjEndpointsResponse struct {
	ID                *string                                                 `json:"ID,omitempty"`                /*  终端节点ID  */
	EndpointServiceID *string                                                 `json:"endpointServiceID,omitempty"` /*  终端节点服务ID  */
	RawType           *string                                                 `json:"type,omitempty"`              /*  接口还是反向，interface:接口，reverse:反向  */
	Name              *string                                                 `json:"name,omitempty"`              /*  终端节点名称  */
	VpcID             *string                                                 `json:"vpcID,omitempty"`             /*  所属的专有网络id  */
	VpcAddress        *string                                                 `json:"vpcAddress,omitempty"`        /*  私网地址  */
	Whitelist         *string                                                 `json:"whitelist,omitempty"`         /*  白名单  */
	Status            int32                                                   `json:"status"`                      /*  endpoint状态, 1 表示已链接，2 表示未链接  */
	Description       *string                                                 `json:"description,omitempty"`       /*  描述  */
	EndpointObj       *CtvpcListEndpointReturnObjEndpointsEndpointObjResponse `json:"endpointObj"`                 /*  后端节点信息，可能为 null  */
	CreatedTime       *string                                                 `json:"createdTime,omitempty"`       /*  创建时间  */
	UpdatedTime       *string                                                 `json:"updatedTime,omitempty"`       /*  更新时间  */
}

type CtvpcListEndpointReturnObjEndpointsEndpointObjResponse struct {
	SubnetID  *string   `json:"subnetID,omitempty"` /*  子网id  */
	PortID    *string   `json:"portID,omitempty"`   /*  端口id  */
	Ip        *string   `json:"ip,omitempty"`       /*  私网地址  */
	EnableDns int32     `json:"enableDns"`          /*  是否开启 dns  */
	DnsNames  []*string `json:"dnsNames"`           /*  dns名称列表  */
}
