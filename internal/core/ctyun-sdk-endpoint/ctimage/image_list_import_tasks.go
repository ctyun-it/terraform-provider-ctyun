package ctimage

import (
	"context"
	"net/http"
	"strconv"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// ImageListImportTasksApi 查询创建私有镜像（镜像文件）任务列表
// https://www.ctyun.cn/document/10027726/10087256
type ImageListImportTasksApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageListImportTasksApi(client *ctyunsdk.CtyunClient) *ImageListImportTasksApi {
	return &ImageListImportTasksApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/image/list-import-tasks",
		},
	}
}

func (this *ImageListImportTasksApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageListImportTasksRequest) (*ImageListImportTasksResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	request.AddParam("regionID", req.RegionId)
	request.AddParam("pageNo", strconv.Itoa(req.PageNo))
	request.AddParam("pageSize", strconv.Itoa(req.PageSize))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	var resp imageListImportTasksRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	return &ImageListImportTasksResponse{
		PageNo:       resp.PageNo,
		CurrentPage:  resp.CurrentPage,
		PageSize:     resp.PageSize,
		CurrentCount: resp.CurrentCount,
		TotalCount:   resp.TotalCount,
		ImageImportTasks: ImageListImportTasksImageImportTasksResponse{
			ImageName:  resp.ImageImportTasks.ImageName,
			OsType:     resp.ImageImportTasks.OsType,
			TaskId:     resp.ImageImportTasks.TaskID,
			TaskStatus: resp.ImageImportTasks.TaskStatus,
		},
	}, nil
}

type imageListImportTasksImageImportTasksRealResponse struct {
	ImageName  string `json:"imageName"`
	OsType     string `json:"osType"`
	TaskID     string `json:"taskID"`
	TaskStatus string `json:"taskStatus"`
}

type imageListImportTasksRealResponse struct {
	PageNo           int                                              `json:"pageNo"`
	CurrentPage      int                                              `json:"currentPage"`
	PageSize         int                                              `json:"pageSize"`
	CurrentCount     int                                              `json:"currentCount"`
	TotalCount       int                                              `json:"totalCount"`
	ImageImportTasks imageListImportTasksImageImportTasksRealResponse `json:"imageImportTasks"`
}

type ImageListImportTasksRequest struct {
	RegionId string // 资源池id
	PageNo   int    // 页码，取值范围：最小 1（默认值）
	PageSize int    // 每页记录数目，取值范围：最小 1，最大 50，默认值 10
}

type ImageListImportTasksImageImportTasksResponse struct {
	ImageName  string
	OsType     string
	TaskId     string
	TaskStatus string
}

type ImageListImportTasksResponse struct {
	PageNo           int
	CurrentPage      int
	PageSize         int
	CurrentCount     int
	TotalCount       int
	ImageImportTasks ImageListImportTasksImageImportTasksResponse
}
