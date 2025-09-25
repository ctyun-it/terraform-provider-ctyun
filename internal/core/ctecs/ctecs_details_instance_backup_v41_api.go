package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsDetailsInstanceBackupV41Api
/* 查询云主机备份详情<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsDetailsInstanceBackupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsDetailsInstanceBackupV41Api(client *core.CtyunClient) *CtecsDetailsInstanceBackupV41Api {
	return &CtecsDetailsInstanceBackupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup/details",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsDetailsInstanceBackupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsDetailsInstanceBackupV41Request) (*CtecsDetailsInstanceBackupV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("instanceBackupID", req.InstanceBackupID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsDetailsInstanceBackupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsDetailsInstanceBackupV41Request struct {
	RegionID         string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceBackupID string /*  云主机备份ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033738">产品定义-云主机备份</a>来了解云主机备份<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8354&data=87&isNormal=1&vid=81">查询云主机备份列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8332&data=87&isNormal=1&vid=81">创建云主机备份</a>   */
}

type CtecsDetailsInstanceBackupV41Response struct {
	StatusCode  int32                                           `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                          `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                          `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                          `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                          `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsDetailsInstanceBackupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsDetailsInstanceBackupV41ReturnObjResponse struct {
	InstanceBackupID          string `json:"instanceBackupID,omitempty"`          /*  云主机备份ID  */
	InstanceBackupName        string `json:"instanceBackupName,omitempty"`        /*  云主机备份名称  */
	InstanceBackupStatus      string `json:"instanceBackupStatus,omitempty"`      /*  备份状态，取值范围：<br />CREATING: 备份创建中, <br />ACTIVE: 可用， <br />RESTORING: 备份恢复中，<br />DELETING: 删除中，<br />EXPIRED：到期，<br />ERROR：错误  */
	InstanceBackupDescription string `json:"instanceBackupDescription,omitempty"` /*  云主机备份描述  */
	InstanceID                string `json:"instanceID,omitempty"`                /*  云主机ID  */
	InstanceName              string `json:"instanceName,omitempty"`              /*  云主机名称  */
	RepositoryID              string `json:"repositoryID,omitempty"`              /*  云主机备份存储库ID  */
	RepositoryName            string `json:"repositoryName,omitempty"`            /*  云主机备份存储库名称  */
	RepositoryExpired         *bool  `json:"repositoryExpired"`                   /*  云主机备份存储库是否过期  */
	RepositoryFreeze          *bool  `json:"repositoryFreeze"`                    /*  存储库是否冻结  */
	DiskTotalSize             int32  `json:"diskTotalSize,omitempty"`             /*  云硬盘总容量大小  */
	UsedSize                  int64  `json:"usedSize,omitempty"`                  /*  磁盘备份已使用大小  */
	DiskCount                 int32  `json:"diskCount,omitempty"`                 /*  云硬盘数目  */
	RestoreFinishedTime       string `json:"restoreFinishedTime,omitempty"`       /*  备份恢复完成时间  */
	CreatedTime               string `json:"createdTime,omitempty"`               /*  创建时间  */
	FinishedTime              string `json:"finishedTime,omitempty"`              /*  完成时间  */
	ProjectID                 string `json:"projectID,omitempty"`                 /*  企业项目ID  */
	BackupType                string `json:"backupType,omitempty"`                /*  备份类型，取值范围：●FULL：全量备份。●INCREMENT：增量备份。只有4.0资源池此参数有效  */
}
