package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbmMonitorHistoryInterfaceMetricApi
/* 查询物理机指定时间段内的网卡监控数据
 */type EbmMonitorHistoryInterfaceMetricApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbmMonitorHistoryInterfaceMetricApi(client *core.CtyunClient) *EbmMonitorHistoryInterfaceMetricApi {
	return &EbmMonitorHistoryInterfaceMetricApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebm/monitor_history_interface_metric",
			ContentType:  "application/json",
		},
	}
}

func (a *EbmMonitorHistoryInterfaceMetricApi) Do(ctx context.Context, credential core.Credential, req *EbmMonitorHistoryInterfaceMetricRequest) (*EbmMonitorHistoryInterfaceMetricResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("azName", req.AzName)
	ctReq.AddParam("instanceUUID", req.InstanceUUID)
	ctReq.AddParam("startTime", strconv.FormatInt(int64(req.StartTime), 10))
	ctReq.AddParam("endTime", strconv.FormatInt(int64(req.EndTime), 10))
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbmMonitorHistoryInterfaceMetricResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbmMonitorHistoryInterfaceMetricRequest struct {
	RegionID string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81">资源池列表查询</a>获取最新的天翼云资源池列表
	 */AzName string /*  可用区名称，您可以查看地域和可用区来了解可用区<br/>您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5855&data=87">资源池可用区查询</a><br/>注：查询结果中zoneList内返回存在可用区名称(即多可用区，本字段填写实际可用区名称)，若查询结果中zoneList为空（即为单可用区，本字段填写default）
	 */InstanceUUID string /*  实例UUID，您可以调用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=9715&data=87&isNormal=1">根据订单号查询uuid</a>获取实例UUID
	 */StartTime int32 /*  监控起始时间(unix秒时间戳)
	 */EndTime int32 /*  监控结束时间(unix秒时间戳)
	 */
}

type EbmMonitorHistoryInterfaceMetricResponse struct {
	StatusCode int32 `json:"statusCode"` /*  返回状态码(800为成功，900为失败)，默认值:800
	 */ErrorCode *string `json:"errorCode"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Error *string `json:"error"` /*  业务细分码，为product.module.code三段式码，详见错误码说明
	 */Message *string `json:"message"` /*  响应结果的描述，一般为英文描述
	 */Description *string `json:"description"` /*  响应结果的描述，一般为中文描述
	 */ReturnObj *EbmMonitorHistoryInterfaceMetricReturnObjResponse `json:"returnObj"` /*  返回参数，参考表returnObj
	 */
}

type EbmMonitorHistoryInterfaceMetricReturnObjResponse struct {
	Result []*EbmMonitorHistoryInterfaceMetricReturnObjResultResponse `json:"result"` /*  返回数据结果
	 */
}

type EbmMonitorHistoryInterfaceMetricReturnObjResultResponse struct {
	InstanceUUID *string `json:"instanceUUID"` /*  唯一键
	 */FuserLastUpdated *string `json:"fuserLastUpdated"` /*  用户最近更新时间
	 */ItemAggregateList *EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListResponse `json:"itemAggregateList"` /*  监控项值列表
	 */
}

type EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListResponse struct {
	NetOutBytesRate []*EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListNetOutBytesRateResponse `json:"netOutBytesRate"` /*  网卡输出流量速率监控数据列表
	 */NetInBytesRate []*EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListNetInBytesRateResponse `json:"netInBytesRate"` /*  网卡输入流量速率监控数据列表
	 */
}

type EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListNetOutBytesRateResponse struct {
	SamplingTime int32 `json:"samplingTime"` /*  监控获取时间
	 */Value *string `json:"value"` /*  网卡输出流量速率(kb/s)
	 */
}

type EbmMonitorHistoryInterfaceMetricReturnObjResultItemAggregateListNetInBytesRateResponse struct {
	SamplingTime int32 `json:"samplingTime"` /*  监控获取时间
	 */Value *string `json:"value"` /*  网卡输入流量速率(kb/s)
	 */
}
