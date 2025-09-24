package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpdateNodePoolKubeletConfigApi
/* 调用该接口配置节点池Kubelet。
 */type CcseUpdateNodePoolKubeletConfigApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpdateNodePoolKubeletConfigApi(client *core.CtyunClient) *CcseUpdateNodePoolKubeletConfigApi {
	return &CcseUpdateNodePoolKubeletConfigApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/kubeletconfig/update",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpdateNodePoolKubeletConfigApi) Do(ctx context.Context, credential core.Credential, req *CcseUpdateNodePoolKubeletConfigRequest) (*CcseUpdateNodePoolKubeletConfigResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder = builder.ReplaceUrl("nodePoolId", req.NodePoolId)
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
	var resp CcseUpdateNodePoolKubeletConfigResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpdateNodePoolKubeletConfigRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>  */
	KubeletArgs *CcseUpdateNodePoolKubeletConfigKubeletArgsRequest `json:"kubeletArgs"` /*  其他可支持的Kubelet参数:registryPullQPS，registryBurst， podPidsLimit，eventRecordQPS，eventBurst    , kubeAPIQPS    , kubeAPIBurst, cpuManagerPolicy, topologyManagerScope, cpuCFSQuota, maxPods,  */
}

type CcseUpdateNodePoolKubeletConfigKubeletArgsRequest struct {
	MaxPods         int32 `json:"maxPods,omitempty"`         /*  最大pod数  */
	RegistryPullQPS int32 `json:"registryPullQPS,omitempty"` /*  注册拉取qps  */
	RegistryBurst   int32 `json:"registryBurst,omitempty"`   /*  注册表突发  */
}

type CcseUpdateNodePoolKubeletConfigResponse struct {
	RequestId string                                            `json:"requestId,omitempty"` /*  请求ID  */
	Message   string                                            `json:"message,omitempty"`   /*  响应信息  */
	ReturnObj *CcseUpdateNodePoolKubeletConfigReturnObjResponse `json:"returnObj"`           /*  响应对象  */
}

type CcseUpdateNodePoolKubeletConfigReturnObjResponse struct {
	TaskId string `json:"taskId,omitempty"` /*  任务ID  */
}
