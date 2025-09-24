package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcListNetworkPathReportApi
/* 获取网络路径分析报告列表
 */type CtvpcListNetworkPathReportApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcListNetworkPathReportApi(client *core.CtyunClient) *CtvpcListNetworkPathReportApi {
	return &CtvpcListNetworkPathReportApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/list-network-path-report",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcListNetworkPathReportApi) Do(ctx context.Context, credential core.Credential, req *CtvpcListNetworkPathReportRequest) (*CtvpcListNetworkPathReportResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("analysisID", req.AnalysisID)
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
	var resp CtvpcListNetworkPathReportResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcListNetworkPathReportRequest struct {
	RegionID   string /*  资源池 ID  */
	AnalysisID string /*  路径分析 ID  */
	PageNumber int32  /*  列表的页码，默认值为 1。  */
	PageSize   int32  /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcListNetworkPathReportResponse struct {
	StatusCode  int32                                          `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                        `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                        `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                        `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   []*CtvpcListNetworkPathReportReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcListNetworkPathReportReturnObjResponse struct {
	Results      []*CtvpcListNetworkPathReportReturnObjResultsResponse `json:"results"`      /*  网卡列表  */
	TotalCount   int32                                                 `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                 `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                 `json:"totalPage"`    /*  总页数  */
}

type CtvpcListNetworkPathReportReturnObjResultsResponse struct {
	NetworkPathReportID *string `json:"networkPathReportID,omitempty"` /*  网络路径分析报告 ID  */
	CreatedAt           *string `json:"createdAt,omitempty"`           /*  创建时间  */
	UpdatedAt           *string `json:"updatedAt,omitempty"`           /*  更新时间  */
	NetworkPathID       *string `json:"networkPathID,omitempty"`       /*  网络路径 ID  */
	AnalysisID          *string `json:"analysisID,omitempty"`          /*  路径分析 ID  */
	NodeID              *string `json:"nodeID,omitempty"`              /*  节点 ID  */
	Cidr                *string `json:"cidr,omitempty"`                /*  网段  */
	NodeType            *string `json:"nodeType,omitempty"`            /*  节点类型  */
	Ip                  *string `json:"ip,omitempty"`                  /*  ip  */
	Mac                 *string `json:"mac,omitempty"`                 /*  物理地址  */
	ErrCode             *string `json:"errCode,omitempty"`             /*  错误码  */
	Level               int32   `json:"level"`                         /*  拓扑排序 1 - 5  */
	NodeStatus          int32   `json:"nodeStatus"`                    /*  节点状态 0 不可使用， 1 可以使用  */
}
