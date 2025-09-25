package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupRestoreBackupApi
/* 恢复云硬盘备份。
 */type EbsbackupRestoreBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupRestoreBackupApi(client *core.CtyunClient) *EbsbackupRestoreBackupApi {
	return &EbsbackupRestoreBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/restore-backup",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupRestoreBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupRestoreBackupRequest) (*EbsbackupRestoreBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupRestoreBackupRequest
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
	var resp EbsbackupRestoreBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupRestoreBackupRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	BackupID string `json:"backupID,omitempty"` /*  云硬盘备份ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10039050">查询云硬盘备份列表</a>获取  */
	DiskID   string `json:"diskID,omitempty"`   /*  云硬盘ID，您可以通过<a href="https://www.ctyun.cn/document/10027696/10096187">查询云硬盘列表</a>获取  */
}

type EbsbackupRestoreBackupResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message"`     /*  错误信息的英文描述  */
	Description string                                   `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupRestoreBackupReturnObjResponse `json:"returnObj"`   /*  无实际对象返回，值为{}  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                   `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupRestoreBackupReturnObjResponse struct {
	TaskID string `json:"taskID"` /*  恢复任务ID  */
}
