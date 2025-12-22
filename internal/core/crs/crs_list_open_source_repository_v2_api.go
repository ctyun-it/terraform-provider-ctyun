package crs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CrsListOpenSourceRepositoryV2Api
/* 分页查询开源镜像列表 */
type CrsListOpenSourceRepositoryV2Api struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCrsListOpenSourceRepositoryV2Api(client *core.CtyunClient) *CrsListOpenSourceRepositoryV2Api {
	return &CrsListOpenSourceRepositoryV2Api{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v2/listOpenSourceRepository",
			ContentType:  "application/json",
		},
	}
}

func (a *CrsListOpenSourceRepositoryV2Api) Do(ctx context.Context, credential core.Credential, req *CrsListOpenSourceRepositoryV2Request) (*CrsListOpenSourceRepositoryV2Response, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddHeader("Content-Type", req.ContentType)
	ctReq.AddHeader("regionId", req.RegionId)
	if req.RepositoryName != nil && *req.RepositoryName != "" {
		ctReq.AddParam("repositoryName", *req.RepositoryName)
	}
	if req.Category != nil && *req.Category != "" {
		ctReq.AddParam("category", *req.Category)
	}
	if req.Architecture != nil && *req.Architecture != "" {
		ctReq.AddParam("architecture", *req.Architecture)
	}
	if req.PageNum != nil && *req.PageNum != 0 {
		ctReq.AddParam("pageNum", strconv.FormatInt(int64(*req.PageNum), 10))
	}
	if req.PageNum != nil && *req.PageSize != 0 {
		ctReq.AddParam("pageSize", strconv.FormatInt(int64(*req.PageSize), 10))
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CrsListOpenSourceRepositoryV2Response
	err = response.Parse(&resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

type CrsListOpenSourceRepositoryV2Request struct {
	ContentType    string  `json:"Content-Type"`             /*  类型  */
	RegionId       string  `json:"regionId"`                 /*  资源池编码（资源池ID，您可以查看<a href="https://www.ctyun.cn/document/10026765/11005378" target="_blank">容器镜像服务资源池</a>获取）  */
	RepositoryName *string `json:"repositoryName,omitempty"` /*  镜像仓库名称，支持模糊搜索  */
	Category       *string `json:"category,omitempty"`       /*  根据分类标签筛选（base：基础镜像，os：操作系统，ai：AI，middleware：中间件，storage：存储，network：网络，other：其他），支持根据多个标签筛选（多个标签使用,分隔，筛选结果为满足任一标签的镜像），如果为空表示包含所有分类  */
	Architecture   *string `json:"architecture,omitempty"`   /*  根据架构筛选（arm64，amd64），支持根据多个标签筛选（多个标签使用,分隔，筛选结果为满足任一标签的镜像），如果为空表示包含所有架构  */
	PageNum        *int32  `json:"pageNum,omitempty"`        /*  当前页码（默认为1）  */
	PageSize       *int32  `json:"pageSize,omitempty"`       /*  每页条数（默认为10,最大值为50）  */
}

type CrsListOpenSourceRepositoryV2Response struct {
	StatusCode int32                                           `json:"statusCode"` /*  响应码 （800为请求成功，900为失败 ）  */
	Message    *string                                         `json:"message"`    /*  返回信息  */
	Error      *string                                         `json:"error"`      /*  错误码  */
	ReturnObj  *CrsListOpenSourceRepositoryV2ReturnObjResponse `json:"returnObj"`  /*  返回结果  */
}

type CrsListOpenSourceRepositoryV2ReturnObjResponse struct {
	Total   *int32                                                   `json:"total"`   /*  总条数  */
	Size    *int32                                                   `json:"size"`    /*  每页条数  */
	Current *int32                                                   `json:"current"` /*  当前页码  */
	Pages   *int32                                                   `json:"pages"`   /*  总页数  */
	Records []*CrsListOpenSourceRepositoryV2ReturnObjRecordsResponse `json:"records"` /*  开源镜像仓库列表  */
}

type CrsListOpenSourceRepositoryV2ReturnObjRecordsResponse struct {
	NamespaceName    *string   `json:"namespaceName"`    /*  命名空间名称  */
	RepositoryName   *string   `json:"repositoryName"`   /*  镜像仓库名称  */
	RepositoryId     *string   `json:"repositoryId"`     /*  镜像仓库id  */
	ImageUrl         *string   `json:"imageUrl"`         /*  公网拉取地址  */
	ImageUrlInternal *string   `json:"imageUrlInternal"` /*  内网拉取地址  */
	RegionId         *string   `json:"regionId"`         /*  资源池编码  */
	Starred          *bool     `json:"starred"`          /*  是否收藏  */
	Category         []*string `json:"category"`         /*  分类标签列表  */
	Architecture     []*string `json:"architecture"`     /*  支持架构列表  */
	Os               []*string `json:"os"`               /*  支持系统列表  */
}
