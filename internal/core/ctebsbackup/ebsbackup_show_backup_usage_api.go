package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupShowBackupUsageApi
/* 查询云硬盘备份实际占用存储大小
 */type EbsbackupShowBackupUsageApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupShowBackupUsageApi(client *core.CtyunClient) *EbsbackupShowBackupUsageApi {
	return &EbsbackupShowBackupUsageApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/show-usage",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupShowBackupUsageApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupShowBackupUsageRequest) (*EbsbackupShowBackupUsageResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("backupID", req.BackupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupShowBackupUsageResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupShowBackupUsageRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	BackupID string `json:"backupID,omitempty"` /*  云硬盘备份ID  */
}

type EbsbackupShowBackupUsageResponse struct {
	StatusCode  int32                                      `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                     `json:"message"`     /*  错误信息的英文描述  */
	Description string                                     `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupShowBackupUsageReturnObjResponse `json:"returnObj"`   /*  是  */
	ErrorCode   string                                     `json:"errorCode"`   /*  业务错误细分码，发生错误时返回  */
	Error       string                                     `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupShowBackupUsageReturnObjResponse struct {
	BackupSize int32 `json:"backupSize"` /*  云硬盘备份大小，单位Byte  */
}
