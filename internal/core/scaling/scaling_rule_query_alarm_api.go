package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingRuleQueryAlarmApi
/* 查询弹性伸缩的告警策略<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingRuleQueryAlarmApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingRuleQueryAlarmApi(client *core.CtyunClient) *ScalingRuleQueryAlarmApi {
	return &ScalingRuleQueryAlarmApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/rule/query-alarm",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingRuleQueryAlarmApi) Do(ctx context.Context, credential core.Credential, req *ScalingRuleQueryAlarmRequest) (*ScalingRuleQueryAlarmResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingRuleQueryAlarmRequest
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
	var resp ScalingRuleQueryAlarmResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingRuleQueryAlarmRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID  int32  `json:"groupID,omitempty"`  /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	RuleID   int32  `json:"ruleID,omitempty"`   /*  伸缩策略ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4990&data=93">查询弹性伸缩组内的策略列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5069&data=93">创建一条伸缩策略  */
}

type ScalingRuleQueryAlarmResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                  `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                  `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                  `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingRuleQueryAlarmReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                  `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingRuleQueryAlarmReturnObjResponse struct {
	TriggerID          string `json:"triggerID"`          /*  告警规则ID  */
	Name               string `json:"name"`               /*  告警规则名称  */
	MetricName         string `json:"metricName"`         /*  监控指标名称  */
	Statistics         string `json:"statistics"`         /*  聚合方法。取值范围：<br/>avg：平均值。<br/>max：最大值。<br/>min：最小值。  */
	ComparisonOperator string `json:"comparisonOperator"` /*  比较符。取值范围：<br/>ge：大于等于。<br/>le：小于等于。<br/>gt：大于。<br/>lt：小于。  */
	Threshold          int32  `json:"threshold"`          /*  阈值  */
	Period             string `json:"period"`             /*  监控周期，单位：分钟  */
	EvaluationCount    int32  `json:"evaluationCount"`    /*  连续出现次数  */
	Cooldown           int32  `json:"cooldown"`           /*  冷却时间，单位：秒  */
	Status             int32  `json:"status"`             /*  启用状态。<br>取值范围：<br>1: 启用；<br>2: 停用。  */
}
