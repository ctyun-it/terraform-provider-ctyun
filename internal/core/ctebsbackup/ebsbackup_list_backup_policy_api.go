package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListBackupPolicyApi
/* 查询云硬盘备份策略列表。
 */type EbsbackupListBackupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListBackupPolicyApi(client *core.CtyunClient) *EbsbackupListBackupPolicyApi {
	return &EbsbackupListBackupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/policy/list-policies",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListBackupPolicyApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListBackupPolicyRequest) (*EbsbackupListBackupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PolicyID != "" {
		ctReq.AddParam("policyID", req.PolicyID)
	}
	if req.PolicyName != "" {
		ctReq.AddParam("policyName", req.PolicyName)
	}
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListBackupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListBackupPolicyRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PageNo     int32  `json:"pageNo,omitempty"`     /*  页码，默认值1  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  每页记录数目 ,默认10  */
	PolicyID   string `json:"policyID,omitempty"`   /*  备份策略ID  */
	PolicyName string `json:"policyName,omitempty"` /*  备份策略名，指定了policyID时，该参数会被忽略  */
	ProjectID  string `json:"projectID,omitempty"`  /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListBackupPolicyResponse struct {
	StatusCode  int32                                       `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                      `json:"message"`     /*  错误信息的英文描述  */
	Description string                                      `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupListBackupPolicyReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                      `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                      `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupListBackupPolicyReturnObjResponse struct {
	PolicyList   []*EbsbackupListBackupPolicyReturnObjPolicyListResponse `json:"policyList"`   /*  策略列表  */
	TotalCount   int32                                                   `json:"totalCount"`   /*  策略总数  */
	CurrentCount int32                                                   `json:"currentCount"` /*  当前页策略数  */
}

type EbsbackupListBackupPolicyReturnObjPolicyListResponse struct {
	RegionID              string                                                                `json:"regionID"`              /*  资源池ID  */
	AccountID             string                                                                `json:"accountID"`             /*  账户ID  */
	Status                int32                                                                 `json:"status"`                /*  状态，0-停用，1-启用  */
	CreatedTime           int32                                                                 `json:"createdTime"`           /*  创建时间  */
	PolicyID              string                                                                `json:"policyID"`              /*  策略ID  */
	PolicyName            string                                                                `json:"policyName"`            /*  策略名  */
	CycleType             string                                                                `json:"cycleType"`             /*  备份周期类型，day-按天备份，week-按星期备份  */
	CycleDay              int32                                                                 `json:"cycleDay"`              /*  备份周期，只有cycleType为day时返回  */
	CycleWeek             string                                                                `json:"cycleWeek"`             /*  备份周期，只有cycleType为week时返回，则取值范围0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开  */
	Time                  string                                                                `json:"time"`                  /*  备份整点时间，取值范围0-23，如果一天内多个时间节点备份，以逗号隔开  */
	RetentionType         string                                                                `json:"retentionType"`         /*  备份保留类型，num-按数量保留，date-按时间保留，all-全部保留  */
	RetentionNum          int32                                                                 `json:"retentionNum"`          /*  保留数量，只有retentionType为num时返回  */
	RetentionDay          int32                                                                 `json:"retentionDay"`          /*  保留天数，只有retentionType为date时返回  */
	RemainFirstOfCurMonth *bool                                                                 `json:"remainFirstOfCurMonth"` /*  是否保留每个月第一个备份，在retentionType为num时返回  */
	BindedDiskCount       int32                                                                 `json:"bindedDiskCount"`       /*  策略绑定的云硬盘数量  */
	BindedDiskIDs         string                                                                `json:"bindedDiskIDs"`         /*  策略绑定的云硬盘ID，以逗号分隔  */
	RepositoryList        []*CtebsListDiskBackupPolicyReturnObjPolicyListRepositoryListResponse `json:"repositoryList"`        /*  策略绑定的云硬盘备份存储库列表  */
	ProjectID             string                                                                `json:"projectID"`             /*  企业项目ID  */
	FullBackupInterval    int32                                                                 `json:"fullBackupInterval"`    /*  启用周期性全量备份。-1代表不开启，默认为-1；取值范围为[-1,100]，即每执行n次增量备份后，执行一次全量备份。  */
	AdvRetentionStatus    *bool                                                                 `json:"advRetentionStatus"`    /*  是否启用高级保留策略，取值范围：
	●true：启用
	●false：不启用  */
	AdvRetention *EbsbackupListBackupPolicyReturnObjPolicyListAdvRetentionResponse `json:"advRetention"` /*  高级保留策略内容。  */
}

type EbsbackupListBackupPolicyReturnObjPolicyListAdvRetentionResponse struct {
	AdvDay int32 `json:"advDay"` /*  ● 保留n天内，每天最新的一个备份。
	● 单位为天，取值范围：[0, 100]，默认值0  */
	AdvWeek int32 `json:"advWeek"` /*  ● 保留n周内，每周最新的一个备份。
	● 单位为周，取值范围：[0, 100]，默认值0  */
	AdvMonth int32 `json:"advMonth"` /*  ● 保留n月内，每月最新的一个备份。
	● 单位为月，取值范围：[0, 100]，默认值0  */
	AdvYear int32 `json:"advYear"` /*  ● 保留n年内，每年最新的一个备份。
	● 单位为年，取值范围：[0, 100]，默认值0  */
}

type CtebsListDiskBackupPolicyReturnObjPolicyListRepositoryListResponse struct {
	RepositoryID   string `json:"repositoryID,omitempty"`   /*  云主机备份库ID  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  云主机备份库名称  */
}
