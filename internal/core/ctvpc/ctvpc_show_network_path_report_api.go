package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcShowNetworkPathReportApi
/* 获取网络路径分析报告详情
 */type CtvpcShowNetworkPathReportApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcShowNetworkPathReportApi(client *core.CtyunClient) *CtvpcShowNetworkPathReportApi {
	return &CtvpcShowNetworkPathReportApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/show-network-path-report",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcShowNetworkPathReportApi) Do(ctx context.Context, credential core.Credential, req *CtvpcShowNetworkPathReportRequest) (*CtvpcShowNetworkPathReportResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("networkPathReportID", req.NetworkPathReportID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcShowNetworkPathReportResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcShowNetworkPathReportRequest struct {
	RegionID            string /*  资源池 ID  */
	NetworkPathReportID string /*  路径分析报告 ID  */
}

type CtvpcShowNetworkPathReportResponse struct {
	StatusCode  int32                                        `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                      `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                      `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                      `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcShowNetworkPathReportReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcShowNetworkPathReportReturnObjResponse struct {
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
