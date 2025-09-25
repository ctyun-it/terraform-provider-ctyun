package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetPolicyGovernanceApi
/* 调用该接口查看集群策略治理详情
 */type CcseGetPolicyGovernanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetPolicyGovernanceApi(client *core.CtyunClient) *CcseGetPolicyGovernanceApi {
	return &CcseGetPolicyGovernanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/policies/governance",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetPolicyGovernanceApi) Do(ctx context.Context, credential core.Credential, req *CcseGetPolicyGovernanceRequest) (*CcseGetPolicyGovernanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetPolicyGovernanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetPolicyGovernanceRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
}

type CcseGetPolicyGovernanceResponse struct {
	StatusCode int32                                     `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                    `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetPolicyGovernanceReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                    `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetPolicyGovernanceReturnObjResponse struct {
	OnState    []*CcseGetPolicyGovernanceReturnObjOnStateResponse  `json:"onState"`    /*  当前集群中开启的不同等级策略计数统计  */
	Violations *CcseGetPolicyGovernanceReturnObjViolationsResponse `json:"violations"` /*  集群中针对不同策略类型的拦截和告警的审计计数统计列表  */
	AdmitLog   *CcseGetPolicyGovernanceReturnObjAdmitLogResponse   `json:"admitLog"`   /*  集群当前策略治理审计日志  */
}

type CcseGetPolicyGovernanceReturnObjOnStateResponse struct {
	Severity     string `json:"severity,omitempty"`     /*  策略治理等级  */
	Total        int64  `json:"total,omitempty"`        /*  该等级下策略种类总数  */
	EnabledCount int64  `json:"enabledCount,omitempty"` /*  当前开启的策略种类计数  */
}

type CcseGetPolicyGovernanceReturnObjViolationsResponse struct {
	Deny      []*CcseGetPolicyGovernanceReturnObjViolationsDenyResponse    `json:"deny"`      /*  拦截模式下不同治理等级的违规计数统计  */
	Warn      []*CcseGetPolicyGovernanceReturnObjViolationsWarnResponse    `json:"warn"`      /*  告警模式下不同治理等级的违规计数统计  */
	DenyCount *CcseGetPolicyGovernanceReturnObjViolationsDenyCountResponse `json:"denyCount"` /*  告警模式下High治理等级的违规计数统计  */
	WarnCount *CcseGetPolicyGovernanceReturnObjViolationsWarnCountResponse `json:"warnCount"` /*  告警模式下Medium治理等级的违规计数统计  */
}

type CcseGetPolicyGovernanceReturnObjAdmitLogResponse struct {
	Progress string                                                 `json:"progress,omitempty"` /*  查询结果的状态:Complete, Incomplete  */
	Count    int64                                                  `json:"count,omitempty"`    /*  当前查询到的日志总数  */
	Log      []*CcseGetPolicyGovernanceReturnObjAdmitLogLogResponse `json:"log"`                /*  策略治理审计日志内容  */
}

type CcseGetPolicyGovernanceReturnObjViolationsDenyResponse struct {
	PolicyName string `json:"policyName,omitempty"` /*  策略名称  */
	PolicyDesc string `json:"policyDesc,omitempty"` /*  策略描述  */
	Violations int32  `json:"violations,omitempty"` /*  集群中对应规则类型下被拦截的违规计数统计  */
	Severity   string `json:"severity,omitempty"`   /*  策略治理等级  */
}

type CcseGetPolicyGovernanceReturnObjViolationsWarnResponse struct {
	PolicyName string `json:"policyName,omitempty"` /*  策略名称  */
	PolicyDesc string `json:"policyDesc,omitempty"` /*  策略描述  */
	Violations int32  `json:"violations,omitempty"` /*  集群中对应规则类型下被拦截的违规计数统计  */
	Severity   string `json:"severity,omitempty"`   /*  策略治理等级  */
}

type CcseGetPolicyGovernanceReturnObjViolationsDenyCountResponse struct {
	High   int32 `json:"high,omitempty"`   /*  高级别违规计数统计  */
	Medium int32 `json:"medium,omitempty"` /*  中级别违规计数统计  */
	Low    int32 `json:"low,omitempty"`    /*  低级别违规计数统计  */
}

type CcseGetPolicyGovernanceReturnObjViolationsWarnCountResponse struct {
	High   int32 `json:"high,omitempty"`   /*  高级别违规计数统计  */
	Medium int32 `json:"medium,omitempty"` /*  中级别违规计数统计  */
	Low    int32 `json:"low,omitempty"`    /*  低级别违规计数统计  */
}

type CcseGetPolicyGovernanceReturnObjAdmitLogLogResponse struct {
	ClusterId         string `json:"clusterId,omitempty"`         /*  目标集群ID  */
	ConstraintKind    int64  `json:"constraintKind,omitempty"`    /*  策略类型名称  */
	ResourceName      string `json:"resourceName,omitempty"`      /*  目标资源名称  */
	ResourceKind      string `json:"resourceKind,omitempty"`      /*  目标资源类型  */
	ResourceNamespace string `json:"resourceNamespace,omitempty"` /*  目标资源命名空间  */
	Msg               string `json:"msg,omitempty"`               /*  策略治理审计日志信息  */
}
