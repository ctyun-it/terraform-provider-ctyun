package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// CtecsCreateInstanceBackupV41Api
/* 创建云主机备份<br/><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsCreateInstanceBackupV41Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsCreateInstanceBackupV41Api(client *core.CtyunClient) *CtecsCreateInstanceBackupV41Api {
	return &CtecsCreateInstanceBackupV41Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/ecs/backup/create",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsCreateInstanceBackupV41Api) Do(ctx context.Context, credential core.Credential, req *CtecsCreateInstanceBackupV41Request) (*CtecsCreateInstanceBackupV41Response, error) {
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
	var resp CtecsCreateInstanceBackupV41Response
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsCreateInstanceBackupV41Request struct {
	RegionID                  string `json:"regionID,omitempty"`                  /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	InstanceID                string `json:"instanceID,omitempty"`                /*  云主机ID，您可以查看<a href="https://www.ctyun.cn/products/ecs">弹性云主机</a>了解云主机的相关信息<br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=8309&data=87">查询云主机列表</a>  */
	InstanceBackupName        string `json:"instanceBackupName,omitempty"`        /*  云主机备份名称。满足以下规则：长度为2-63字符，头尾不支持输入空格  */
	InstanceBackupDescription string `json:"instanceBackupDescription,omitempty"` /*  云主机备份描述，字符长度不超过256字符  */
	RepositoryID              string `json:"repositoryID,omitempty"`              /*  云主机备份存储库ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033742">产品定义-存储库</a>来了解存储库<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6909&data=87&isNormal=1&vid=81">查询存储库列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6910&data=87&isNormal=1&vid=81">创建存储库</a>  */
	FullBackup                bool   `json:"fullBackup,omitempty"`                /*  是否启用全量备份，取值范围：●true：是●false：否若启用该参数，则此次备份的类型为全量备份。注：只有4.0资源池支持该参数。  */

}

type CtecsCreateInstanceBackupV41Response struct {
	StatusCode  int32                                          `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                         `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                         `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                         `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                         `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsCreateInstanceBackupV41ReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsCreateInstanceBackupV41ReturnObjResponse struct {
	Results *CtecsCreateInstanceBackupV41ReturnObjResultsResponse `json:"results"` /*  备份结果  */
}

type CtecsCreateInstanceBackupV41ReturnObjResultsResponse struct {
	InstanceBackupID          string `json:"instanceBackupID,omitempty"`          /*  云主机备份ID  */
	InstanceBackupName        string `json:"instanceBackupName,omitempty"`        /*  云主机备份名称  */
	InstanceBackupStatus      string `json:"instanceBackupStatus,omitempty"`      /*  备份状态，取值范围：<br />CREATING: 备份创建中, <br />ACTIVE: 可用， <br />RESTORING: 备份恢复中，<br />DELETING: 删除中，<br />EXPIRED：到期，<br />ERROR：错误  */
	InstanceBackupDescription string `json:"instanceBackupDescription,omitempty"` /*  云主机备份描述  */
	InstanceID                string `json:"instanceID,omitempty"`                /*  云主机ID  */
	InstanceName              string `json:"instanceName,omitempty"`              /*  云主机名称  */
	RepositoryID              string `json:"repositoryID,omitempty"`              /*  云主机备份存储库ID  */
	RepositoryName            string `json:"repositoryName,omitempty"`            /*  云主机备份存储库名称  */
	DiskTotalSize             int32  `json:"diskTotalSize,omitempty"`             /*  云盘总容量大小，单位为GB  */
	UsedSize                  int32  `json:"usedSize,omitempty"`                  /*  云硬盘备份已使用大小  */
	CreatedTime               string `json:"createdTime,omitempty"`               /*  创建时间  */
	ProjectID                 string `json:"projectID,omitempty"`                 /*  企业项目ID  */
}
