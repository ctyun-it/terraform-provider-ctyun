package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupDeleteEbsBackupApi
/* 删除云硬盘备份。
 */type EbsbackupDeleteEbsBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupDeleteEbsBackupApi(client *core.CtyunClient) *EbsbackupDeleteEbsBackupApi {
	return &EbsbackupDeleteEbsBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/delete",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupDeleteEbsBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupDeleteEbsBackupRequest) (*EbsbackupDeleteEbsBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupDeleteEbsBackupRequest
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
	var resp EbsbackupDeleteEbsBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupDeleteEbsBackupRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	BackupID string `json:"backupID,omitempty"` /*  云硬盘备份ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10039050">查询云硬盘备份列表</a>获取  */
}

type EbsbackupDeleteEbsBackupResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）。  */
	Message     string                                     `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                     `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupDeleteEbsBackupReturnObjResponse `json:"returnObj"`   /*  无实际对象返回，值为{}  */
	ErrorCode   string                                     `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                     `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupDeleteEbsBackupReturnObjResponse struct {
	TaskID string `json:"taskID"` /*  删除任务ID。  */
}
