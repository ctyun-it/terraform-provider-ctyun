package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupCreateEbsBackupApi
/* 创建云硬盘备份
 */type EbsbackupCreateEbsBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupCreateEbsBackupApi(client *core.CtyunClient) *EbsbackupCreateEbsBackupApi {
	return &EbsbackupCreateEbsBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/create",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupCreateEbsBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupCreateEbsBackupRequest) (*EbsbackupCreateEbsBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupCreateEbsBackupRequest
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
	var resp EbsbackupCreateEbsBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupCreateEbsBackupRequest struct {
	VolumeID     string `json:"volumeID,omitempty"`     /*  云硬盘ID，您可以通过<a href="https://www.ctyun.cn/document/10027696/10041052">查询云硬盘列表</a>获取  */
	RegionID     string `json:"regionID,omitempty"`     /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	Description  string `json:"description,omitempty"`  /*  云硬盘备份描述  */
	RepositoryID string `json:"repositoryID,omitempty"` /*  备份存储库ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10092688">查询存储库列表</a>获取  */
	Name         string `json:"name,omitempty"`         /*  备份名称  */
}

type EbsbackupCreateEbsBackupResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message"`     /*  错误信息的英文描述  */
	Description string                                     `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupCreateEbsBackupReturnObjResponse `json:"returnObj"`   /*  返回业务对象  */
	ErrorCode   string                                     `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                     `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupCreateEbsBackupReturnObjResponse struct {
	RegionID          string `json:"regionID"`          /*  资源池ID  */
	BackupID          string `json:"backupID"`          /*  云硬盘备份ID  */
	BackupName        string `json:"backupName"`        /*  云硬盘备份名称  */
	Status            string `json:"status"`            /*  云硬盘备份状态，available-可用， error-失败，creating-创建中  */
	Size              int32  `json:"size"`              /*  云硬盘大小，单位GB  */
	UsedSize          int32  `json:"usedSize"`          /*  云硬盘备份大小，单位Byte  */
	Description       string `json:"description"`       /*  云硬盘备份描述  */
	RepositoryID      string `json:"repositoryID"`      /*  备份存储库ID  */
	RepositoryName    string `json:"repositoryName"`    /*  备份存储库名称  */
	CreatedDate       int32  `json:"CreatedDate"`       /*  创建时间  */
	UpdatedDate       int32  `json:"updatedDate"`       /*  更新时间  */
	FinishDate        int32  `json:"finishDate"`        /*  备份完成时间  */
	RestoreDate       int32  `json:"restoreDate"`       /*  使用该云硬盘备份恢复数据时间  */
	RestoreFinishDate int32  `json:"restoreFinishDate"` /*  使用该云硬盘备份恢复完成时间  */
	Freeze            *bool  `json:"freeze"`            /*  是否冻结  */
	VolumeID          string `json:"volumeID"`          /*  云硬盘ID  */
	VolumeName        string `json:"volumeName"`        /*  云硬盘名称  */
	Entrypted         *bool  `json:"entrypted"`         /*  云硬盘是否加密  */
	CmkUUID           string `json:"cmkUUID"`           /*  云硬盘加密密钥UUID  */
	VolumeType        string `json:"volumeType"`        /*  云硬盘类型  */
	AzName            string `json:"azName"`            /*  可用域  */
	Pass              *bool  `json:"pass"`              /*  是否支持PASS  */
	VMID              string `json:"vMID"`              /*  挂载的云主机ID  */
	VMName            string `json:"vMName"`            /*  挂载的云主机名称  */
	ProjectID         string `json:"projectID"`         /*  企业项目ID  */
	TaskID            string `json:"taskID"`            /*  备份任务ID  */
}
