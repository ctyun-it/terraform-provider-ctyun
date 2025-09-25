package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CcseGetUpgradeClusterStatusApi
/* 调用该接口查看集群升级状态。
 */type CcseGetUpgradeClusterStatusApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseGetUpgradeClusterStatusApi(client *core.CtyunClient) *CcseGetUpgradeClusterStatusApi {
	return &CcseGetUpgradeClusterStatusApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/upgrade/status",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseGetUpgradeClusterStatusApi) Do(ctx context.Context, credential core.Credential, req *CcseGetUpgradeClusterStatusRequest) (*CcseGetUpgradeClusterStatusResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	if req.TaskId != 0 {
		ctReq.AddParam("taskId", strconv.FormatInt(int64(req.TaskId), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseGetUpgradeClusterStatusResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseGetUpgradeClusterStatusRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	TaskId int64 /*  任务ID，您可以在升级集群接口获取任务ID
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=43&api=18045&data=128&isNormal=1&vid=121">升级集群</a>  */
}

type CcseGetUpgradeClusterStatusResponse struct {
	StatusCode int32                                           `json:"statusCode,omitempty"` /*  状态码  */
	Message    string                                          `json:"message,omitempty"`    /*  提示信息  */
	ReturnObj  []*CcseGetUpgradeClusterStatusReturnObjResponse `json:"returnObj"`            /*  返回对象  */
	Error      string                                          `json:"error,omitempty"`      /*  错误码  */
}

type CcseGetUpgradeClusterStatusReturnObjResponse struct {
	TaskId int64                                                `json:"taskId,omitempty"` /*  master0任务id  */
	Tasks  []*CcseGetUpgradeClusterStatusReturnObjTasksResponse `json:"tasks"`            /*  集群升级任务  */
}

type CcseGetUpgradeClusterStatusReturnObjTasksResponse struct {
	TaskId         int64  `json:"taskId,omitempty"`         /*  任务ID  */
	ClusterId      string `json:"clusterId,omitempty"`      /*  集群ID  */
	ParentTaskId   int64  `json:"parentTaskId,omitempty"`   /*  上级任务ID  */
	NodeName       string `json:"nodeName,omitempty"`       /*  升级节点名称  */
	NodeType       string `json:"nodeType,omitempty"`       /*  升级节点类型：master0, master, node  */
	MasterPlanHash string `json:"masterPlanHash,omitempty"` /*  升级master使用crd的latest hash  */
	WorkerPlanHash string `json:"workerPlanHash,omitempty"` /*  升级worker使用crd的latest hash  */
	TaskStatus     string `json:"taskStatus,omitempty"`     /*  升级状态：start-开始升级，end-完成升级，pause-暂停升级  */
	Version        string `json:"version,omitempty"`        /*  升级版本  */
	Status         int32  `json:"status,omitempty"`         /*  状态  0-无效 1-有效 2-删除  */
	CreatedBy      int64  `json:"createdBy,omitempty"`      /*  创建人  */
	CreatedTime    string `json:"createdTime,omitempty"`    /*  创建时间  */
	UpdatedTime    string `json:"updatedTime,omitempty"`    /*  修改时间  */
}
