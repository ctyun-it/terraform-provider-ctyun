package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListBackupTaskApi
/* 无
 */type CtecsListBackupTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListBackupTaskApi(client *core.CtyunClient) *CtecsListBackupTaskApi {
	return &CtecsListBackupTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup-task/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListBackupTaskApi) Do(ctx context.Context, credential core.Credential, req *CtecsListBackupTaskRequest) (*CtecsListBackupTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListBackupTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListBackupTaskRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	TaskID       string `json:"taskID,omitempty"`       /*  任务ID  */
	TaskType     string `json:"taskType,omitempty"`     /*  任务类型，取值范围：create（生成备份副本任务），restore（恢复备份数据任务），delete（删除备份副本任务）<br />注：不传默认全部  */
	TaskStatus   string `json:"taskStatus,omitempty"`   /*  任务状态，取值范围：successed（成功），failed（失败），in-progress（执行中），canceling（取消中），canceled（已取消）<br />注：不传默认全部  */
	InstanceID   string `json:"instanceID,omitempty"`   /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
	BackupID     string `json:"backupID,omitempty"`     /*  备份ID  */
	BackupName   string `json:"backupName,omitempty"`   /*  备份名称  */
	QueryContent string `json:"queryContent,omitempty"` /*  模糊查询，可匹配查询字段：任务ID、备份副本ID、备份副本名称、云主机ID、云主机名称、存储库ID、存储库名称  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListBackupTaskResponse struct {
	StatusCode  int32                                 `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsListBackupTaskReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListBackupTaskReturnObjResponse struct {
	CurrentCount int32                                          `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                          `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                          `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListBackupTaskReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListBackupTaskReturnObjResultsResponse struct {
	TaskID         string `json:"taskID,omitempty"`         /*  任务ID  */
	TaskStatus     string `json:"taskStatus,omitempty"`     /*  任务状态  */
	TaskType       string `json:"taskType,omitempty"`       /*  任务类型  */
	InstanceID     string `json:"instanceID,omitempty"`     /*  云主机ID  */
	InstanceName   string `json:"instanceName,omitempty"`   /*  云主机名称  */
	BackupID       string `json:"backupID,omitempty"`       /*  备份ID  */
	BackupName     string `json:"backupName,omitempty"`     /*  备份名称  */
	RepositoryID   string `json:"repositoryID,omitempty"`   /*  存储库ID  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  存储库名称  */
	StrategyID     string `json:"strategyID,omitempty"`     /*  策略的ID，非策略触发执行时为空  */
	StrategyName   string `json:"strategyName,omitempty"`   /*  策略的名称，非策略触发执行时为空  */
	StartTime      string `json:"startTime,omitempty"`      /*  开始时间  */
	FinishTime     string `json:"finishTime,omitempty"`     /*  完成时间  */
	TaskDetailDesc string `json:"taskDetailDesc,omitempty"` /*  任务详情  */
}
