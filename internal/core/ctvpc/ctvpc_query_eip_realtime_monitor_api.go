package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtvpcQueryEipRealtimeMonitorApi
/* 查看弹性 IP 实时监控数据。
 */type CtvpcQueryEipRealtimeMonitorApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtvpcQueryEipRealtimeMonitorApi(client *core.CtyunClient) *CtvpcQueryEipRealtimeMonitorApi {
	return &CtvpcQueryEipRealtimeMonitorApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/eip/query-realtime-monitor",
			ContentType:  "application/json",
		},
	}
}

func (a *CtvpcQueryEipRealtimeMonitorApi) Do(ctx context.Context, credential core.Credential, req *CtvpcQueryEipRealtimeMonitorRequest) (*CtvpcQueryEipRealtimeMonitorResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtvpcQueryEipRealtimeMonitorResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtvpcQueryEipRealtimeMonitorRequest struct {
	RegionID   string    `json:"regionID,omitempty"` /*  资源池 ID  */
	DeviceIDs  []*string `json:"deviceIDs"`          /*  弹性 IP 地址列表，仅支持 IP 地址  */
	PageNumber int32     `json:"pageNumber"`         /*  列表的页码，默认值为 1  */
	PageNo     int32     `json:"pageNo"`             /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32     `json:"pageSize"`           /*  每页数据量大小，取值 1-50  */
}

type CtvpcQueryEipRealtimeMonitorResponse struct {
	StatusCode   int32                                            `json:"statusCode"`            /*  返回状态码（800为成功，900为失败）  */
	Message      *string                                          `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  *string                                          `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    *string                                          `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtvpcQueryEipRealtimeMonitorReturnObjResponse `json:"returnObj"`             /*  object  */
	TotalCount   int32                                            `json:"totalCount"`            /*  列表条目数  */
	CurrentCount int32                                            `json:"currentCount"`          /*  分页查询时每页的行数。  */
	TotalPage    int32                                            `json:"totalPage"`             /*  总页数  */
	Error        *string                                          `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtvpcQueryEipRealtimeMonitorReturnObjResponse struct {
	LastUpdated *string                                                  `json:"lastUpdated,omitempty"` /*  最近更新时间  */
	RegionID    *string                                                  `json:"regionID,omitempty"`    /*  资源池 ID  */
	DeviceID    *string                                                  `json:"deviceID,omitempty"`    /*  弹性公网 IP  */
	ItemList    []*CtvpcQueryEipRealtimeMonitorReturnObjItemListResponse `json:"itemList"`              /*  监控项值列表  */
}

type CtvpcQueryEipRealtimeMonitorReturnObjItemListResponse struct {
	IngressThroughput *string `json:"ingressThroughput,omitempty"` /*  入吞吐量  */
	EgressThroughput  *string `json:"egressThroughput,omitempty"`  /*  出吞吐量  */
	IngressBandwidth  *string `json:"ingressBandwidth,omitempty"`  /*  入带宽峰值  */
	EgressBandwidth   *string `json:"egressBandwidth,omitempty"`   /*  出带宽峰值  */
	SamplingTime      *string `json:"samplingTime,omitempty"`      /*  采样时间  */
}
