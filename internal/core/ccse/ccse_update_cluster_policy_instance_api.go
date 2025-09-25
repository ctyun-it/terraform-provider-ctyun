package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateClusterPolicyInstanceApi
/* 调用该接口在指定集群中更新策略治理实例
 */type CcseUpdateClusterPolicyInstanceApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateClusterPolicyInstanceApi(client *core.CtyunClient) *CcseUpdateClusterPolicyInstanceApi {
	return &CcseUpdateClusterPolicyInstanceApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPut,
			UrlPath:      "/v2/cce/clusters/{clusterId}/policies/{policyName}",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateClusterPolicyInstanceApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateClusterPolicyInstanceRequest) (*CcseUpdateClusterPolicyInstanceResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("policyName", req.PolicyName)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseUpdateClusterPolicyInstanceResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateClusterPolicyInstanceRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	PolicyName string /*  策略治理实例名称，您可以通过查询策略治理实例接口来获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18035&data=128&isNormal=1&vid=121">查询策略治理实例</a>  */
	RegionId string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	InstanceId int64 `json:"instanceId,omitempty"` /*  实例Id，您可以通过查询策略治理实例接口来获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18035&data=128&isNormal=1&vid=121">查询策略治理实例</a>  */
	PolicyParameters *CcseUpdateClusterPolicyInstancePolicyParametersRequest `json:"policyParameters"`       /*  实例参数  */
	PolicyScope      string                                                  `json:"policyScope,omitempty"`  /*  所属命名空间  */
	PolicyAction     string                                                  `json:"policyAction,omitempty"` /*  规则治理动作，取值：<br />deny：拦截违规部署<br />warn：告警  */
}

type CcseUpdateClusterPolicyInstancePolicyParametersRequest struct{}

type CcseUpdateClusterPolicyInstanceResponse struct {
	StatusCode int32                                             `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                            `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseUpdateClusterPolicyInstanceReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                            `json:"error,omitempty"`      /*  错误码  */
}

type CcseUpdateClusterPolicyInstanceReturnObjResponse struct {
	Instances []string `json:"instances"` /*  部署的实例名称列表  */
}
