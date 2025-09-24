package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListVpcPeerConnectionApi
/* 查询对等连接列表
 */type CtvpcListVpcPeerConnectionApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListVpcPeerConnectionApi(client *core.CtyunClient) *CtvpcListVpcPeerConnectionApi {
	return &CtvpcListVpcPeerConnectionApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/list-vpc-peer-connection",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListVpcPeerConnectionApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListVpcPeerConnectionRequest) (*CtvpcListVpcPeerConnectionResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PageNumber != 0 {
		ctReq.AddParam("pageNumber", strconv.FormatInt(int64(req.PageNumber), 10))
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	ctReq.AddParam("regionID", req.RegionID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcListVpcPeerConnectionResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListVpcPeerConnectionRequest struct {
	PageSize   int32  /*  当前页数据条数  */
	PageNumber int32  /*  当前页  */
	PageNo     int32  /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	RegionID   string /*  区域id  */
}

type CtvpcListVpcPeerConnectionResponse struct {
	StatusCode   int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcListVpcPeerConnectionReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	TotalCount   int32                                          `json:"totalCount"`            /*  总条数  */
	TotalPage    int32                                          `json:"totalPage"`             /*  总页数  */
	CurrentCount int32                                          `json:"currentCount"`          /*  总页数  */
	Error        *string                                        `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcListVpcPeerConnectionReturnObjResponse struct {
	RequestVpcID   *string `json:"requestVpcID,omitempty"`   /*  本端vpcId值  */
	RequestVpcName *string `json:"requestVpcName,omitempty"` /*  本端vpc的名称  */
	RequestVpcCidr *string `json:"requestVpcCidr,omitempty"` /*  本端vpc的网段  */
	AcceptVpcID    *string `json:"acceptVpcID,omitempty"`    /*  对端的vpcId值  */
	AcceptVpcName  *string `json:"acceptVpcName,omitempty"`  /*  对端vpc的名称，  */
	AcceptVpcCidr  *string `json:"acceptVpcCidr,omitempty"`  /*  对端vpc的网段  */
	AcceptEmail    *string `json:"acceptEmail,omitempty"`    /*  对端vpc账户的邮箱  */
	UserType       *string `json:"userType,omitempty"`       /*  对等连接类型：current(同一个租户) / other(不同租户)  */
	Name           *string `json:"name,omitempty"`           /*  对等连接名称  */
	InstanceID     *string `json:"instanceID,omitempty"`     /*  对等连接 ID  */
	Status         *string `json:"status,omitempty"`         /*  对等连接状态：agree(同意建立对等连接) / pending(待同意建立对等链接)  */
}
