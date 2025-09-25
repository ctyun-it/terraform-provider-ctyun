package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// EbsbackupEbsBackupPolicyUnbindDisksApi
/* 备份策略解绑云硬盘
 */type EbsbackupEbsBackupPolicyUnbindDisksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupEbsBackupPolicyUnbindDisksApi(client *core.CtyunClient) *EbsbackupEbsBackupPolicyUnbindDisksApi {
	return &EbsbackupEbsBackupPolicyUnbindDisksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ebs-backup/policy/unbind-disks",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupEbsBackupPolicyUnbindDisksApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupEbsBackupPolicyUnbindDisksRequest) (*EbsbackupEbsBackupPolicyUnbindDisksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(struct {
		*EbsbackupEbsBackupPolicyUnbindDisksRequest
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
	var resp EbsbackupEbsBackupPolicyUnbindDisksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupEbsBackupPolicyUnbindDisksRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID string `json:"policyID,omitempty"` /*  备份策略ID，您可以通过<a href="https://www.ctyun.cn/document/10026752/10040084">查询备份策略</a>获取  */
	DiskIDs  string `json:"diskIDs,omitempty"`  /*  云硬盘ID,如果解绑多个,请使用逗号隔开，您可以通过<a href="https://www.ctyun.cn/document/10026752/10076567">查询备份策略绑定的云硬盘列表</a>获取  */
}

type EbsbackupEbsBackupPolicyUnbindDisksResponse struct {
	StatusCode  int32                                                 `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                                `json:"message"`     /*  错误信息的英文描述  */
	Description string                                                `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupEbsBackupPolicyUnbindDisksReturnObjResponse `json:"returnObj"`   /*  无实际对象返回，值为{}  */
	ErrorCode   string                                                `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                                `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupEbsBackupPolicyUnbindDisksReturnObjResponse struct{}
