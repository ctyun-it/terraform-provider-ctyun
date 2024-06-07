package ctimage

import (
	"context"
	"net/http"
	"strconv"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// ImageShareListApi 查询私有镜像的共享列表
// https://eop.ctyun.cn/ebp/ctapiDocument/search?sid=23&api=6764&data=89
type ImageShareListApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewimageShareListApi(client *ctyunsdk.CtyunClient) *ImageShareListApi {
	return &ImageShareListApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodGet,
			UrlPath: "/v4/image/show-shared-list",
		},
	}
}

func (this *ImageShareListApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageShareListRequest) (*ImageShareListResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	request.AddParam("imageID", req.ImageId).
		AddParam("regionID", req.RegionId).
		AddParam("pageNo", strconv.Itoa(req.PageNo)).
		AddParam("pageSize", strconv.Itoa(req.PageSize))

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	var resp imageShareListRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var images []ImageShareListImagesResponse
	for _, image := range resp.Images {
		images = append(images, ImageShareListImagesResponse{
			Architecture:    image.Architecture,
			BootMode:        image.BootMode,
			ContainerFormat: image.ContainerFormat,
			CreatedTime:     image.CreatedTime,
			Description:     image.Description,
			DestinationUser: image.DestinationUser,
			DiskFormat:      image.DiskFormat,
			DiskId:          image.DiskID,
			DiskSize:        image.DiskSize,
			ImageClass:      image.ImageClass,
			ImageId:         image.ImageID,
			ImageName:       image.ImageName,
			ImageType:       image.ImageType,
			MaximumRam:      image.MaximumRAM,
			MinimumRam:      image.MinimumRAM,
			OsDistro:        image.OsDistro,
			OsType:          image.OsType,
			OsVersion:       image.OsVersion,
			ProjectId:       image.ProjectID,
			Size:            image.Size,
			SourceServerId:  image.SourceServerID,
			SourceUser:      image.SourceUser,
			Status:          image.Status,
			UpdatedTime:     image.UpdatedTime,
			Visibility:      image.Visibility,
		})
	}
	return &ImageShareListResponse{
		Images:       images,
		PageNo:       resp.PageNo,
		CurrentPage:  resp.CurrentPage,
		PageSize:     resp.PageSize,
		CurrentCount: resp.CurrentCount,
		TotalCount:   resp.TotalCount,
	}, nil
}

type imageShareListRealResponse struct {
	Images []struct {
		Architecture    string `json:"architecture"`
		BootMode        string `json:"bootMode"`
		ContainerFormat string `json:"containerFormat"`
		CreatedTime     int    `json:"createdTime"`
		Description     string `json:"description"`
		DestinationUser string `json:"destinationUser"`
		DiskFormat      string `json:"diskFormat"`
		DiskID          string `json:"diskID"`
		DiskSize        int    `json:"diskSize"`
		ImageClass      string `json:"imageClass"`
		ImageID         string `json:"imageID"`
		ImageName       string `json:"imageName"`
		ImageType       string `json:"imageType"`
		MaximumRAM      int    `json:"maximumRAM"`
		MinimumRAM      int    `json:"minimumRAM"`
		OsDistro        string `json:"osDistro"`
		OsType          string `json:"osType"`
		OsVersion       string `json:"osVersion"`
		ProjectID       string `json:"projectID"`
		Size            int64  `json:"size"`
		SourceServerID  string `json:"sourceServerID"`
		SourceUser      string `json:"sourceUser"`
		Status          string `json:"status"`
		UpdatedTime     int    `json:"updatedTime"`
		Visibility      string `json:"visibility"`
	} `json:"images"`
	PageNo       int `json:"pageNo"`
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
}

type ImageShareListRequest struct {
	ImageId  string
	RegionId string
	PageNo   int
	PageSize int
}

type ImageShareListImagesResponse struct {
	Architecture    string
	BootMode        string
	ContainerFormat string
	CreatedTime     int
	Description     string
	DestinationUser string
	DiskFormat      string
	DiskId          string
	DiskSize        int
	ImageClass      string
	ImageId         string
	ImageName       string
	ImageType       string
	MaximumRam      int
	MinimumRam      int
	OsDistro        string
	OsType          string
	OsVersion       string
	ProjectId       string
	Size            int64
	SourceServerId  string
	SourceUser      string
	Status          string
	UpdatedTime     int
	Visibility      string
}

type ImageShareListResponse struct {
	Images       []ImageShareListImagesResponse
	PageNo       int
	CurrentPage  int
	PageSize     int
	CurrentCount int
	TotalCount   int
}
