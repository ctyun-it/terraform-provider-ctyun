package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseListClusterPolicyInstancesApi
/* 调用该接口查看集群中当前部署的策略实例
 */type CcseListClusterPolicyInstancesApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseListClusterPolicyInstancesApi(client *core.CtyunClient) *CcseListClusterPolicyInstancesApi {
	return &CcseListClusterPolicyInstancesApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/policies/instances",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseListClusterPolicyInstancesApi) Do(ctx context.Context, credential core.Credential, req *CcseListClusterPolicyInstancesRequest) (*CcseListClusterPolicyInstancesResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PolicyName != "" {
		ctReq.AddParam("policyName", req.PolicyName)
	}
	if req.InstanceName != "" {
		ctReq.AddParam("instanceName", req.InstanceName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseListClusterPolicyInstancesResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseListClusterPolicyInstancesRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	PolicyName string /*  策略治理规则名称，您可以通过查询策略治理规则列表接口来获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18948&data=128&isNormal=1&vid=121">查询策略治理规则列表</a>  */
	InstanceName string /*  策略规则实例名称  */
}

type CcseListClusterPolicyInstancesResponse struct {
	StatusCode int32                                              `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                             `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  []*CcseListClusterPolicyInstancesReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                             `json:"error,omitempty"`      /*  错误码  */
}

type CcseListClusterPolicyInstancesReturnObjResponse struct {
	InstanceId       int64  `json:"instanceId,omitempty"`       /*  规则ID  */
	InstanceName     string `json:"instanceName,omitempty"`     /*  实例名称  */
	ClusterId        string `json:"clusterId,omitempty"`        /*  集群标识ID  */
	PolicyName       string `json:"policyName,omitempty"`       /*  策略规则名称  */
	PolicyCategory   string `json:"policyCategory,omitempty"`   /*  策略模板类型  */
	PolicyDesc       string `json:"policyDesc,omitempty"`       /*  策略模板描述  */
	PolicyParameters string `json:"policyParameters,omitempty"` /*  当前规则实例的配置参数  */
	PolicySeverity   string `json:"policySeverity,omitempty"`   /*  规则治理等级，取值：<br /> high<br /> medium<br /> low  */
	PolicyScope      string `json:"policyScope,omitempty"`      /*  策略实例实施范围，默认*代表全部命名空间，否则返回命名空间名称，多个命名空间用逗号分隔  */
	PolicyAction     string `json:"policyAction,omitempty"`     /*  规则治理动作，取值：<br /> deny：拦截违规部署<br /> warn：告警  */
	Status           int32  `json:"status,omitempty"`           /*  状态，取值：<br /> 1：有效<br /> 2：删除  */
	CreatedTime      string `json:"createdTime,omitempty"`      /*  创建时间  */
	CreatedBy        int64  `json:"createdBy,omitempty"`        /*  创建人  */
	ModifiedTime     string `json:"modifiedTime,omitempty"`     /*  修改时间  */
	ModifiedBy       int64  `json:"modifiedBy,omitempty"`       /*  修改人  */
}
