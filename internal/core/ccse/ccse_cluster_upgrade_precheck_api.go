package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CcseClusterUpgradePrecheckApi
/* 调用该接口获取集群升级检查结果
 */type CcseClusterUpgradePrecheckApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCcseClusterUpgradePrecheckApi(client *core.CtyunClient) *CcseClusterUpgradePrecheckApi {
	return &CcseClusterUpgradePrecheckApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/cce/clusters/{clusterId}/upgrade/precheck",
			ContentType:  "application/json",
		},
	}
}

func (a *CcseClusterUpgradePrecheckApi) Do(ctx context.Context, credential core.Credential, req *CcseClusterUpgradePrecheckRequest) (*CcseClusterUpgradePrecheckResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder = builder.ReplaceUrl("clusterId", req.ClusterId)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("regionId", req.RegionId)
	ctReq.AddParam("checklistName", req.ChecklistName)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CcseClusterUpgradePrecheckResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CcseClusterUpgradePrecheckRequest struct {
	ClusterId string /*  集群ID，获取方式请参见<a href="https://www.ctyun.cn/document/10083472/11002105">如何获取接口URI中参数</a>。  */
	RegionId  string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10083472/11004422" target="_blank">云容器引擎资源池</a>
	另外您通过<a href="https://www.ctyun.cn/document/10026730/10028695" target="_blank">地域和可用区</a>来了解资源池
	获取：
	<span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&vid=81" target="_blank">资源池列表查询</a>
	*/
	ChecklistName string /*  集群升级检查项，目前支持的检查项及含义如下
	ClusterComponentChecklist：集群组件检查
	ClusterConfigurationChecklist：集群配置检查
	ClusterResourceChecklist：集群资源检查  */
}

type CcseClusterUpgradePrecheckResponse struct {
	StatusCode int32                                        `json:"statusCode,omitempty"` /*  响应状态码  */
	Message    string                                       `json:"message,omitempty"`    /*  响应信息  */
	ReturnObj  *CcseClusterUpgradePrecheckReturnObjResponse `json:"returnObj"`            /*  响应对象  */
	Error      string                                       `json:"error,omitempty"`      /*  错误码，参见错误码说明  */
}

type CcseClusterUpgradePrecheckReturnObjResponse struct {
	ErrorMessage     string                                                         `json:"errorMessage,omitempty"` /*  集群升级中的错误信息  */
	CheckStatus      string                                                         `json:"checkStatus,omitempty"`  /*  检查状态：checking 检查中；completed 检查完成  */
	CheckEntryResult []*CcseClusterUpgradePrecheckReturnObjCheckEntryResultResponse `json:"checkEntryResult"`       /*  检查结果  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultResponse struct {
	EntryName            string                                                                             `json:"entryName,omitempty"`      /*  检查项名称  */
	EntryGroupName       string                                                                             `json:"entryGroupName,omitempty"` /*  检查项所属分组  */
	NormalInstanceResult []*CcseClusterUpgradePrecheckReturnObjCheckEntryResultNormalInstanceResultResponse `json:"normalInstanceResult"`     /*  正常节点实例检查结果  */
	ErrorInstanceResult  []*CcseClusterUpgradePrecheckReturnObjCheckEntryResultErrorInstanceResultResponse  `json:"errorInstanceResult"`      /*  异常节点实例检查结果  */
	CheckSummary         *CcseClusterUpgradePrecheckReturnObjCheckEntryResultCheckSummaryResponse           `json:"checkSummary"`             /*  检查汇总  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultNormalInstanceResultResponse struct {
	InstanceId        string                                                                                            `json:"instanceId,omitempty"`   /*  节点实例id  */
	InstanceName      string                                                                                            `json:"instanceName,omitempty"` /*  节点实例名称  */
	CheckpointResults *CcseClusterUpgradePrecheckReturnObjCheckEntryResultNormalInstanceResultCheckpointResultsResponse `json:"checkpointResults"`      /*  检查结果  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultErrorInstanceResultResponse struct {
	InstanceId        string                                                                                           `json:"instanceId,omitempty"`   /*  节点实例id  */
	InstanceName      string                                                                                           `json:"instanceName,omitempty"` /*  节点实例名称  */
	CheckpointResults *CcseClusterUpgradePrecheckReturnObjCheckEntryResultErrorInstanceResultCheckpointResultsResponse `json:"checkpointResults"`      /*  检查结果  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultCheckSummaryResponse struct {
	NormalCount int32  `json:"normalCount,omitempty"` /*  正常数量  */
	ErrorCount  int32  `json:"errorCount,omitempty"`  /*  错误数量  */
	Process     string `json:"process,omitempty"`     /*  检查进度  */
	Code        string `json:"code,omitempty"`        /*  结果代码  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultNormalInstanceResultCheckpointResultsResponse struct {
	MessageLevel   string `json:"messageLevel,omitempty"`   /*  消息级别  */
	CheckpointName string `json:"checkpointName,omitempty"` /*  检查内容  */
	AdviseCode     string `json:"adviseCode,omitempty"`     /*  修复建议  */
	MessageCode    string `json:"messageCode,omitempty"`    /*  结果  */
	AffectCode     string `json:"affectCode,omitempty"`     /*  影响  */
}

type CcseClusterUpgradePrecheckReturnObjCheckEntryResultErrorInstanceResultCheckpointResultsResponse struct {
	MessageLevel   string `json:"messageLevel,omitempty"`   /*  消息级别  */
	CheckpointName string `json:"checkpointName,omitempty"` /*  检查内容  */
	AdviseCode     string `json:"adviseCode,omitempty"`     /*  修复建议  */
	MessageCode    string `json:"messageCode,omitempty"`    /*  结果  */
	AffectCode     string `json:"affectCode,omitempty"`     /*  影响  */
}
