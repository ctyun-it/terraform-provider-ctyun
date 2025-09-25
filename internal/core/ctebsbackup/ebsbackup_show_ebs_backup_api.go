package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupShowEbsBackupApi
/* 查询云硬盘备份信息。
 */type EbsbackupShowEbsBackupApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupShowEbsBackupApi(client *core.CtyunClient) *EbsbackupShowEbsBackupApi {
	return &EbsbackupShowEbsBackupApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/show",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupShowEbsBackupApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupShowEbsBackupRequest) (*EbsbackupShowEbsBackupResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("backupID", req.BackupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupShowEbsBackupResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupShowEbsBackupRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID。  */
	BackupID string `json:"backupID,omitempty"` /*  云硬盘备份ID。  */
}

type EbsbackupShowEbsBackupResponse struct {
	StatusCode  int32                                    `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                   `json:"message"`     /*  成功或失败时的描述，一般为英文描述。  */
	Description string                                   `json:"description"` /*  成功或失败时的描述，一般为中文描述。  */
	ReturnObj   *EbsbackupShowEbsBackupReturnObjResponse `json:"returnObj"`   /*  返回对象。  */
	ErrorCode   string                                   `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
	Error       string                                   `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码。  */
}

type EbsbackupShowEbsBackupReturnObjResponse struct{}
