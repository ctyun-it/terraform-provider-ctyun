package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListEbsBackupApi
/* 查询云硬盘备份列表.
 */type EbsbackupListEbsBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListEbsBackupApi(client *core.CtyunClient) *EbsbackupListEbsBackupApi {
	return &EbsbackupListEbsBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListEbsBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListEbsBackupRequest) (*EbsbackupListEbsBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.VolumeID != "" {
		ctReq.AddParam("volumeID", req.VolumeID)
	}
	if req.VolumeName != "" {
		ctReq.AddParam("volumeName", req.VolumeName)
	}
	if req.BackupName != "" {
		ctReq.AddParam("backupName", req.BackupName)
	}
	if req.RepositoryID != "" {
		ctReq.AddParam("repositoryID", req.RepositoryID)
	}
	if req.Status != "" {
		ctReq.AddParam("status", req.Status)
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
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListEbsBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListEbsBackupRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	VolumeID     string `json:"volumeID,omitempty"`     /*  云硬盘ID。  */
	VolumeName   string `json:"volumeName,omitempty"`   /*  云硬盘名称，模糊过滤。  */
	BackupName   string `json:"backupName,omitempty"`   /*  云硬盘备份名称，模糊过滤。  */
	RepositoryID string `json:"repositoryID,omitempty"` /*  云硬盘备份存储库ID。  */
	Status       string `json:"status,omitempty"`       /*  云硬盘备份状态，取值范围：
	●available：可用。
	●error：错误。
	●restoring：恢复中。
	●creating：创建中。
	●deleting：删除中。
	●merging_backup：合并中。
	●frozen：已冻结。  */
	QueryContent string `json:"queryContent,omitempty"` /*  该参数，可用于模糊过滤 云硬盘ID/云硬盘名称/备份ID/备份名称，即上述4个字段如果包含该参数的值，则会被过滤出来。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，默认1。  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目 ,默认10。  */
	ProjectID    string `json:"projectID,omitempty"`    /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListEbsBackupResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                   `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                   `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupListEbsBackupReturnObjResponse `json:"returnObj"`   /*  返回对象。  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                   `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupListEbsBackupReturnObjResponse struct{}
