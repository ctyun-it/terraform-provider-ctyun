package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtelbGwlbListApi
/* 查看网关负载均衡列表
 */type CtelbGwlbListApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbGwlbListApi(client *core.CtyunClient) *CtelbGwlbListApi {
	return &CtelbGwlbListApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/gwlb/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbGwlbListApi) Do(ctx context.Context, credential core.Credential, req *CtelbGwlbListRequest) (*CtelbGwlbListResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.GwLbID != "" {
		ctReq.AddParam("gwLbID", req.GwLbID)
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
	var resp CtelbGwlbListResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbGwlbListRequest struct {
	RegionID   string /*  资源池 ID  */
	ProjectID  string /*  企业项目ID，默认"0"  */
	GwLbID     string /*  网关负载均衡ID  */
	PageNumber int32  /*  列表的页码，默认值为 1。  */
	PageSize   int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtelbGwlbListResponse struct {
	StatusCode  int32                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbGwlbListReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtelbGwlbListReturnObjResponse struct {
	Results      []*CtelbGwlbListReturnObjResultsResponse `json:"results"`                /*  接口业务数据  */
	TotalCount   int32                                    `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                    `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                    `json:"totalPage,omitempty"`    /*  总页数  */
}

type CtelbGwlbListReturnObjResultsResponse struct {
	GwLbID           string `json:"gwLbID,omitempty"`           /*  网关负载均衡 ID  */
	Name             string `json:"name,omitempty"`             /*  名字  */
	Description      string `json:"description,omitempty"`      /*  描述  */
	VpcID            string `json:"vpcID,omitempty"`            /*  虚拟私有云 ID  */
	SubnetID         string `json:"subnetID,omitempty"`         /*  子网 ID  */
	PortID           string `json:"portID,omitempty"`           /*  网卡 ID  */
	Ipv6Enabled      *bool  `json:"ipv6Enabled"`                /*  是否开启 ipv6  */
	PrivateIpAddress string `json:"privateIpAddress,omitempty"` /*  私有 IP 地址  */
	Ipv6Address      string `json:"ipv6Address,omitempty"`      /*  ipv6 地址  */
	SlaName          string `json:"slaName,omitempty"`          /*  规格  */
	DeleteProtection *bool  `json:"deleteProtection"`           /*  是否开启删除保护  */
	CreatedAt        string `json:"createdAt,omitempty"`        /*  创建时间  */
	UpdatedAt        string `json:"updatedAt,omitempty"`        /*  更新时间  */
}
