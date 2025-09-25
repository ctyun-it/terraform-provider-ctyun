package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListSubnetUsedIPsApi
/* 查看某个子网已使用IP
 */type CtvpcListSubnetUsedIPsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListSubnetUsedIPsApi(client *core.CtyunClient) *CtvpcListSubnetUsedIPsApi {
	return &CtvpcListSubnetUsedIPsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/list-used-ips",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListSubnetUsedIPsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListSubnetUsedIPsRequest) (*CtvpcListSubnetUsedIPsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("subnetID", req.SubnetID)
	if req.Ip != nil {
		ctReq.AddParam("ip", *req.Ip)
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
	var resp CtvpcListSubnetUsedIPsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListSubnetUsedIPsRequest struct {
	RegionID   string  /*  资源池 ID  */
	SubnetID   string  /*  子网 ID  */
	Ip         *string /*  子网内的 IP 地址  */
	PageNumber int32   /*  列表的页码，默认值为 1。  */
	PageNo     int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListSubnetUsedIPsResponse struct {
	StatusCode  int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListSubnetUsedIPsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListSubnetUsedIPsReturnObjResponse struct {
	UsedIPs      []*CtvpcListSubnetUsedIPsReturnObjUsedIPsResponse `json:"usedIPs"`      /*  已使用的 IP 数组  */
	TotalCount   int32                                             `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                             `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                             `json:"totalPage"`    /*  总页数  */
}

type CtvpcListSubnetUsedIPsReturnObjUsedIPsResponse struct {
	Ipv4Address *string `json:"ipv4Address,omitempty"` /*  ipv4 地址  */
	Ipv6Address *string `json:"ipv6Address,omitempty"` /*  ipv6 地址  */
	UseDesc     *string `json:"useDesc,omitempty"`     /*  用途中文描述:云主机, 裸金属, 高可用虚 IP, SNAT, 负载均衡, 预占内网 IP, 内网网关接口, system  */
}
