package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupCreateBackupApi
/* 创建云硬盘备份。
 */type EbsbackupCreateBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupCreateBackupApi(client *core.CtyunClient) *EbsbackupCreateBackupApi {
	return &EbsbackupCreateBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/create-backup",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupCreateBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupCreateBackupRequest) (*EbsbackupCreateBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupCreateBackupRequest
	}{
		req,
	}, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupCreateBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupCreateBackupRequest struct {
	DiskID       string `json:"diskID,omitempty"`       /*  云硬盘ID，您可以通过<a href="https://www.ctyun.cn/document/10027696/10096187">查询云硬盘列表</a>获取。  */
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	Description  string `json:"description,omitempty"`  /*  云硬盘备份描述。  */
	RepositoryID string `json:"repositoryID,omitempty"` /*  备份存储库ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10039480">查询存储库列表</a>获取。  */
	BackupName   string `json:"backupName,omitempty"`   /*  备份名称，长度为 2~63 个字符，只能由数字、字母、-、_ 组成，不能以数字、-、_ 开头。  */
	FullBackup   *bool  `json:"fullBackup"`             /*  是否启用全量备份，取值范围：
	●true：是
	●false：否
	若启用该参数，则此次备份的类型为全量备份。
	某些特殊情况下，如您是第一次备份，或者切换了存储库，则本次备份为全量备份，不受该参数影响。  */
}

type EbsbackupCreateBackupResponse struct {
	StatusCode  int32                                   `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                  `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                  `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupCreateBackupReturnObjResponse `json:"returnObj"`   /*  成功时返回对象。  */
	ErrorCode   string                                  `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                  `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupCreateBackupReturnObjResponse struct {
	RegionID            string `json:"regionID"`            /*  资源池ID。  */
	BackupID            string `json:"backupID"`            /*  云硬盘备份ID。  */
	BackupName          string `json:"backupName"`          /*  云硬盘备份名称。  */
	BackupStatus        string `json:"backupStatus"`        /*  云硬盘备份状态，该接口会返回creating状态。  */
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
	DiskType            string `json:"diskType"`            /*  云硬盘类型，取值范围为：
	●SATA：普通IO。
	●SAS：高IO。
	●SSD：超高IO。
	●FAST-SSD：极速型SSD。
	●XSSD-0、XSSD-1、XSSD-2：X系列云硬盘。  */
	Paas         *bool  `json:"paas"`         /*  是否支持PAAS。  */
	InstanceID   string `json:"instanceID"`   /*  云硬盘挂载的云主机ID。  */
	InstanceName string `json:"instanceName"` /*  云硬盘挂载的云主机名称。  */
	ProjectID    string `json:"projectID"`    /*  企业项目ID。  */
	TaskID       string `json:"taskID"`       /*  云硬盘备份任务ID。  */
}
