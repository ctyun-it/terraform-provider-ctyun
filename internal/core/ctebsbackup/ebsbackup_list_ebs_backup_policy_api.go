package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListEbsBackupPolicyApi
/* 查询云硬盘备份策略列表
 */type EbsbackupListEbsBackupPolicyApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListEbsBackupPolicyApi(client *core.CtyunClient) *EbsbackupListEbsBackupPolicyApi {
	return &EbsbackupListEbsBackupPolicyApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/policy/list",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListEbsBackupPolicyApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListEbsBackupPolicyRequest) (*EbsbackupListEbsBackupPolicyResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.PolicyID != "" {
		ctReq.AddParam("policyID", req.PolicyID)
	}
	if req.PolicyName != "" {
		ctReq.AddParam("policyName", req.PolicyName)
	}
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListEbsBackupPolicyResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListEbsBackupPolicyRequest struct {
	RegionID   string `json:"regionID,omitempty"`   /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PageNo     int32  `json:"pageNo,omitempty"`     /*  页码，默认值1  */
	PageSize   int32  `json:"pageSize,omitempty"`   /*  每页记录数目 ,默认10  */
	PolicyID   string `json:"policyID,omitempty"`   /*  备份策略ID  */
	PolicyName string `json:"policyName,omitempty"` /*  备份策略名，指定了policyID时，该参数会被忽略  */
	ProjectID  string `json:"projectID,omitempty"`  /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10026730/10238876">创建企业项目</a>了解如何创建企业项目<br />注：默认值为"0"  */
}

type EbsbackupListEbsBackupPolicyResponse struct {
	StatusCode  int32                                          `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                         `json:"message"`     /*  错误信息的英文描述  */
	Description string                                         `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupListEbsBackupPolicyReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                         `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                         `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupListEbsBackupPolicyReturnObjResponse struct{}
