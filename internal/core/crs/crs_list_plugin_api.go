package crs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CrsListPluginApi
/* 查询插件市场列表
 */type CrsListPluginApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCrsListPluginApi(client *core.CtyunClient) *CrsListPluginApi {
	return &CrsListPluginApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v1/listPlugin",
			ContentType:  "application/json",
		},
	}
}

func (a *CrsListPluginApi) Do(ctx context.Context, credential core.Credential, req *CrsListPluginRequest) (*CrsListPluginResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("Content-Type", req.ContentType)
	ctReq.AddHeader("regionId", req.RegionIdHeader)
	ctReq.AddParam("regionId", req.RegionIdParam)
	if req.NamespaceName != "" {
		ctReq.AddParam("namespaceName", req.NamespaceName)
	}
	if req.RepositoryName != "" {
		ctReq.AddParam("repositoryName", req.RepositoryName)
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
	var resp CrsListPluginResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CrsListPluginRequest struct {
	ContentType    string `json:"Content-Type,omitempty"`   /*  类型  */
	RegionIdHeader string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	RegionIdParam  string `json:"regionId,omitempty"`       /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	NamespaceName  string `json:"namespaceName,omitempty"`  /*  命名空间名称，默认为plugins(plugins:专有版和托管版集群插件市场, sce:Serverless集群插件市场)  */
	RepositoryName string `json:"repositoryName,omitempty"` /*  插件名称，支持模糊查询  */
	PageNum        int32  `json:"pageNum,omitempty"`        /*  当前页码（默认为1）  */
	PageSize       int32  `json:"pageSize,omitempty"`       /*  每页条数（默认为10）  */
	OrderBy        string `json:"orderBy,omitempty"`        /*  排序字段(repositoryName：插件名称，createdTime：创建时间，modifiedTime：修改时间，默认值：modifiedTime)  */
	Order          string `json:"order,omitempty"`          /*  排序方式(desc：降序排序, asc：升序排序，默认值：desc)  */
}

type CrsListPluginResponse struct {
	StatusCode int32                           `json:"statusCode,omitempty"` /*  响应码（0为请求成功，其它为失败 ）  */
	Message    string                          `json:"message,omitempty"`    /*  返回信息  */
	Error      string                          `json:"error,omitempty"`      /*  错误码  */
	ReturnObj  *CrsListPluginReturnObjResponse `json:"returnObj"`            /*  返回结果  */
}

type CrsListPluginReturnObjResponse struct {
	Total   int32                                    `json:"total,omitempty"`   /*  总条数  */
	Size    int32                                    `json:"size,omitempty"`    /*  每页条数  */
	Current int32                                    `json:"current,omitempty"` /*  当前页码  */
	Pages   int32                                    `json:"pages,omitempty"`   /*  总页数  */
	Records []*CrsListPluginReturnObjRecordsResponse `json:"records"`           /*  插件市场列表  */
}

type CrsListPluginReturnObjRecordsResponse struct {
	NamespaceName    string `json:"namespaceName,omitempty"`    /*  命名空间名称  */
	RepositoryName   string `json:"repositoryName,omitempty"`   /*  插件名称  */
	ImageUrl         string `json:"imageUrl,omitempty"`         /*  公网地址  */
	ImageUrlInternal string `json:"imageUrlInternal,omitempty"` /*  内网地址  */
}
