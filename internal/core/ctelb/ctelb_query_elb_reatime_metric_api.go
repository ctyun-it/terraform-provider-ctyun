package ctelb

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtelbQueryElbReatimeMetricApi
/* 查看负载均衡实时监控。
 */type CtelbQueryElbReatimeMetricApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtelbQueryElbReatimeMetricApi(client *core.CtyunClient) *CtelbQueryElbReatimeMetricApi {
	return &CtelbQueryElbReatimeMetricApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/elb/query-realtime-monitor",
			ContentType:  "application/json",
		},
	}
}

func (a *CtelbQueryElbReatimeMetricApi) Do(ctx context.Context, credential core.Credential, req *CtelbQueryElbReatimeMetricRequest) (*CtelbQueryElbReatimeMetricResponse, error) {
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
	var resp CtelbQueryElbReatimeMetricResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtelbQueryElbReatimeMetricRequest struct {
	RegionID   string   `json:"regionID,omitempty"`   /*  资源池 ID  */
	DeviceIDs  []string `json:"deviceIDs"`            /*  负载均衡 ID 列表  */
	PageNumber int32    `json:"pageNumber,omitempty"` /*  列表的页码，默认值为 1  */
	PageNo     int32    `json:"pageNo,omitempty"`     /*  列表的页码，默认值为 1, 推荐使用该字段, pageNumber 后续会废弃  */
	PageSize   int32    `json:"pageSize,omitempty"`   /*  每页数据量大小，取值 1-50  */
}

type CtelbQueryElbReatimeMetricResponse struct {
	StatusCode   int32                                          `json:"statusCode,omitempty"`   /*  返回状态码（800为成功，900为失败）  */
	Message      string                                         `json:"message,omitempty"`      /*  statusCode为900时的错误信息; statusCode为800时为success, 英文  */
	Description  string                                         `json:"description,omitempty"`  /*  statusCode为900时的错误信息; statusCode为800时为成功, 中文  */
	ErrorCode    string                                         `json:"errorCode,omitempty"`    /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
	ReturnObj    []*CtelbQueryElbReatimeMetricReturnObjResponse `json:"returnObj"`              /*  返回结果  */
	TotalCount   int32                                          `json:"totalCount,omitempty"`   /*  列表条目数  */
	CurrentCount int32                                          `json:"currentCount,omitempty"` /*  分页查询时每页的行数。  */
	TotalPage    int32                                          `json:"totalPage,omitempty"`    /*  总页数  */
	Error        string                                         `json:"error,omitempty"`        /*  statusCode为900时为业务细分错误码，三段式：product.module.code; statusCode为800时为SUCCESS  */
}

type CtelbQueryElbReatimeMetricReturnObjResponse struct {
	LastUpdated string                                                 `json:"lastUpdated,omitempty"` /*  最近更新时间  */
	RegionID    string                                                 `json:"regionID,omitempty"`    /*  资源池 ID  */
	DeviceID    string                                                 `json:"deviceID,omitempty"`    /*  弹性公网 IP  */
	ItemList    []*CtelbQueryElbReatimeMetricReturnObjItemListResponse `json:"itemList"`              /*  监控项值列表  */
}

type CtelbQueryElbReatimeMetricReturnObjItemListResponse struct {
	LbReqRate    string `json:"lbReqRate,omitempty"`    /*  请求频率  */
	LbLbin       string `json:"lbLbin,omitempty"`       /*  出吞吐量  */
	LbLbout      string `json:"lbLbout,omitempty"`      /*  入带宽峰值  */
	LbHrspOther  string `json:"lbHrspOther,omitempty"`  /*  HTTP 其它状态码统计数量  */
	LbHrsp2xx    string `json:"lbHrsp2xx,omitempty"`    /*  HTTP 2xx 状态码统计数量  */
	LbHrsp3xx    string `json:"lbHrsp3xx,omitempty"`    /*  HTTP 3xx 状态码统计数量  */
	LbHrsp4xx    string `json:"lbHrsp4xx,omitempty"`    /*  HTTP 4xx 状态码统计数量  */
	LbHrsp5xx    string `json:"lbHrsp5xx,omitempty"`    /*  HTTP 5xx 状态码统计数量  */
	LbNewcreate  string `json:"lbNewcreate,omitempty"`  /*  HTTP 新创建的链接数  */
	LbScur       string `json:"lbScur,omitempty"`       /*  HTTP  */
	LbInpkts     string `json:"lbInpkts,omitempty"`     /*  出带宽峰值  */
	LbOutpkts    string `json:"lbOutpkts,omitempty"`    /*  出带宽峰值  */
	LbActconn    string `json:"lbActconn,omitempty"`    /*  HTTP 活跃的链接数  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  采样时间  */
}
