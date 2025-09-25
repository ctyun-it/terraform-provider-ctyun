package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsQueryVmNetworkHistoryMetricDataV41Api
/* 该接口提供用户查询指定时间段内的网卡监控数据功能<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br /><b>注意事项：</b><br />&emsp;&emsp;当前查询结果以分页形式进行展示，单次查询最多显示50条数据<br />&emsp;&emsp;调用接口时，如果监控项返回的值为"[]"则说明未获取到监控项
 */type CtecsQueryVmNetworkHistoryMetricDataV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsQueryVmNetworkHistoryMetricDataV41Api(client *core.CtyunClient) *CtecsQueryVmNetworkHistoryMetricDataV41Api {
	return &CtecsQueryVmNetworkHistoryMetricDataV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/vm-network-history-metric-data",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsQueryVmNetworkHistoryMetricDataV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsQueryVmNetworkHistoryMetricDataV41Request) (*CtecsQueryVmNetworkHistoryMetricDataV41Response, error) {
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
	var resp CtecsQueryVmNetworkHistoryMetricDataV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsQueryVmNetworkHistoryMetricDataV41Request struct {
	RegionID     string   `json:"regionID,omitempty"`  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DeviceIDList []string `json:"deviceIDList"`        /*  云主机ID列表，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	Period       int32    `json:"period,omitempty"`    /*  聚合周期，单位秒，注：默认值为300  */
	StartTime    string   `json:"startTime,omitempty"` /*  必传参数，查询起始时间戳  */
	EndTime      string   `json:"endTime,omitempty"`   /*  必传参数，查询终止时间戳  */
	PageNo       int32    `json:"pageNo,omitempty"`    /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	Page         int32    `json:"page,omitempty"`      /*  页码，取值范围：正整数（≥1），注：默认值为1，后续该字段可能废弃  */
	PageSize     int32    `json:"pageSize,omitempty"`  /*  每页记录数目，取值范围：[1, 50]，注：默认值为20  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41Response struct {
	StatusCode  int32                                                     `json:"statusCode,omitempty"`  /*  返回码状态(800为成功，900为失败)，默认值：800  */
	ErrorCode   string                                                    `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                                    `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                                    `json:"message,omitempty"`     /*  英文描述信息  */
	MsgDesc     string                                                    `json:"msgDesc,omitempty"`     /*  中文描述信息  */
	Description string                                                    `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResponse `json:"returnObj"`             /*  返回参数，参考表returnObj  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResponse struct {
	Result       []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultResponse `json:"result"`                 /*  result对象  */
	CurrentCount int32                                                             `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                             `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                             `json:"totalPage,omitempty"`    /*  总页数  */
	PageSize     int32                                                             `json:"pageSize,omitempty"`     /*  每页记录数目  */
	Page         int32                                                             `json:"page,omitempty"`         /*  页码  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultResponse struct {
	FUID              string                                                                           `json:"fUID,omitempty"`             /*  唯一键  */
	FuserLastUpdated  string                                                                           `json:"fuserLastUpdated,omitempty"` /*  用户最近更新时间  */
	RegionID          string                                                                           `json:"regionID,omitempty"`         /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DeviceUUID        string                                                                           `json:"deviceUUID,omitempty"`       /*  云主机ID  */
	ItemAggregateList *CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListResponse `json:"itemAggregateList"`          /*  监控信息  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListResponse struct {
	Net_in_bytes_rate                    []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNet_in_bytes_rateResponse                    `json:"net_in_bytes_rate"`                    /*  网络流入速率，下级对象中value的单位为（bit/s）  */
	Net_out_bytes_rate                   []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNet_out_bytes_rateResponse                   `json:"net_out_bytes_rate"`                   /*  网络流出速率，下级对象中value的单位为（bit/s）  */
	Network_incoming_packets_rate_inband []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_packets_rate_inbandResponse `json:"network_incoming_packets_rate_inband"` /*  网卡包接收的正确率，下级对象中value的单位为（%）  */
	Network_outing_packets_rate_inband   []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_packets_rate_inbandResponse   `json:"network_outing_packets_rate_inband"`   /*  网卡包发送的正确率，下级对象中value的单位为（%）  */
	Network_incoming_errs_rate_inband    []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_errs_rate_inbandResponse    `json:"network_incoming_errs_rate_inband"`    /*  接收误包率，下级对象中value的单位为（%）  */
	Network_outing_drop_rate_inband      []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_drop_rate_inbandResponse      `json:"network_outing_drop_rate_inband"`      /*  发送丢包率，下级对象中value的单位为（%）  */
	Network_outing_errs_rate_inband      []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_errs_rate_inbandResponse      `json:"network_outing_errs_rate_inband"`      /*  发送误包率，下级对象中value的单位为（%）  */
	Network_incoming_drop_rate_inband    []*CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_drop_rate_inbandResponse    `json:"network_incoming_drop_rate_inband"`    /*  接收丢包率，下级对象中value的单位为（%）  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNet_in_bytes_rateResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNet_out_bytes_rateResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_packets_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_packets_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_errs_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_drop_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_outing_errs_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}

type CtecsQueryVmNetworkHistoryMetricDataV41ReturnObjResultItemAggregateListNetwork_incoming_drop_rate_inbandResponse struct {
	Value        string `json:"value,omitempty"`        /*  监控项值  */
	SamplingTime int32  `json:"samplingTime,omitempty"` /*  监控获取时间  */
}
