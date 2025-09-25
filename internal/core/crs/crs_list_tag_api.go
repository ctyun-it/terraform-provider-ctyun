package crs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CrsListTagApi
/* 查询版本列表
 */type CrsListTagApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCrsListTagApi(client *core.CtyunClient) *CrsListTagApi {
	return &CrsListTagApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/listTag",
			ContentType:  "application/json",
		},
	}
}

func (a *CrsListTagApi) Do(ctx context.Context, credential core.Credential, req *CrsListTagRequest) (*CrsListTagResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("Content-Type", req.ContentType)
	ctReq.AddHeader("regionId", req.RegionIdHeader)
	ctReq.AddParam("regionId", req.RegionIdParam)
	ctReq.AddParam("namespaceName", req.NamespaceName)
	ctReq.AddParam("repositoryName", req.RepositoryName)
	if req.TagName != "" {
		ctReq.AddParam("tagName", req.TagName)
	}
	if req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(req.PageNum), 10))
	}
	if req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(req.PageSize), 10))
	}
	if req.OrderBy != "" {
		ctReq.AddParam("orderBy", req.OrderBy)
	}
	if req.Order != "" {
		ctReq.AddParam("order", req.Order)
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CrsListTagResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CrsListTagRequest struct {
	ContentType    string `json:"Content-Type,omitempty"`   /*  类型  */
	RegionIdHeader string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	RegionIdParam  string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	NamespaceName  string `json:"namespaceName,omitempty"`  /*  命名空间名称  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  插件名称  */
	TagName        string `json:"tagName,omitempty"`        /*  版本名称  */
	PageNum        int32  `json:"pageNum,omitempty"`        /*  当前页码（默认为1）  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  每页条数（默认为10）  */
	OrderBy        string `json:"orderBy,omitempty"`        /*  排序字段(name：名称，pushTime：推送时间，默认值：name)  */
	Order          string `json:"order,omitempty"`          /*  排序方式(desc：降序排序, asc：升序排序，默认值：desc)  */
}

type CrsListTagResponse struct {
	StatusCode int32                        `json:"statusCode,omitempty"` /*  响应码（0为请求成功，其它为失败 ）  */
	Message    string                       `json:"message,omitempty"`    /*  返回信息  */
	Error      string                       `json:"error,omitempty"`      /*  错误码  */
	ReturnObj  *CrsListTagReturnObjResponse `json:"returnObj"`            /*  返回结果  */
}

type CrsListTagReturnObjResponse struct {
	Total   int32                                 `json:"total,omitempty"`   /*  总条数  */
	Size    int32                                 `json:"size,omitempty"`    /*  每页条数  */
	Current int32                                 `json:"current,omitempty"` /*  当前页码  */
	Pages   int32                                 `json:"pages,omitempty"`   /*  总页数  */
	Records []*CrsListTagReturnObjRecordsResponse `json:"records"`           /*  版本列表  */
}

type CrsListTagReturnObjRecordsResponse struct {
	Name        string                                         `json:"name,omitempty"`        /*  名称  */
	Size        string                                         `json:"size,omitempty"`        /*  大小  */
	Description string                                         `json:"description,omitempty"` /*  描述  */
	Digest      string                                         `json:"digest,omitempty"`      /*  Hash 值  */
	Push_time   string                                         `json:"push_time,omitempty"`   /*  推送时间  */
	Extra_attrs *CrsListTagReturnObjRecordsExtra_attrsResponse `json:"extra_attrs"`           /*  额外属性  */
}

type CrsListTagReturnObjRecordsExtra_attrsResponse struct {
	ApiVersion  string `json:"apiVersion,omitempty"`  /*  API 版本  */
	AppVersion  string `json:"appVersion,omitempty"`  /*  APP 版本  */
	KubeVersion string `json:"kubeVersion,omitempty"` /*  兼容的 Kubernetes 版本  */
}
