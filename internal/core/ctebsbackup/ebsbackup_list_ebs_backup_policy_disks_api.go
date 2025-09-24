package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// EbsbackupListEbsBackupPolicyDisksApi
/* 查询云硬盘备份策略绑定的云硬盘列表
 */type EbsbackupListEbsBackupPolicyDisksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewEbsbackupListEbsBackupPolicyDisksApi(client *core.CtyunClient) *EbsbackupListEbsBackupPolicyDisksApi {
	return &EbsbackupListEbsBackupPolicyDisksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ebs-backup/policy/list-disks",
			ContentType:  "application/json",
		},
	}
}

func (a *EbsbackupListEbsBackupPolicyDisksApi) Do(ctx context.Context, credential core.Credential, req *EbsbackupListEbsBackupPolicyDisksRequest) (*EbsbackupListEbsBackupPolicyDisksResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("policyID", req.PolicyID)
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.DiskID != "" {
		ctReq.AddParam("diskID", req.DiskID)
	}
	if req.DiskName != "" {
		ctReq.AddParam("diskName", req.DiskName)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp EbsbackupListEbsBackupPolicyDisksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type EbsbackupListEbsBackupPolicyDisksRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PolicyID string `json:"policyID,omitempty"` /*  备份策略ID  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码，默认1  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页显示条目，默认10  */
	DiskID   string `json:"diskID,omitempty"`   /*  云硬盘ID  */
	DiskName string `json:"diskName,omitempty"` /*  云硬盘名称，模糊过滤，指定diskID时，该参数无效  */
}

type EbsbackupListEbsBackupPolicyDisksResponse struct {
	StatusCode  int32                                               `json:"statusCode"`  /*  返回状态码（800为成功，900为失败）  */
	Message     string                                              `json:"message"`     /*  错误信息的英文描述  */
	Description string                                              `json:"description"` /*  错误信息的本地化描述（中文）  */
	ReturnObj   *EbsbackupListEbsBackupPolicyDisksReturnObjResponse `json:"returnObj"`   /*  返回对象  */
	ErrorCode   string                                              `json:"errorCode"`   /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
	Error       string                                              `json:"error"`       /*  业务错误细分码，发生错误时返回，为product.module.code三段式码  */
}

type EbsbackupListEbsBackupPolicyDisksReturnObjResponse struct {
	DiskList     []*EbsbackupListEbsBackupPolicyDisksReturnObjDiskListResponse `json:"diskList"`     /*  云硬盘列表  */
	TotalCount   int32                                                         `json:"totalCount"`   /*  云硬盘总数  */
	CurrentCount int32                                                         `json:"currentCount"` /*  当前页云硬盘数  */
}

type EbsbackupListEbsBackupPolicyDisksReturnObjDiskListResponse struct {
	DiskID       string `json:"diskID"`       /*  云硬盘ID  */
	DiskName     string `json:"diskName"`     /*  云硬盘名称  */
	DiskSize     int32  `json:"diskSize"`     /*  云硬盘大小，单位GB  */
	DiskStatus   string `json:"diskStatus"`   /*  云硬盘使用状态 available、in-use等，具体请参考[云硬盘使用状态](https://www.ctyun.cn/document/10027696/10168629)  */
	DiskType     string `json:"diskType"`     /*  云硬盘类型 SATA/SAS/SSD/FAST-SSD  */
	DiskMode     string `json:"diskMode"`     /*  云硬盘模式 VBD/ISCSI/FCSAN  */
	CreatedTime  int32  `json:"createdTime"`  /*  创建时间  */
	ExpiredTime  int32  `json:"expiredTime"`  /*  过期时间  */
	InstanceID   string `json:"instanceID"`   /*  云硬盘挂载的云主机ID  */
	InstanceName string `json:"instanceName"` /*  云硬盘挂载的云主机名称  */
}
