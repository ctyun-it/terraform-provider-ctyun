package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListVpcPeerRequestsApi
/* 获取待处理的对等连接请求
 */type CtvpcListVpcPeerRequestsApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListVpcPeerRequestsApi(client *core.CtyunClient) *CtvpcListVpcPeerRequestsApi {
	return &CtvpcListVpcPeerRequestsApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/vpcpeer/requests",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListVpcPeerRequestsApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListVpcPeerRequestsRequest) (*CtvpcListVpcPeerRequestsResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Page != 0 {
		ctReq.AddParam("page", strconv.FormatInt(int64(req.Page), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListVpcPeerRequestsResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListVpcPeerRequestsRequest struct {
	PageSize int32  /*  当前页数据条数  */
	Page     int32  /*  当前页  */
	PageNo   int32  /*  列表的页码，默认值为 1, 推荐使用该字段, page 后续会废弃  */
	RegionID string /*  区域id  */
}

type CtvpcListVpcPeerRequestsResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcListVpcPeerRequestsReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                                    `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListVpcPeerRequestsReturnObjResponse struct {
	VpcPeerRequests []*CtvpcListVpcPeerRequestsReturnObjVpcPeerRequestsResponse `json:"vpcPeerRequests"` /*  待处理的对等链接建立请求  */
	TotalCount      int32                                                       `json:"totalCount"`      /*  列表条目数  */
	CurrentCount    int32                                                       `json:"currentCount"`    /*  分页查询时每页的行数。  */
	TotalPage       int32                                                       `json:"totalPage"`       /*  总页数  */
}

type CtvpcListVpcPeerRequestsReturnObjVpcPeerRequestsResponse struct {
	RequestVpcID   *string `json:"requestVpcID,omitempty"`   /*  本端的vpcId值  */
	RequestVpcName *string `json:"requestVpcName,omitempty"` /*  本端vpc的名称  */
	RequestVpcCidr *string `json:"requestVpcCidr,omitempty"` /*  本端vpc的网段  */
	AcceptVpcID    *string `json:"acceptVpcID,omitempty"`    /*  对端的vpcId值  */
	AcceptVpcName  *string `json:"acceptVpcName,omitempty"`  /*  对端vpc的名称，  */
	AcceptVpcCidr  *string `json:"acceptVpcCidr,omitempty"`  /*  对端vpc的网段  */
	InstanceID     *string `json:"instanceID,omitempty"`     /*  对等连接 ID  */
	Status         *string `json:"status,omitempty"`         /*  pending(待同意建立对等链接)  */
}
