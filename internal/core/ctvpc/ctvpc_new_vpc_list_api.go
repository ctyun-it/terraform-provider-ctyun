package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcNewVpcListApi
/* 查询用户专有网络列表
 */type CtvpcNewVpcListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcNewVpcListApi(client *core.CtyunClient) *CtvpcNewVpcListApi {
	return &CtvpcNewVpcListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vpc/new-list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcNewVpcListApi) Do(ctx context.Context, credential core.Credential, req *CtvpcNewVpcListRequest) (*CtvpcNewVpcListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VpcID != nil {
		ctReq.AddParam("vpcID", *req.VpcID)
	}
	if req.VpcName != nil {
		ctReq.AddParam("vpcName", *req.VpcName)
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
	if req.ProjectID != nil {
		ctReq.AddParam("projectID", *req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcNewVpcListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcNewVpcListRequest struct {
	RegionID   string  /*  资源池 ID  */
	VpcID      *string /*  多个 VPC 的 ID 之间用半角逗号（,）隔开。  */
	VpcName    *string /*  vpc 名字  */
	PageNumber int32   /*  列表的页码，默认值为 1。  */
	PageNo     int32   /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
	ProjectID  *string /*  企业项目 ID，默认为0  */
}

type CtvpcNewVpcListResponse struct {
	StatusCode  int32                             `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcNewVpcListReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
	Error       *string                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcNewVpcListReturnObjResponse struct {
	Vpcs         []*CtvpcNewVpcListReturnObjVpcsResponse `json:"vpcs"`         /*  vpc 组  */
	TotalCount   int32                                   `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                   `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                   `json:"totalPage"`    /*  总页数  */
}

type CtvpcNewVpcListReturnObjVpcsResponse struct {
	VpcID          *string   `json:"vpcID,omitempty"`       /*  vpc 示例 ID  */
	Name           *string   `json:"name,omitempty"`        /*  名称  */
	Description    *string   `json:"description,omitempty"` /*  描述  */
	CIDR           *string   `json:"CIDR,omitempty"`        /*  CIDR  */
	Ipv6Enabled    *bool     `json:"ipv6Enabled"`           /*  是否开启 ipv6  */
	EnableIpv6     *bool     `json:"enableIpv6"`            /*  是否开启 ipv6  */
	Ipv6CIDRS      []*string `json:"ipv6CIDRS"`             /*  ipv6CIDRS  */
	SubnetIDs      []*string `json:"subnetIDs"`             /*  子网 id 列表  */
	NatGatewayIDs  []*string `json:"natGatewayIDs"`         /*  网关 id 列表  */
	SecondaryCIDRS []*string `json:"secondaryCIDRS"`        /*  附加网段  */
	ProjectID      *string   `json:"projectID,omitempty"`   /*  企业项目 ID，默认为0  */
}
