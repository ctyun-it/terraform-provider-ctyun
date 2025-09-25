package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtecsListInstanceBackupPolicyApi
/* 查询云主机备份策略列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权
 */type CtecsListInstanceBackupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceBackupPolicyApi(client *core.CtyunClient) *CtecsListInstanceBackupPolicyApi {
	return &CtecsListInstanceBackupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup-policy/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceBackupPolicyApi) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceBackupPolicyRequest) (*CtecsListInstanceBackupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PolicyID != "" {
		ctReq.AddParam("policyID", req.PolicyID)
	}
	if req.PolicyName != "" {
		ctReq.AddParam("policyName", req.PolicyName)
	}
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListInstanceBackupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceBackupPolicyRequest struct {
	RegionID   string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID   string /*  云主机备份策略ID，32字节<br />获取：<br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6914&data=87&isNormal=1&vid=81">创建云主机备份策略</a>  */
	PolicyName string /*  云主机备份策略名称。满足以下规则：只能由数字、英文字母、中划线-、下划线_、点.组成，长度为2-64字符<br />注：在所有资源池不可重复。模糊过滤，存在policyID过滤参数时此参数无效  */
	ProjectID  string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	PageNo     int32  /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize   int32  /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
}

type CtecsListInstanceBackupPolicyResponse struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  错误信息的英文描述  */
	Description string                                          `json:"description,omitempty"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *CtecsListInstanceBackupPolicyReturnObjResponse `json:"returnObj"`             /*  返回对象  */
}

type CtecsListInstanceBackupPolicyReturnObjResponse struct {
	CurrentCount int32                                                       `json:"currentCount,omitempty"` /*  当前页记录数目  */
	CurrentPage  int32                                                       `json:"currentPage,omitempty"`  /*  当前页码  */
	TotalCount   int32                                                       `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                       `json:"totalPage,omitempty"`    /*  总页数  */
	PolicyList   []*CtecsListInstanceBackupPolicyReturnObjPolicyListResponse `json:"policyList"`             /*  策略列表  */
}

type CtecsListInstanceBackupPolicyReturnObjPolicyListResponse struct {
	RegionID      string `json:"regionID,omitempty"`      /*  资源池ID，您可以调用[regionID](https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87)查看最新的天翼云资源池列表  */
	Status        int32  `json:"status,omitempty"`        /*  是否启用策略，取值范围：<br />0：停用，<br />1：启用  */
	PolicyID      string `json:"policyID,omitempty"`      /*  云主机备份策略ID  */
	PolicyName    string `json:"policyName,omitempty"`    /*  云主机备份策略名称  */
	CycleType     string `json:"cycleType,omitempty"`     /*  云主机备份周期类型，取值范围： <br />day：按天备份<br />week：按星期备份  */
	CycleDay      int32  `json:"cycleDay,omitempty"`      /*  只有cycleType为day时返回备份周期值  */
	CycleWeek     string `json:"cycleWeek,omitempty"`     /*  只有cycleType为week时返回备份周期，取值范围：0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开  */
	Time          string `json:"time,omitempty"`          /*  备份整点时间，取值范围：0-23，如果一天内多个时间节点备份，以逗号隔开  */
	RetentionType string `json:"retentionType,omitempty"` /*  云主机备份保留类型，取值范围：<br />date：按时间保留，<br />num：按数量保留，<br />all：永久保留  */
	RetentionNum  int32  `json:"retentionNum,omitempty"`  /*  只有retentionType为num时返回保留数量值  */
	RetentionDay  int32  `json:"retentionDay,omitempty"`  /*  只有retentionType为date时返回保留天数值  */
	ResourceCount int32  `json:"resourceCount,omitempty"` /*  策略已绑定的云主机数量  */
	//TODO  openapi返回字段多了个空格，待修复后修改
	ResourceIDs        string                                                                    `json:"resourceIDs ,omitempty"`       /*  策略已绑定的云主机ID，以逗号分隔  */
	RepositoryList     []*CtecsListInstanceBackupPolicyReturnObjPolicyListRepositoryListResponse `json:"repositoryList"`               /*  策略已绑定的云主机备份库列表  */
	ProjectID          string                                                                    `json:"projectID,omitempty"`          /*  企业项目ID  */
	FullBackupInterval int32                                                                     `json:"fullBackupInterval,omitempty"` /*  是否启用周期性全量备份。-1代表不开启，默认为-1；取值范围为[-1,100]，即每执行n次增量备份后，执行一次全量备份；若传入为0，代表每一次均为全量备份。  */
	AdvRetentionStatus bool                                                                      `json:"advRetentionStatus,omitempty"` /*  是否开启高级保留策略，false（不启用），true(启用)，默认值为false。需校验云主机备份保留类型（retentionType），若保留类型为按数量保存（num），可开启高级保留策略；若保留类型为date（按时间保存）或all（永久保存），不可开启高级保留策略。  */
	AdvRetention       *AdvRetention                                                             `json:"advRetention,omitempty"`       /*  高级保留策略内容，只有retentionType为num且advRetentionStatus为true才生效  */
}

type CtecsListInstanceBackupPolicyReturnObjPolicyListRepositoryListResponse struct {
	RepositoryID   string `json:"repositoryID,omitempty"`   /*  云主机备份库ID  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  云主机备份库名称  */
}
