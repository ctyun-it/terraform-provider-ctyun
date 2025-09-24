package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListEbsBackupRepoApi
/* 查询云硬盘备份存储库列表.
 */type EbsbackupListEbsBackupRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListEbsBackupRepoApi(client *core.CtyunClient) *EbsbackupListEbsBackupRepoApi {
	return &EbsbackupListEbsBackupRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/repo/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListEbsBackupRepoApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListEbsBackupRepoRequest) (*EbsbackupListEbsBackupRepoResponse, error) {
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
	var resp EbsbackupListEbsBackupRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListEbsBackupRepoRequest struct {
	RegionID       string `json:"regionID,omitempty"`       /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  云硬盘备份存储库名称。  */
	RepositoryID   string `json:"repositoryID,omitempty"`   /*  云硬盘备份存储库ID。  */
	Status         string `json:"status,omitempty"`         /*  云硬盘备份存储库状态，取值范围：
	●active：可用。
	●master_order_creating：主订单未完成。
	●freezing：已冻结。
	●expired：已过期。  */
	HideExpire   *bool  `json:"hideExpire"`             /*  是否隐藏过期的云硬盘备份存储库。  */
	QueryContent string `json:"queryContent,omitempty"` /*  目前，仅支持备份存储库名称的过滤。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，默认1。  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目 ,默认10。  */
	Asc          *bool  `json:"asc"`                    /*  和sort配合使用，是否升序排列。  */
	Sort         string `json:"sort,omitempty"`         /*  和asc配合使用，指定用于排序的字段。可选字段：createdDate/expiredDate/size/freeSize/usedSize/name  */
	ProjectID    string `json:"projectID,omitempty"`    /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListEbsBackupRepoResponse struct {
	StatusCode  int32                                        `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                       `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                       `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupListEbsBackupRepoReturnObjResponse `json:"returnObj"`   /*  返回对象数组。  */
	ErrorCode   string                                       `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                       `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupListEbsBackupRepoReturnObjResponse struct{}
