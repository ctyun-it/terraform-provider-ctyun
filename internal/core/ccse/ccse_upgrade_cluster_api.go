package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseUpgradeClusterApi
/* 调用该接口升级集群。
 */type CcseUpgradeClusterApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseUpgradeClusterApi(client *core.CtyunClient) *CcseUpgradeClusterApi {
	return &CcseUpgradeClusterApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v2/cce/clusters/{clusterId}/upgrade",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseUpgradeClusterApi) Do(ctx context.Context, credential core.Credential, req *CcseUpgradeClusterRequest) (*CcseUpgradeClusterResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
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
	var resp CcseUpgradeClusterResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseUpgradeClusterRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	NextVersion string `json:"nextVersion,omitempty"` /*  集群可升级版本  */
	Version     string `json:"version,omitempty"`     /*  当前版本  */
	Concurrency int32  `json:"concurrency,omitempty"` /*  worker升级并发数量  */
}

type CcseUpgradeClusterResponse struct {
	StatusCode int32                                `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                               `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  *CcseUpgradeClusterReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                               `json:"error,omitempty"`      /*  错误码  */
}

type CcseUpgradeClusterReturnObjResponse struct {
	TaskId         int64  `json:"taskId,omitempty"`         /*  任务id  */
	ClusterId      string `json:"clusterId,omitempty"`      /*  集群ID  */
	ParentTaskId   int64  `json:"parentTaskId,omitempty"`   /*  parentTaskId  */
	NodeName       string `json:"nodeName,omitempty"`       /*  节点名称  */
	NodeType       string `json:"nodeType,omitempty"`       /*  节点类型，包括master0, master, node  */
	MasterPlanHash string `json:"masterPlanHash,omitempty"` /*  masterPlanHash  */
	WorkerPlanHash string `json:"workerPlanHash,omitempty"` /*  workerPlanHash  */
	TaskStatus     string `json:"taskStatus,omitempty"`     /*  任务状态，start：开始升级，end：完成升级，pause：暂停升级  */
	Version        string `json:"version,omitempty"`        /*  版本  */
	Status         int32  `json:"status,omitempty"`         /*  状态码，0：无效 1：有效 2：删除  */
	CreatedBy      int64  `json:"createdBy,omitempty"`      /*  创建者  */
	CreatedTime    string `json:"createdTime,omitempty"`    /*  创建时间  */
	UpdatedTime    string `json:"updatedTime,omitempty"`    /*  更新时间  */
}
