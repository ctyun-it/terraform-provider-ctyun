package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcVniaShowNetworkPathApi
/* 获取网络路径详情
 */type CtvpcVniaShowNetworkPathApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcVniaShowNetworkPathApi(client *core.CtyunClient) *CtvpcVniaShowNetworkPathApi {
	return &CtvpcVniaShowNetworkPathApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/vnia/show-network-path",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcVniaShowNetworkPathApi) Do(ctx context.Context, credential core.Credential, req *CtvpcVniaShowNetworkPathRequest) (*CtvpcVniaShowNetworkPathResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("networkPathID", req.NetworkPathID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcVniaShowNetworkPathResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcVniaShowNetworkPathRequest struct {
	RegionID      string /*  资源池 ID  */
	NetworkPathID string /*  网络路径 ID  */
}

type CtvpcVniaShowNetworkPathResponse struct {
	StatusCode  int32                                      `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message     *string                                    `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description *string                                    `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   *string                                    `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtvpcVniaShowNetworkPathReturnObjResponse `json:"returnObj"`             /*  接口业务数据  */
}

type CtvpcVniaShowNetworkPathReturnObjResponse struct {
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
