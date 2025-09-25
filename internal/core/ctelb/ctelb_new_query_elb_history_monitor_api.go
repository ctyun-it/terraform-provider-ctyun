package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbNewQueryElbHistoryMonitorApi
/* 查看负载均衡历史监控。
 */type CtelbNewQueryElbHistoryMonitorApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbNewQueryElbHistoryMonitorApi(client *core.CtyunClient) *CtelbNewQueryElbHistoryMonitorApi {
	return &CtelbNewQueryElbHistoryMonitorApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/new-query-history-monitor",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbNewQueryElbHistoryMonitorApi) Do(ctx context.Context, credential core.Credential, req *CtelbNewQueryElbHistoryMonitorRequest) (*CtelbNewQueryElbHistoryMonitorResponse, error) {
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
	var resp CtelbNewQueryElbHistoryMonitorResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbNewQueryElbHistoryMonitorRequest struct {
	RegionID    string   `json:"regionID,omitempty"`   /*  资源池 ID  */
	DeviceIDs   []string `json:"deviceIDs"`            /*  负载均衡 ID 列表  */
	MetricNames []string `json:"metricNames"`          /*  监控指标  */
	StartTime   string   `json:"startTime,omitempty"`  /*  开始时间，YYYY-mmm-dd HH:MM:SS（只允许dd和HH中间有一个空格）  */
	EndTime     string   `json:"endTime,omitempty"`    /*  开始时间，YYYY-mmm-dd HH:MM:SS（只允许dd和HH中间有一个空格）  */
	Period      int32    `json:"period,omitempty"`     /*  可选参数，聚合周期，单位：秒，默认60，例14400  */
	PageNumber  int32    `json:"pageNumber,omitempty"` /*  列表的页码，默认值为 1  */
	PageNo      int32    `json:"pageNo,omitempty"`     /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize    int32    `json:"pageSize,omitempty"`   /*  每页数据量大小，取值 1-50  */
}

type CtelbNewQueryElbHistoryMonitorResponse struct {
	StatusCode  int32                                            `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                           `json:"message,omitempty"`     /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description string                                           `json:"description,omitempty"` /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode   string                                           `json:"errorCode,omitempty"`   /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj   *CtelbNewQueryElbHistoryMonitorReturnObjResponse `json:"returnObj"`             /*  返回结果  */
	Error       string                                           `json:"error,omitempty"`       /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbNewQueryElbHistoryMonitorReturnObjResponse struct {
	Monitors     []*CtelbNewQueryElbHistoryMonitorReturnObjMonitorsResponse `json:"monitors"`               /*  监控数据  */
	TotalCount   int32                                                      `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                                      `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                                      `json:"totalPage,omitempty"`    /*  总页数  */
}

type CtelbNewQueryElbHistoryMonitorReturnObjMonitorsResponse struct {
	LastUpdated       string                                                                      `json:"lastUpdated,omitempty"` /*  最近更新时间  */
	RegionID          string                                                                      `json:"regionID,omitempty"`    /*  资源池 ID  */
	DeviceID          string                                                                      `json:"deviceID,omitempty"`    /*  弹性公网 IP  */
	ItemAggregateList []*CtelbNewQueryElbHistoryMonitorReturnObjMonitorsItemAggregateListResponse `json:"itemAggregateList"`     /*  监控项值列表，这里 itemAggregateList 是一个数组，每个数组是一个字典，每个字典的对应一个数组，每个字典的 key 对应 metricNames 中的一个元素.  */
}

type CtelbNewQueryElbHistoryMonitorReturnObjMonitorsItemAggregateListResponse struct{}
