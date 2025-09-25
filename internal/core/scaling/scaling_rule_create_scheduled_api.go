package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingRuleCreateScheduledApi
/* 在伸缩组中创建一个定时策略<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingRuleCreateScheduledApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingRuleCreateScheduledApi(client *core.CtyunClient) *ScalingRuleCreateScheduledApi {
	return &ScalingRuleCreateScheduledApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/rule/create-scheduled",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingRuleCreateScheduledApi) Do(ctx context.Context, credential core.Credential, req *ScalingRuleCreateScheduledRequest) (*ScalingRuleCreateScheduledResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingRuleCreateScheduledRequest
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
	var resp ScalingRuleCreateScheduledResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingRuleCreateScheduledRequest struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID       int32  `json:"groupID,omitempty"`       /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	Name          string `json:"name,omitempty"`          /*  伸缩策略名称<br>请注意不能与当前伸缩组内的其他策略冲突  */
	ExecutionTime string `json:"executionTime,omitempty"` /*  定时策略执行时间，格式为：%Y-%m-%d %H:%M:%S  */
	Action        int32  `json:"action,omitempty"`        /*  执行动作。<br>取值范围：<br>1：增加<br>2：减少<br>3：设置为  */
	OperateUnit   int32  `json:"operateUnit,omitempty"`   /*  操作单位。<br> 取值范围：<br>1：个数。<br>2：百分比。  */
	OperateCount  int32  `json:"operateCount,omitempty"`  /*  调整值  */
}

type ScalingRuleCreateScheduledResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                       `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                       `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                       `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingRuleCreateScheduledReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                       `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingRuleCreateScheduledReturnObjResponse struct {
	RuleID int32 `json:"ruleID"` /*  伸缩策略ID  */
}
