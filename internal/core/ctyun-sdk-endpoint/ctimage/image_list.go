package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"net/http"
	"strconv"
)

// ImageListApi 查询可以使用的镜像资源
// https://www.ctyun.cn/document/10027726/10040047
type ImageListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageListApi(client *ctyunsdk.CtyunClient) *ImageListApi {
	return &ImageListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/image/list",
		},
	}
}

func (this *ImageListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageListRequest) (*ImageListResponse, ctyunsdk.CtyunRequestError) {
	request := this.
		WithCredential(&credential).
		AddParam("regionID", req.RegionId).
		AddParam("azName", req.AzName).
		AddParam("flavorName", req.FlavorName).
		AddParam("pageNo", strconv.Itoa(req.PageNo)).
		AddParam("pageSize", strconv.Itoa(req.PageSize)).
		AddParam("queryContent", req.QueryContent).
		AddParam("status", req.Status).
		AddParam("visibility", strconv.Itoa(req.Visibility))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	realResponse := &imageListRealResponse{}
	err = response.ParseByStandardModelWithCheck(realResponse)
	if err != nil {
		return nil, err
	}
	var images []ImageListImagesResponse
	for _, img := range realResponse.Images {
		images = append(images, ImageListImagesResponse{
			Architecture:     img.Architecture,
			AzName:           img.AzName,
			BootMode:         img.BootMode,
			ContainerFormat:  img.ContainerFormat,
			CreatedTime:      img.CreatedTime,
			Description:      img.Description,
			DestinationUser:  img.DestinationUser,
			DiskFormat:       img.DiskFormat,
			DiskId:           img.DiskID,
			DiskSize:         img.DiskSize,
			ImageClass:       img.ImageClass,
			ImageId:          img.ImageID,
			ImageName:        img.ImageName,
			ImageType:        img.ImageType,
			MaximumRam:       img.MaximumRAM,
			MinimumRam:       img.MinimumRAM,
			OsDistro:         img.OsDistro,
			OsType:           img.OsType,
			OsVersion:        img.OsVersion,
			ProjectId:        img.ProjectID,
			SharedListLength: img.SharedListLength,
			Size:             img.Size,
			SourceServerId:   img.SourceServerID,
			SourceUser:       img.SourceUser,
			Status:           img.Status,
			Tags:             img.Tags,
			UpdatedTime:      img.UpdatedTime,
			Visibility:       img.Visibility,
		})
	}
	return &ImageListResponse{
		Images:       images,
		CurrentCount: realResponse.CurrentCount,
		TotalPage:    realResponse.TotalPage,
	}, nil
}

type imageListRealResponse struct {
	Images       []imageListImagesRealResponse
	CurrentCount int
	TotalCount   int
	TotalPage    int
}

type imageListImagesRealResponse struct {
	Architecture     string `json:"architecture"`
	AzName           string `json:"azName"`
	BootMode         string `json:"bootMode"`
	ContainerFormat  string `json:"containerFormat"`
	CreatedTime      int    `json:"createdTime"`
	Description      string `json:"description"`
	DestinationUser  string `json:"destinationUser"`
	DiskFormat       string `json:"diskFormat"`
	DiskID           string `json:"diskID"`
	DiskSize         int    `json:"diskSize"`
	ImageClass       string `json:"imageClass"`
	ImageID          string `json:"imageID"`
	ImageName        string `json:"imageName"`
	ImageType        string `json:"imageType"`
	MaximumRAM       int    `json:"maximumRAM"`
	MinimumRAM       int    `json:"minimumRAM"`
	OsDistro         string `json:"osDistro"`
	OsType           string `json:"osType"`
	OsVersion        string `json:"osVersion"`
	ProjectID        string `json:"projectID,omitempty"`
	SharedListLength int    `json:"sharedListLength"`
	Size             int64  `json:"size"`
	SourceServerID   string `json:"sourceServerID"`
	SourceUser       string `json:"sourceUser"`
	Status           string `json:"status"`
	Tags             string `json:"tags"`
	UpdatedTime      int    `json:"updatedTime"`
	Visibility       string `json:"visibility"`
}

type ImageListRequest struct {
	RegionId     string // 资源池id
	AzName       string // 可用区名称
	FlavorName   string // 云主机规格名称
	Visibility   int    // 镜像可见类型，取值范围（值：描述）：0：私有镜像1：公共镜像（默认值）2：共享镜像3：安全产品镜像4：甄选应用镜像
	PageNo       int    // 页码，取值范围：最小 1（默认值）
	PageSize     int    // 分页查询时每页的行数，最大值为 50，默认值为 10
	QueryContent string // 查询内容
	Status       string // 镜像状态 accepted：已接受共享镜像 rejected：已拒绝共享镜像 waiting：等待接受/拒绝共享镜像
	// ProjectId    string // 企业项目 ID
}

type ImageListImagesResponse struct {
	Architecture     string
	AzName           string
	BootMode         string
	ContainerFormat  string
	CreatedTime      int
	Description      string
	DestinationUser  string
	DiskFormat       string
	DiskId           string
	DiskSize         int
	ImageClass       string
	ImageId          string
	ImageName        string
	ImageType        string
	MaximumRam       int
	MinimumRam       int
	OsDistro         string
	OsType           string
	OsVersion        string
	ProjectId        string
	SharedListLength int
	Size             int64
	SourceServerId   string
	SourceUser       string
	Status           string
	Tags             string
	UpdatedTime      int
	Visibility       string
}

type ImageListResponse struct {
	Images       []ImageListImagesResponse
	CurrentCount int
	TotalCount   int
	TotalPage    int
}
