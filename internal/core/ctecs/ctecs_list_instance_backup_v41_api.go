package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsListInstanceBackupV41Api
/* 查询云主机备份列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsListInstanceBackupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceBackupV41Api(client *core.CtyunClient) *CtecsListInstanceBackupV41Api {
	return &CtecsListInstanceBackupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceBackupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceBackupV41Request) (*CtecsListInstanceBackupV41Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListInstanceBackupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceBackupV41Request struct {
	RegionID             string `json:"regionID,omitempty"`             /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	PageNo               int32  `json:"pageNo,omitempty"`               /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize             int32  `json:"pageSize,omitempty"`             /*  每页记录数目，取值范围：[1, 50]，注：默认值为10  */
	InstanceID           string `json:"instanceID,omitempty"`           /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8281&data=87">创建一台按量付费或包年包月的云主机</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8282&data=87">批量创建按量付费或包年包月云主机</a>  */
	RepositoryID         string `json:"repositoryID,omitempty"`         /*  云主机备份存储库ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033742">产品定义-存储库</a>来了解存储库<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=6909&data=100">查询存储库列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=6910&data=100">创建存储库</a>  */
	InstanceBackupID     string `json:"instanceBackupID,omitempty"`     /*  云主机备份ID，您可以查看<a href="https://www.ctyun.cn/document/10051003/10051023">产品定义-云主机备份</a>来了解云主机备份<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=6910&data=100">查询云主机备份列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=14&api=8332&data=100">创建云主机备份</a>  */
	QueryContent         string `json:"queryContent,omitempty"`         /*  模糊匹配查询内容（匹配字段：instanceBackupName、instanceBackupID、instanceBackupStatus、instanceName）  */
	InstanceBackupStatus string `json:"instanceBackupStatus,omitempty"` /*  云主机备份状态  */
	ProjectID            string `json:"projectID,omitempty"`            /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
}

type CtecsListInstanceBackupV41Response struct {
	StatusCode  int32                                        `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                       `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                       `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                       `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                       `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsListInstanceBackupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListInstanceBackupV41ReturnObjResponse struct {
	CurrentCount int32                                                 `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                 `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                 `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListInstanceBackupV41ReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListInstanceBackupV41ReturnObjResultsResponse struct {
	InstanceBackupID          string `json:"instanceBackupID,omitempty"`          /*  云主机备份ID  */
	InstanceBackupName        string `json:"instanceBackupName,omitempty"`        /*  云主机备份名称  */
	InstanceBackupStatus      string `json:"instanceBackupStatus,omitempty"`      /*  云主机备份状态，取值范围：<br />CREATING: 备份创建中, <br />ACTIVE: 可用， <br />RESTORING: 备份恢复中，<br />DELETING: 删除中，<br />EXPIRED：到期，<br />ERROR：错误  */
	InstanceBackupDescription string `json:"instanceBackupDescription,omitempty"` /*  云主机备份描述  */
	InstanceID                string `json:"instanceID,omitempty"`                /*  云主机ID  */
	InstanceName              string `json:"instanceName,omitempty"`              /*  云主机名称  */
	RepositoryID              string `json:"repositoryID,omitempty"`              /*  云主机备份存储库ID  */
	RepositoryName            string `json:"repositoryName,omitempty"`            /*  云主机备份存储库名称  */
	RepositoryExpired         *bool  `json:"repositoryExpired"`                   /*  云主机备份存储库是否过期  */
	RepositoryFreeze          *bool  `json:"repositoryFreeze"`                    /*  存储库是否冻结  */
	DiskTotalSize             int32  `json:"diskTotalSize,omitempty"`             /*  云盘总容量大小，单位为GB  */
	UsedSize                  int64  `json:"usedSize,omitempty"`                  /*  磁盘备份已使用大小  */
	DiskCount                 int32  `json:"diskCount,omitempty"`                 /*  云盘数目  */
	RestoreFinishedTime       string `json:"restoreFinishedTime,omitempty"`       /*  备份恢复完成时间  */
	CreatedTime               string `json:"createdTime,omitempty"`               /*  创建时间  */
	FinishedTime              string `json:"finishedTime,omitempty"`              /*  完成时间  */
	ProjectID                 string `json:"projectID,omitempty"`                 /*  企业项目ID  */
	BackupType                string `json:"backupType,omitempty"`                /*  备份类型，取值范围：●FULL：全量备份。●INCREMENT：增量备份。只有4.0资源池此参数有效  */
}
