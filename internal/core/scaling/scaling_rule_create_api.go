package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingRuleCreateApi
/* 创建一个弹性伸缩策略<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingRuleCreateApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingRuleCreateApi(client *core.CtyunClient) *ScalingRuleCreateApi {
	return &ScalingRuleCreateApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/rule/create",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingRuleCreateApi) Do(ctx context.Context, credential core.Credential, req *ScalingRuleCreateRequest) (*ScalingRuleCreateResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingRuleCreateRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ScalingRuleCreateResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingRuleCreateRequest struct {
	RegionID      string                              `json:"regionID,omitempty"`      /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID       int64                               `json:"groupID,omitempty"`       /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	Name          string                              `json:"name,omitempty"`          /*  伸缩策略名称  */
	RawType       int32                               `json:"type,omitempty"`          /*  策略类型。<br>取值范围：<br>1：告警策略<br>2：定时策略<br>3：周期策略<br>4：目标追踪策略  */
	OperateUnit   int32                               `json:"operateUnit,omitempty"`   /*  操作单位，type=1、type=2、type=3时必填。<br> 取值范围：<br>1：个数。<br>2：百分比。  */
	OperateCount  int32                               `json:"operateCount,omitempty"`  /*  调整值，type=1、type=2、type=3时必填。  */
	Action        int32                               `json:"action,omitempty"`        /*  执行动作，type=1、type=2、type=3时必填。<br>取值范围：<br>1：增加<br>2：减少<br>3：设置为  */
	Cycle         int32                               `json:"cycle,omitempty"`         /*  循环方式，type=3时必填。取值范围：<br>1：按月循环。<br>2：按周循环。<br>3：按天循环。  */
	Day           []int32                             `json:"day"`                     /*  执行日期，当cycle为1时必填，且Array元素限制为[1-31]中的整数且不重复；<br/>当cycle为2时必填，且Array元素限制为[1-7]中的整数且不重复；  */
	EffectiveFrom string                              `json:"effectiveFrom,omitempty"` /*  周期策略生效开始时间，type=3时必填，格式为：%Y-%m-%d %H:%M:%S  */
	EffectiveTill string                              `json:"effectiveTill,omitempty"` /*  周期策略生效截止时间，type=3时必填，格式为：%Y-%m-%d %H:%M:%S  */
	ExecutionTime string                              `json:"executionTime,omitempty"` /*  执行时间，type=2、type=3时必填，格式为：%Y-%m-%d %H:%M:%S  */
	Cooldown      int32                               `json:"cooldown,omitempty"`      /*  冷却时间（type=1）或预热时间（type=4），type=1、type=4时必填，单位：秒  */
	TriggerObj    *ScalingRuleCreateTriggerObjRequest `json:"triggerObj"`              /*  告警策略参数表，type=1时必填  */
	TargetObj     *ScalingRuleCreateTargetObjRequest  `json:"targetObj"`               /*  目标追踪策略参数表，type=4时必填  */
}

type ScalingRuleCreateTriggerObjRequest struct {
	Name               string `json:"name,omitempty"`               /*  告警规则名称  */
	MetricName         string `json:"metricName,omitempty"`         /*  监控指标名称，支持：</br>cpu_util：CPU使用率，单位“%”</br>mem_util：内存使用率，单位“%”</br>network_incoming_bytes_rate_inband：网络流入速率，单位“Kbps”</br>network_outing_bytes_rate_inband：网络流出速率，单位“Kbps”</br>disk_read_bytes_rate：磁盘读速率，单位“KBps”</br>disk_write_bytes_rate：磁盘写速率，单位“KBps”</br>disk_read_requests_rate：磁盘读请求速率，单位“IOPS”</br>disk_write_requests_rate：磁盘写请求速率，单位“IOPS”  */
	Statistics         string `json:"statistics,omitempty"`         /*  聚合方法。取值范围：<br/>avg：平均值。<br/>max：最大值。<br/>min：最小值。  */
	ComparisonOperator string `json:"comparisonOperator,omitempty"` /*  比较符。取值范围：<br/>ge：大于等于。<br/>le：小于等于。<br/>gt：大于。<br/>lt：小于。  */
	Threshold          int32  `json:"threshold,omitempty"`          /*  阈值  */
	Period             string `json:"period,omitempty"`             /*  监控周期，单位：分钟  */
	EvaluationCount    int32  `json:"evaluationCount,omitempty"`    /*  连续出现次数  */
}

type ScalingRuleCreateTargetObjRequest struct {
	MetricName              string `json:"metricName,omitempty"`              /*  监控指标名称，仅支持：</br>cpu_util：CPU使用率，单位“%”</br>network_incoming_bytes_rate_inband：网络流入速率，单位“Kbps”</br>network_outing_bytes_rate_inband：网络流出速率，单位“Kbps”  */
	TargetValue             int32  `json:"targetValue,omitempty"`             /*  追踪目标值。伸缩组监控指标维持的目标值，将通过添加或删除实例来将指标维持在目标值附近。取值范围：</br>cpu_util：[1,100]</br>network_incoming_bytes_rate_inband：[1,99999999]</br>network_outcoming_bytes_rate_inband：[1,99999999]  */
	ScaleOutEvaluationCount int32  `json:"scaleOutEvaluationCount,omitempty"` /*  扩容连续告警次数。创建目标追踪策略后，会自动创建告警规则。本参数用于指定目标追踪策略触发扩容告警时，所需连续满足告警条件的次数。取值范围：[1,100]  */
	ScaleInEvaluationCount  int32  `json:"scaleInEvaluationCount,omitempty"`  /*  缩容连续告警次数。创建目标追踪策略后，会自动创建告警规则。本参数用于指定目标追踪策略触发缩容告警时，所需连续满足告警条件的次数。取值范围：[1,100]  */
	OperateRange            int32  `json:"operateRange,omitempty"`            /*  缩容波动范围。目标追踪策略触发缩容活动时的目标时，当伸缩组监控指标<目标值*（1-波动范围）时，触发缩容活动。取值范围：[10, 20]  */
	DisableScaleIn          bool   `json:"disableScaleIn"`                    /*  是否禁用缩容，默认false。取值范围：</br>true：开启禁用缩容，目标追踪策略不会触发缩容活动</br>false：关闭禁用缩容，目标追踪策略允许触发伸缩活动  */
}

type ScalingRuleCreateResponse struct {
	StatusCode  int32                               `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                              `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                              `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                              `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingRuleCreateReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                              `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingRuleCreateReturnObjResponse struct {
	RuleID int64 `json:"ruleID"` /*  伸缩策略ID  */
}
