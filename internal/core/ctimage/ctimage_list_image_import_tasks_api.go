package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
	"strconv"
)

// CtimageListImageImportTasksApi
/* 查询使用镜像文件来创建私有镜像的任务列表<br />准备：<br />1. 在调用前需了解如何构造请求，可参见：如何调用 API - 构造请求<br />2. OpenAPI 请求需进行加密调用，可参见：如何调用 API - 认证鉴权<br />注意：<br />1. 此接口已弃用，目前仍可调用。若您欲查询镜像状态为 queued 等的私有镜像，则可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4763&data=89&isNormal=1&vid=83" target="_blank">查询可以使用的镜像资源</a>接口来查询您可使用的镜像资源，可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=4764&data=89&isNormal=1&vid=83" target="_blank">查询镜像详细信息</a>接口来查询 1 份镜像的详细信息<br />2. 在调用前，请您认真阅读此文档，包括但不限于参数描述中的“注意”部分
 */type CtimageListImageImportTasksApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewCtimageListImageImportTasksApi(client *core.CtyunClient) *CtimageListImageImportTasksApi {
	return &CtimageListImageImportTasksApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/image/list-import-tasks",
			ContentType:  "application/json",
		},
	}
}

func (a *CtimageListImageImportTasksApi) Do(ctx context.Context, credential core.Credential, req *CtimageListImageImportTasksRequest) (*CtimageListImageImportTasksResponse, error) {
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
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp CtimageListImageImportTasksResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type CtimageListImageImportTasksRequest struct {
	RegionID string `json:"regionID,omitempty"` /*  资源池 ID。可使用<a href="https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=25&api=5851&data=87&isNormal=1&vid=81" target="_blank">资源池列表查询</a>接口来查询您可见的资源池的列表  */
	PageNo   int32  `json:"pageNo,omitempty"`   /*  页码。取值范围：最小 1（默认值）  */
	PageSize int32  `json:"pageSize,omitempty"` /*  每页记录数目。取值范围：最小 1，最大 50，默认 10  */
}

type CtimageListImageImportTasksResponse struct {
	StatusCode  int32                                         `json:"statusCode,omitempty"`  /*  状态码。取值范围（值：描述）：<br />800：成功，<br />900：失败  */
	Error       string                                        `json:"error,omitempty"`       /*  错误码（product.module.code 三段式码）  */
	ErrorCode   string                                        `json:"errorCode,omitempty"`   /*  同 error 参数  */
	Message     string                                        `json:"message,omitempty"`     /*  英文描述信息  */
	Description string                                        `json:"description,omitempty"` /*  中文描述信息  */
	ReturnObj   *CtimageListImageImportTasksReturnObjResponse `json:"returnObj"`             /*  成功时返回的数据  */
}

type CtimageListImageImportTasksReturnObjResponse struct {
	ImageImportTasks []string `json:"imageImportTasks"`       /*  创建私有镜像（镜像文件）任务列表  */
	CurrentPage      int32    `json:"currentPage,omitempty"`  /*  当前页码  */
	CurrentCount     int32    `json:"currentCount,omitempty"` /*  当前页记录数  */
	TotalPage        int32    `json:"totalPage,omitempty"`    /*  总页数  */
	TotalCount       int32    `json:"totalCount,omitempty"`   /*  总记录数  */
}
