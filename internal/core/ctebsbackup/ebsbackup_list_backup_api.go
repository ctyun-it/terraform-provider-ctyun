package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListBackupApi
/* 查询云硬盘备份列表。
 */type EbsbackupListBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListBackupApi(client *core.CtyunClient) *EbsbackupListBackupApi {
	return &EbsbackupListBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/list-backups",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListBackupRequest) (*EbsbackupListBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.DiskID != "" {
		ctReq.AddParam("diskID", req.DiskID)
	}
	if req.DiskName != "" {
		ctReq.AddParam("diskName", req.DiskName)
	}
	if req.BackupName != "" {
		ctReq.AddParam("backupName", req.BackupName)
	}
	if req.RepositoryID != "" {
		ctReq.AddParam("repositoryID", req.RepositoryID)
	}
	if req.BackupStatus != "" {
		ctReq.AddParam("backupStatus", req.BackupStatus)
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
	var resp EbsbackupListBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListBackupRequest struct {
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	DiskID       string `json:"diskID,omitempty"`       /*  云硬盘ID。  */
	DiskName     string `json:"diskName,omitempty"`     /*  云硬盘名称，模糊过滤。  */
	BackupName   string `json:"backupName,omitempty"`   /*  云硬盘备份名称，模糊过滤。  */
	RepositoryID string `json:"repositoryID,omitempty"` /*  云硬盘备份存储库ID。  */
	BackupStatus string `json:"backupStatus,omitempty"` /*  云硬盘备份状态，取值范围：
	●available：可用。
	●error：错误。
	●restoring：恢复中。
	●creating：创建中。
	●deleting：删除中。
	●merging_backup：合并中。
	●frozen：已冻结。  */
	QueryContent string `json:"queryContent,omitempty"` /*  该参数，可用于模糊过滤 云硬盘ID/云硬盘名称/备份ID/备份名称，即上述4个字段如果包含该参数的值，则会被过滤出来。  */
	PageNo       int32  `json:"pageNo,omitempty"`       /*  页码，默认1。  */
	PageSize     int32  `json:"pageSize,omitempty"`     /*  每页记录数目，默认10。  */
	ProjectID    string `json:"projectID,omitempty"`    /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListBackupResponse struct {
	StatusCode  int32                                 `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupListBackupReturnObjResponse `json:"returnObj"`   /*  返回对象。  */
	ErrorCode   string                                `json:"errorCode"`   /*  参考错误码。  */
	Error       string                                `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupListBackupReturnObjResponse struct {
	BackupList   []*EbsbackupListBackupReturnObjBackupListResponse `json:"backupList"`   /*  云硬盘备份列表。  */
	TotalCount   int32                                             `json:"totalCount"`   /*  云硬盘备份总数。  */
	CurrentCount int32                                             `json:"currentCount"` /*  当前页云硬盘备份数。  */
}

type EbsbackupListBackupReturnObjBackupListResponse struct {
	RegionID     string `json:"regionID"`     /*  资源池ID。  */
	BackupID     string `json:"backupID"`     /*  云硬盘备份ID。  */
	BackupName   string `json:"backupName"`   /*  云硬盘备份名称。  */
	BackupStatus string `json:"backupStatus"` /*  云硬盘备份状态，取值范围：
	●available：可用。
	●error：错误。
	●restoring：恢复中。
	●creating：创建中。
	●deleting：删除中。
	●merging_backup：合并中。
	●frozen：已冻结。  */
	DiskSize            int32  `json:"diskSize"`            /*  云硬盘大小，单位GB。  */
	BackupSize          int32  `json:"backupSize"`          /*  云硬盘备份大小，单位Byte。  */
	Description         string `json:"description"`         /*  云硬盘备份描述。  */
	RepositoryID        string `json:"repositoryID"`        /*  备份存储库ID。  */
	RepositoryName      string `json:"repositoryName"`      /*  备份存储库名称。  */
	CreatedTime         int32  `json:"createdTime"`         /*  备份创建时间。  */
	UpdatedTime         int32  `json:"updatedTime"`         /*  备份更新时间。  */
	FinishedTime        int32  `json:"finishedTime"`        /*  备份完成时间。  */
	RestoredTime        int32  `json:"restoredTime"`        /*  使用该云硬盘备份恢复数据时间。  */
	RestoreFinishedTime int32  `json:"restoreFinishedTime"` /*  使用该云硬盘备份恢复完成时间。  */
	Freeze              *bool  `json:"freeze"`              /*  备份是否冻结。  */
	DiskID              string `json:"diskID"`              /*  云硬盘ID。  */
	DiskName            string `json:"diskName"`            /*  云硬盘名称。  */
	Encrypted           *bool  `json:"encrypted"`           /*  云硬盘是否加密。  */
	DiskType            string `json:"diskType"`            /*  磁盘规格类型，取值为：
	●SATA：普通IO。
	●SAS：高IO。
	●SSD：超高IO。
	●FAST-SSD：极速型SSD。
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
	Paas         *bool  `json:"paas"`         /*  是否支持PAAS。  */
	InstanceID   string `json:"instanceID"`   /*  云硬盘挂载的云主机ID。  */
	InstanceName string `json:"instanceName"` /*  云硬盘挂载的云主机名称。  */
	ProjectID    string `json:"projectID"`    /*  企业项目ID。  */
	BackupType   string `json:"backupType"`   /*  备份类型，取值范围：
	●full-backup：全量备份
	●incremental-backup：增量备份  */
}
