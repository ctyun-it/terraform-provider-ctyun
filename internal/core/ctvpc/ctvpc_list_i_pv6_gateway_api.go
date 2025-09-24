package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListIPv6GatewayApi
/* 调用此接口可查询 IPv6 网关列表。
 */type CtvpcListIPv6GatewayApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListIPv6GatewayApi(client *core.CtyunClient) *CtvpcListIPv6GatewayApi {
	return &CtvpcListIPv6GatewayApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/list-ipv6-gateway",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListIPv6GatewayApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListIPv6GatewayRequest) (*CtvpcListIPv6GatewayResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
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
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListIPv6GatewayResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListIPv6GatewayRequest struct {
	ProjectID  *string /*  企业项目 ID，默认为0  */
	PageNumber int32   /*  列表的页码，默认值为 1。  */
	PageNo     int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
	RegionID   string  /*  资源池ID  */
}

type CtvpcListIPv6GatewayResponse struct {
	StatusCode   int32                                    `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                  `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                  `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                  `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListIPv6GatewayReturnObjResponse `json:"returnObj"`             /*  object  */
	TotalCount   int32                                    `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                                    `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                                    `json:"totalPage"`             /*  总页数  */
	Error        *string                                  `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListIPv6GatewayReturnObjResponse struct {
	Name          *string `json:"name,omitempty"`          /*  ipv6网关名称  */
	Status        *string `json:"status,omitempty"`        /*  仅有一个状态：ACTIVE（正常）  */
	ProjectIdEcs  *string `json:"projectIdEcs,omitempty"`  /*  企业项目  */
	VpcID         *string `json:"vpcID,omitempty"`         /*  vpcID值  */
	VpcCidr       *string `json:"vpcCidr,omitempty"`       /*  无类别域间路由  */
	VpcName       *string `json:"vpcName,omitempty"`       /*  vpc名称  */
	Ipv6GatewayID *string `json:"ipv6GatewayID,omitempty"` /*  ipv6网关惟一值  */
	CreationTime  *string `json:"creationTime,omitempty"`  /*  创建时间  */
}
