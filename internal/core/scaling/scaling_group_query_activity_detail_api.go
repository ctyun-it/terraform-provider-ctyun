package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ScalingGroupQueryActivityDetailApi
/* 根据一个伸缩活动的ID查询一个伸缩活动的详细信息<br /><b>准备工作：</b><br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u6784%u9020%u8BF7%u6C42&data=93">构造请求</a><br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=%u8BA4%u8BC1%u9274%u6743&data=93">认证鉴权</a><br />
 */type ScalingGroupQueryActivityDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewScalingGroupQueryActivityDetailApi(client *core.CtyunClient) *ScalingGroupQueryActivityDetailApi {
	return &ScalingGroupQueryActivityDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/scaling/group/query-activity-detail",
			ContentType:  "application/json",
		},
	}
}

func (a *ScalingGroupQueryActivityDetailApi) Do(ctx context.Context, credential core.Credential, req *ScalingGroupQueryActivityDetailRequest) (*ScalingGroupQueryActivityDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*ScalingGroupQueryActivityDetailRequest
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
	var resp ScalingGroupQueryActivityDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ScalingGroupQueryActivityDetailRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	GroupID  int32  `json:"groupID,omitempty"`  /*  伸缩组ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=4996&data=93">查询伸缩组列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5081&data=93">创建一个伸缩组</a>  */
	ActiveID int32  `json:"activeID,omitempty"` /*  伸缩活动ID <br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=19&api=5005&data=93">查询伸缩活动ID</a>  */
}

type ScalingGroupQueryActivityDetailResponse struct {
	StatusCode  int32                                             `json:"statusCode"`  /*  返回码：800表示成功，900表示失败  */
	ErrorCode   string                                            `json:"errorCode"`   /*  业务细分码，为product.module.code三段式码  */
	Message     string                                            `json:"message"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                            `json:"description"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *ScalingGroupQueryActivityDetailReturnObjResponse `json:"returnObj"`   /*  成功时返回的数据，参见表returnObj  */
	Error       string                                            `json:"error"`       /*  业务细分码，为product.module.code三段式码  */
}

type ScalingGroupQueryActivityDetailReturnObjResponse struct {
	RuleFailReason      string                                                          `json:"ruleFailReason"`      /*  规则失败原因  */
	AfterCount          int32                                                           `json:"afterCount"`          /*  活动后计数  */
	EndTime             string                                                          `json:"endTime"`             /*  结束时间  */
	BeforeCount         int32                                                           `json:"beforeCount"`         /*  活动前计数  */
	RuleID              string                                                          `json:"ruleID"`              /*  伸缩策略ID  */
	StartTime           string                                                          `json:"startTime"`           /*  开始时间  */
	FailReason          string                                                          `json:"failReason"`          /*  失败原因  */
	InstanceList        []*ScalingGroupQueryActivityDetailReturnObjInstanceListResponse `json:"instanceList"`        /*  虚机列表  */
	ExecutionMode       int32                                                           `json:"executionMode"`       /*  执行方式。<br>取值范围：<br>1：自动执行策略。<br>2：手动执行策略。<br>3：手动移入实例。<br>4：手动移出实例。<br>5：新建伸缩组满足最小数。<br>6：修改伸缩组满足最大最小限制。<br>7：健康检查移入。<br>8：健康检查移出。  */
	GroupID             int32                                                           `json:"groupID"`             /*  伸缩组ID  */
	RuleExpectDelta     int32                                                           `json:"ruleExpectDelta"`     /*  策略预期可变化数量  */
	ExecutionResult     int32                                                           `json:"executionResult"`     /*  执行结果。<br>取值范围：<br>0：执行中。<br>1：成功。<br>2：失败。  */
	ExecutionDate       string                                                          `json:"executionDate"`       /*  执行时间  */
	RuleExecutionResult int32                                                           `json:"ruleExecutionResult"` /*  规则执行结果。<br/>取值范围：<br/>0：执行中。<br/>1: 成功。<br/>2：失败。  */
	ActiveID            int32                                                           `json:"activeID"`            /*  伸缩活动ID  */
	RuleRealDelta       int32                                                           `json:"ruleRealDelta"`       /*  策略实际可变化数量  */
	Description         string                                                          `json:"description"`         /*  说明  */
}

type ScalingGroupQueryActivityDetailReturnObjInstanceListResponse struct {
	InstanceID   string `json:"instanceID"`   /*  实例ID  */
	InstanceName string `json:"instanceName"` /*  实例名称  */
}
