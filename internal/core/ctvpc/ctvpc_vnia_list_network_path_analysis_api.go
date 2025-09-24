package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtvpcVniaListNetworkPathAnalysisApi
/* 获取网络路径分析列表
 */type CtvpcVniaListNetworkPathAnalysisApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaListNetworkPathAnalysisApi(client *core.CtyunClient) *CtvpcVniaListNetworkPathAnalysisApi {
	return &CtvpcVniaListNetworkPathAnalysisApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/list-network-path-analysis",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaListNetworkPathAnalysisApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaListNetworkPathAnalysisRequest) (*CtvpcVniaListNetworkPathAnalysisResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.NetworkPathID != nil {
		ctReq.AddParam("networkPathID", *req.NetworkPathID)
	}
	if req.AnalysisID != nil {
		ctReq.AddParam("analysisID", *req.AnalysisID)
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
	var resp CtvpcVniaListNetworkPathAnalysisResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaListNetworkPathAnalysisRequest struct {
	RegionID      string  /*  资源池 ID  */
	NetworkPathID *string /*  网络路径 ID  */
	AnalysisID    *string /*  路径分析 ID  */
	PageNumber    int32   /*  列表的页码，默认值为 1。  */
	PageSize      int32   /*  分页查询时每页的行数，最大值为 50，默认值为 10。  */
}

type CtvpcVniaListNetworkPathAnalysisResponse struct {
	StatusCode  int32                                                `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                                              `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                                              `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                                              `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   []*CtvpcVniaListNetworkPathAnalysisReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaListNetworkPathAnalysisReturnObjResponse struct {
	Results      []*CtvpcVniaListNetworkPathAnalysisReturnObjResultsResponse `json:"results"`      /*  网卡列表  */
	TotalCount   int32                                                       `json:"totalCount"`   /*  列表条目数  */
	CurrentCount int32                                                       `json:"currentCount"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                       `json:"totalPage"`    /*  总页数  */
}

type CtvpcVniaListNetworkPathAnalysisReturnObjResultsResponse struct {
	AnalysisID     *string `json:"analysisID,omitempty"`     /*  路径分析 ID  */
	CreatedAt      *string `json:"createdAt,omitempty"`      /*  创建时间  */
	UpdatedAt      *string `json:"updatedAt,omitempty"`      /*  更新时间  */
	NetworkPathID  *string `json:"networkPathID,omitempty"`  /*  网络路径 ID  */
	AnalysisStatus *string `json:"analysisStatus,omitempty"` /*  分析状态  */
	ReachableFlag  *string `json:"reachableFlag,omitempty"`  /*  可访问性标记：unreachable / reachable / unknown  */
	ErrMsg         *string `json:"errMsg,omitempty"`         /*  错误信息  */
}
