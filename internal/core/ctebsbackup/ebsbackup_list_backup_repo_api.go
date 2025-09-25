package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListBackupRepoApi
/* 查询云硬盘备份存储库列表。
 */type EbsbackupListBackupRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListBackupRepoApi(client *core.CtyunClient) *EbsbackupListBackupRepoApi {
	return &EbsbackupListBackupRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/repo/list-repos",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListBackupRepoApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListBackupRepoRequest) (*EbsbackupListBackupRepoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.RepositoryName != "" {
		ctReq.AddParam("repositoryName", req.RepositoryName)
	}
	if req.RepositoryID != "" {
		ctReq.AddParam("repositoryID", req.RepositoryID)
	}
	if req.Status != "" {
		ctReq.AddParam("status", req.Status)
	}
	if req.HideExpire != nil {
		ctReq.AddParam("hideExpire", strconv.FormatBool(*req.HideExpire))
	}
	if req.QueryContent != "" {
		ctReq.AddParam("queryContent", req.QueryContent)
	}
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
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListBackupRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListBackupRepoRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  云硬盘备份存储库名称。  */
	RepositoryID   string `json:"repositoryID,omitempty"`   /*  云硬盘备份存储库ID。  */
	Status         string `json:"status,omitempty"`         /*  云硬盘备份存储库状态，取值范围：
	●active：可用。
	●master_order_creating：主订单未完成。
	●freezing：已冻结。
	●expired：已过期。  */
	HideExpire   *bool  `json:"hideExpire"`             /*  是否隐藏过期的云硬盘备份存储库。  */
	QueryContent string `json:"queryContent,omitempty"` /*  模糊查询，目前仅支持备份存储库名称的过滤。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，默认为1。  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目 ,默认为10。  */
	Asc          *bool  `json:"asc"`                    /*  和sort配合使用，是否是升序排列。  */
	Sort         string `json:"sort,omitempty"`         /*  和asc配合使用，指定用于排序的字段。可选字段：
	●createdTime：创建时间。
	●expiredTime：过期时间。
	●size：存储库空间大小。
	●freeSize：存储库剩余空间。
	●usedSize：存储库已使用空间大小。
	●repositoryName：存储库名称。  */
	ProjectID string `json:"projectID,omitempty"` /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListBackupRepoResponse struct {
	StatusCode  int32                                     `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                    `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                    `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupListBackupRepoReturnObjResponse `json:"returnObj"`   /*  返回对象数组。  */
	ErrorCode   string                                    `json:"errorCode"`   /*  业务错误细分码，发生错误时返回。  */
	Error       string                                    `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupListBackupRepoReturnObjResponse struct {
	RepositoryList []*EbsbackupListBackupRepoReturnObjRepositoryListResponse `json:"repositoryList"` /*  云硬盘备份存储库列表。  */
	TotalCount     int32                                                     `json:"totalCount"`     /*  云硬盘备份存储库总数。  */
	CurrentCount   int32                                                     `json:"currentCount"`   /*  当前页云硬盘备份存储库数。  */
}

type EbsbackupListBackupRepoReturnObjRepositoryListResponse struct {
	RegionID       string `json:"regionID"`       /*  资源池ID。  */
	RepositoryID   string `json:"repositoryID"`   /*  备份存储库ID。  */
	RepositoryName string `json:"repositoryName"` /*  备份存储库名称。  */
	Status         string `json:"status"`         /*  云硬盘备份存储库状态，取值范围：
	●active：可用。
	●master_order_creating：主订单未完成。
	●freezing：已冻结。
	●expired：已过期。  */
	Size        int32    `json:"size"`        /*  云硬盘备份存储库总容量，单位为GB。  */
	FreeSize    float64  `json:"freeSize"`    /*  云硬盘备份存储库剩余大小，单位为GB。  */
	UsedSize    float64  `json:"usedSize"`    /*  云硬盘备份存储库使用大小，单位为GB。  */
	CreatedTime int32    `json:"createdTime"` /*  存储库创建时间。  */
	ExpiredTime int32    `json:"expiredTime"` /*  存储库到期时间。  */
	Expired     *bool    `json:"expired"`     /*  备份存储库是否到期。  */
	Freeze      *bool    `json:"freeze"`      /*  存储库是否冻结。  */
	Paas        *bool    `json:"paas"`        /*  是否支持PAAS。  */
	BackupList  []string `json:"backupList"`  /*  备份存储库下的可用的备份列表，元素为备份ID。  */
	ProjectID   string   `json:"projectID"`   /*  企业项目ID。  */
}
