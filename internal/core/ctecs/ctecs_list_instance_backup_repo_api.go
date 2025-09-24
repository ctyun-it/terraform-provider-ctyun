package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtecsListInstanceBackupRepoApi
/* 根据过滤参数查询满足条件的云主机备份存储库，并返回云主机备份存储库列表<br /><b>准备工作</b>：<br />&emsp;&emsp;构造请求：在调用前需要了解如何构造请求，详情查看构造请求<br />&emsp;&emsp;认证鉴权：openapi请求需要进行加密调用，详细查看认证鉴权<br />
 */type CtecsListInstanceBackupRepoApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtecsListInstanceBackupRepoApi(client *core.CtyunClient) *CtecsListInstanceBackupRepoApi {
	return &CtecsListInstanceBackupRepoApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/ecs/backup-repo/list",
			ContentType:  "application/json",
		},
	}
}

func (a *CtecsListInstanceBackupRepoApi) Do(ctx context.Context, credential core.Credential, req *CtecsListInstanceBackupRepoRequest) (*CtecsListInstanceBackupRepoResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	if req.ProjectID != "" {
		ctReq.AddParam("projectID", req.ProjectID)
	}
	if req.RepositoryName != "" {
		ctReq.AddParam("repositoryName", req.RepositoryName)
	}
	if req.RepositoryID != "" {
		ctReq.AddParam("repositoryID", req.RepositoryID)
	}
	if req.Status != "" {
		ctReq.AddParam("status", req.Status)
	}
	if req.PageNo != 0 {
		ctReq.AddParam("pageNo", strconv.FormatInt(int64(req.PageNo), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtecsListInstanceBackupRepoResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtecsListInstanceBackupRepoRequest struct {
	RegionID       string /*  资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026730/10028695">地域和可用区</a>来了解资源池 <br />获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87">资源池列表查询</a>  */
	ProjectID      string /*  企业项目ID，企业项目管理服务提供统一的云资源按企业项目管理，以及企业项目内的资源管理，成员管理。您可以通过查看<a href="https://www.ctyun.cn/document/10017248/10017961">创建企业项目</a>了解如何创建企业项目  */
	RepositoryName string /*  云主机备份存储库名称，满足以下规则：只能由数字、字母、-组成，不能以数字和-开头、且不能以-结尾，长度为2-63字符  */
	RepositoryID   string /*  云主机备份存储库ID，您可以查看<a href="https://www.ctyun.cn/document/10026751/10033742">产品定义-存储库</a>来了解存储库<br /> 获取：<br /><span style="background-color: rgb(73, 204, 144);color: rgb(255,255,255);padding: 2px; margin:2px">查</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6909&data=87&isNormal=1&vid=81">查询存储库列表</a><br /><span style="background-color: rgb(97, 175, 254);color: rgb(255,255,255);padding: 2px; margin:2px">创</span> <a  href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=6910&data=87&isNormal=1&vid=81">创建存储库</a>  */
	Status         string /*  存储库状态  */
	PageNo         int32  /*  页码，取值范围：正整数（≥1），注：默认值为1  */
	PageSize       int32  /*  每页记录数目，取值范围：[1~50]，默认值：10，单页最大记录不超过50  */
}

type CtecsListInstanceBackupRepoResponse struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  返回状态码（800为成功，900为失败）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  错误码，为product.module.code三段式码  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码，为product.module.code三段式码  */
	Message     string                                        `json:"message,omitempty"`     /*  失败时的错误描述，一般为英文描述  */
	Description string                                        `json:"description,omitempty"` /*  失败时的错误描述，一般为中文描述  */
	ReturnObj   *CtecsListInstanceBackupRepoReturnObjResponse `json:"returnObj"`             /*  返回参数  */
}

type CtecsListInstanceBackupRepoReturnObjResponse struct {
	CurrentCount int32                                                  `json:"currentCount,omitempty"` /*  当前页记录数目  */
	TotalCount   int32                                                  `json:"totalCount,omitempty"`   /*  总记录数  */
	TotalPage    int32                                                  `json:"totalPage,omitempty"`    /*  总页数  */
	Results      []*CtecsListInstanceBackupRepoReturnObjResultsResponse `json:"results"`                /*  分页明细  */
}

type CtecsListInstanceBackupRepoReturnObjResultsResponse struct {
	RegionID       string   `json:"regionID,omitempty"`       /*  资源池ID  */
	RepositoryID   string   `json:"repositoryID,omitempty"`   /*  存储库ID  */
	ProjectID      string   `json:"projectID,omitempty"`      /*  企业项目ID  */
	RepositoryName string   `json:"repositoryName,omitempty"` /*  存储库名称  */
	Status         string   `json:"status,omitempty"`         /*  云主机存储库状态，<br />expired: 已到期，<br />active: 可用  */
	Size           int32    `json:"size,omitempty"`           /*  云主机存储库总容量，单位GB  */
	FreeSize       float32  `json:"freeSize"`                 /*  云主机存储库剩余大小，单位GB(废弃该字段)   */
	RemainingSize  float32  `json:"remainingSize"`            /*  云主机存储库剩余大小，单位GB  */
	UsedSize       int64    `json:"usedSize,omitempty"`       /*  云主机存储库使用大小，单位Byte  */
	CreatedAt      string   `json:"createdAt,omitempty"`      /*  创建时间  */
	ExpiredAt      string   `json:"expiredAt,omitempty"`      /*  到期时间  */
	Expired        bool     `json:"expired"`                  /*  存储库是否到期  */
	Freeze         *bool    `json:"freeze"`                   /*  是否冻结  */
	Paas           *bool    `json:"paas"`                     /*  是否支持PAAS  */
	BackupList     []string `json:"backupList"`               /*  存储库下可用的备份ID列表  */
	BackupCount    int32    `json:"backupCount,omitempty"`    /*  存储库中备份数量  */
}
