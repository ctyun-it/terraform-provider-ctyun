package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcVniaListNetworkPathApi
/* 获取网络路径列表
 */type CtvpcVniaListNetworkPathApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaListNetworkPathApi(client *core.CtyunClient) *CtvpcVniaListNetworkPathApi {
	return &CtvpcVniaListNetworkPathApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/list-network-path",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaListNetworkPathApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaListNetworkPathRequest) (*CtvpcVniaListNetworkPathResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.NetworkPathID != nil {
		ctReq.AddParam("networkPathID", *req.NetworkPathID)
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
	var resp CtvpcVniaListNetworkPathResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaListNetworkPathRequest struct {
	RegionID      string  /*  资源池 ID  */
	NetworkPathID *string /*  网络路径 ID  */
	PageNumber    int32   /*  列表的页码，默认值为 1。  */
	PageSize      int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcVniaListNetworkPathResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   []*CtvpcVniaListNetworkPathReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaListNetworkPathReturnObjResponse struct {
	Results      []*CtvpcVniaListNetworkPathReturnObjResultsResponse `json:"results"`      /*  网卡列表  */
	TotalCount   int32                                               `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                               `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                               `json:"totalPage"`    /*  总页数  */
}

type CtvpcVniaListNetworkPathReturnObjResultsResponse struct {
	NetworkPathID  *string `json:"networkPathID,omitempty"`  /*  网络路径 ID  */
	Name           *string `json:"name,omitempty"`           /*  路径分析名字  */
	SourceID       *string `json:"sourceID,omitempty"`       /*  源设备  */
	SourceType     *string `json:"sourceType,omitempty"`     /*  源类型，目前仅支持 ecs / internet / subnet  */
	SourcePort     int32   `json:"sourcePort"`               /*  源端口, 1 - 65535  */
	TargetType     *string `json:"targetType,omitempty"`     /*  目标类型，目前仅支持 ecs / internet / subnet / elb  */
	TargetPort     int32   `json:"targetPort"`               /*  目标端口, 1 - 65535  */
	TargetID       *string `json:"targetID,omitempty"`       /*  目标设备  */
	SourceIP       *string `json:"sourceIP,omitempty"`       /*  源 IP  */
	TargetIP       *string `json:"targetIP,omitempty"`       /*  目的 IP  */
	Protocol       *string `json:"protocol,omitempty"`       /*  协议，仅支持 ICMP / TCP / UDP  */
	AnalysisStatus *string `json:"analysisStatus,omitempty"` /*  分析状态  */
	ReachableFlag  *string `json:"reachableFlag,omitempty"`  /*  可访问性标记：unreachable / reachable / unknown  */
	CreatedAt      *string `json:"createdAt,omitempty"`      /*  创建时间  */
	UpdatedAt      *string `json:"updatedAt,omitempty"`      /*  更新时间  */
}
