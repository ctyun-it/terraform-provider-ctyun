package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseGetClusterPolicyStatusApi
/* 调用该接口查看集群中策略治理规则下的实例详情
 */type CcseGetClusterPolicyStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetClusterPolicyStatusApi(client *core.CtyunClient) *CcseGetClusterPolicyStatusApi {
	return &CcseGetClusterPolicyStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/policies/status",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetClusterPolicyStatusApi) Do(ctx context.Context, credential core.Credential, req *CcseGetClusterPolicyStatusRequest) (*CcseGetClusterPolicyStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.PolicyName != "" {
		ctReq.AddParam("policyName", req.PolicyName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetClusterPolicyStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetClusterPolicyStatusRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	PolicyName string /*  策略名称，您可以通过查询策略治理规则列表接口来获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18948&data=128&isNormal=1&vid=121">查询策略治理规则列表</a>  */
}

type CcseGetClusterPolicyStatusResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                       `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseGetClusterPolicyStatusReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetClusterPolicyStatusReturnObjResponse struct {
	PolicyInstances        []*CcseGetClusterPolicyStatusReturnObjPolicyInstancesResponse      `json:"policyInstances"`        /*  不同策略类型下的策略实例计数列表  */
	InstancesSeverityCount *CcseGetClusterPolicyStatusReturnObjInstancesSeverityCountResponse `json:"instancesSeverityCount"` /*  Map<String,Object>类型，集群中当前部署的不同治理等级的策略实例计数  */
}

type CcseGetClusterPolicyStatusReturnObjPolicyInstancesResponse struct {
	PolicyCategory       string `json:"policyCategory,omitempty"`       /*  策略类型  */
	PolicyName           string `json:"policyName,omitempty"`           /*  策略名称  */
	PolicyDescription    string `json:"policyDescription,omitempty"`    /*  策略描述  */
	PolicySeverity       string `json:"policySeverity,omitempty"`       /*  策略治理等级  */
	PolicyInstancesCount int32  `json:"policyInstancesCount,omitempty"` /*  已部署的策略实例计数  */
}

type CcseGetClusterPolicyStatusReturnObjInstancesSeverityCountResponse struct {
	High   int32 `json:"high,omitempty"`   /*  高级别实例数量  */
	Medium int32 `json:"medium,omitempty"` /*  中级别实例数量  */
	Low    int32 `json:"low,omitempty"`    /*  低级别实例数量  */
}
