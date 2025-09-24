package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListEbsBackupPolicyTasksApi
/* 查询备份策略创建的备份任务列表
 */type EbsbackupListEbsBackupPolicyTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListEbsBackupPolicyTasksApi(client *core.CtyunClient) *EbsbackupListEbsBackupPolicyTasksApi {
	return &EbsbackupListEbsBackupPolicyTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/policy/list-tasks",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListEbsBackupPolicyTasksApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListEbsBackupPolicyTasksRequest) (*EbsbackupListEbsBackupPolicyTasksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("policyID", req.PolicyID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.Asc != nil {
		ctReq.AddParam("asc", strconv.FormatBool(*req.Asc))
	}
	if req.Sort != "" {
		ctReq.AddParam("sort", req.Sort)
	}
	if req.TaskStatus != 0 {
		ctReq.AddParam("taskStatus", strconv.FormatInt(int64(req.TaskStatus), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListEbsBackupPolicyTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListEbsBackupPolicyTasksRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID   string `json:"policyID,omitempty"`   /*  备份策略ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10040084">查询备份策略</a>获取  */
	PageNo     int32  `json:"pageNo,omitempty"`     /*  页码，默认值1  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  每页记录数目 ,默认10  */
	Asc        *bool  `json:"asc"`                  /*  和sort配合使用，是否升序排列，默认降序  */
	Sort       string `json:"sort,omitempty"`       /*  和asc配合使用，指定用于排序的字段。可选字段：createdTime/completedTime，默认createdTime  */
	TaskStatus int32  `json:"taskStatus,omitempty"` /*  备份任务状态，-1-失败，0-执行中，1-成功  */
}

type EbsbackupListEbsBackupPolicyTasksResponse struct {
	StatusCode  int32                                               `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                              `json:"message"`     /*  错误信息的英文描述  */
	Description string                                              `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupListEbsBackupPolicyTasksReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                              `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                              `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupListEbsBackupPolicyTasksReturnObjResponse struct {
	TaskList     []*EbsbackupListEbsBackupPolicyTasksReturnObjTaskListResponse `json:"taskList"`     /*  备份任务列表  */
	TotalCount   int32                                                         `json:"totalCount"`   /*  备份任务总数  */
	CurrentCount int32                                                         `json:"currentCount"` /*  当前页备份任务数  */
}

type EbsbackupListEbsBackupPolicyTasksReturnObjTaskListResponse struct {
	TaskID        string `json:"taskID"`        /*  备份任务ID  */
	TaskStatus    int32  `json:"taskStatus"`    /*  备份任务状态，-1-失败，0-执行中，1-成功  */
	DiskID        string `json:"diskID"`        /*  云硬盘ID  */
	DiskName      string `json:"diskName"`      /*  云硬盘名称  */
	BackupName    string `json:"backupName"`    /*  云硬盘备份名称  */
	CreatedTime   int32  `json:"createdTime"`   /*  创建时间  */
	CompletedTime int32  `json:"completedTime"` /*  完成时间  */
}
