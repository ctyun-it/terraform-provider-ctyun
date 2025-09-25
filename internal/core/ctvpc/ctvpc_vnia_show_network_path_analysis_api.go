package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaShowNetworkPathAnalysisApi
/* 获取网络路径分析详情
 */type CtvpcVniaShowNetworkPathAnalysisApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaShowNetworkPathAnalysisApi(client *core.CtyunClient) *CtvpcVniaShowNetworkPathAnalysisApi {
	return &CtvpcVniaShowNetworkPathAnalysisApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/show-network-path-analysis",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaShowNetworkPathAnalysisApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaShowNetworkPathAnalysisRequest) (*CtvpcVniaShowNetworkPathAnalysisResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("analysisID", req.AnalysisID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVniaShowNetworkPathAnalysisResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaShowNetworkPathAnalysisRequest struct {
	RegionID   string /*  资源池 ID  */
	AnalysisID string /*  路径分析 ID  */
}

type CtvpcVniaShowNetworkPathAnalysisResponse struct {
	StatusCode  int32                                              `json:"statusCode"`            /*  返回状态码（800 为成功，900 为失败）  */
	Message     *string                                            `json:"message,omitempty"`     /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为 success, 英文  */
	Description *string                                            `json:"description,omitempty"` /*  statusCode 为 900 时的错误信息; statusCode 为 800 时为成功, 中文  */
	ErrorCode   *string                                            `json:"errorCode,omitempty"`   /*  statusCode 为 900 时为业务细分错误码，三段式：product.module.code; statusCode 为 800 时为 SUCCESS  */
	ReturnObj   *CtvpcVniaShowNetworkPathAnalysisReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaShowNetworkPathAnalysisReturnObjResponse struct {
	AnalysisID     *string `json:"analysisID,omitempty"`     /*  路径分析 ID  */
	CreatedAt      *string `json:"createdAt,omitempty"`      /*  创建时间  */
	UpdatedAt      *string `json:"updatedAt,omitempty"`      /*  更新时间  */
	NetworkPathID  *string `json:"networkPathID,omitempty"`  /*  网络路径 ID  */
	AnalysisStatus *string `json:"analysisStatus,omitempty"` /*  分析状态  */
	ReachableFlag  *string `json:"reachableFlag,omitempty"`  /*  可访问性标记：unreachable / reachable / unknown  */
	ErrMsg         *string `json:"errMsg,omitempty"`         /*  错误信息  */
}
