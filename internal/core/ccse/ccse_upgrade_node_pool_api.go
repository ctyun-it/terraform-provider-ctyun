package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpgradeNodePoolApi
/* 调用该接口创建节点池升级运维任务。
 */type CcseUpgradeNodePoolApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpgradeNodePoolApi(client *core.CtyunClient) *CcseUpgradeNodePoolApi {
	return &CcseUpgradeNodePoolApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/nodepool/{nodePoolId}/upgrade",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpgradeNodePoolApi) Do(ctx context.Context, credential core.Credential, req *CcseUpgradeNodePoolRequest) (*CcseUpgradeNodePoolResponse, error) {
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
	var resp CcseUpgradeNodePoolResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpgradeNodePoolRequest struct {
	ClusterId  string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	NodePoolId string /*  节点池ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105" target="_blank">如何获取接口URI中参数</a>。  */
	RegionId   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&amp;api=5851&amp;data=87">资源池列表查询</a>  */
	OsUpgrade             bool   `json:"osUpgrade"`                       /*  是否升级操作系统  */
	UpgradeImageName      string `json:"upgradeImageName,omitempty"`      /*  目标操作系统名称（最新操作系统名称可参考产品文档”升级节点池“章节）  */
	UpgradeImageID        string `json:"upgradeImageID,omitempty"`        /*  目标操作系统镜像ID，请联系相关技术人员获取对应操作系统名称的ID  */
	KubeletUpgrade        bool   `json:"kubeletUpgrade"`                  /*  是否升级kubelet  */
	UpgradeKubeletVersion string `json:"upgradeKubeletVersion,omitempty"` /*  目标kubelet版本名称（最新kubelet版本可参考产品文档”升级节点池“章节）  */
	RuntimeUpgrade        bool   `json:"runtimeUpgrade"`                  /*  否升级容器运行时  */
	RuntimeType           string `json:"runtimeType,omitempty"`           /*  节点池容器运行时类型，目前支持containerd运行时升级  */
	UpgradeRuntimeVersion string `json:"upgradeRuntimeVersion,omitempty"` /*  目标运行时版本名称（最新运行时版本可参考产品文档”升级节点池“章节）  */
}

type CcseUpgradeNodePoolResponse struct {
	ReturnObj *CcseUpgradeNodePoolReturnObjResponse `json:"returnObj"` /*  响应对象  */
}

type CcseUpgradeNodePoolReturnObjResponse struct {
	TaskId string `json:"taskId,omitempty"` /*  任务ID  */
}
