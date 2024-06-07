package ctimage

import (
	"context"
	"net/http"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// ImageImportApi 创建私有镜像（镜像文件）
// https://www.ctyun.cn/document/10027726/10040048
type ImageImportApi struct {
	ctyunsdk.CtyunRequestBuilder
	client *ctyunsdk.CtyunClient
}

func NewImageImportApi(client *ctyunsdk.CtyunClient) *ImageImportApi {
	return &ImageImportApi{
		client: client,
		CtyunRequestBuilder: ctyunsdk.CtyunRequestBuilder{
			Method:  http.MethodPost,
			UrlPath: "/v4/image/import",
		},
	}
}

func (this *ImageImportApi) Do(ctx context.Context, credential ctyunsdk.Credential, req *ImageImportRequest) (*ImageImportResponse, ctyunsdk.CtyunRequestError) {
	request := this.WithCredential(&credential)
	_, err := request.WriteJson(&imageImportRealRequest{
		RegionID:        req.RegionId,
		ImageFileSource: req.ImageFileSource,
		ImageProperties: imageImportImagePropertiesRealRequest{
			ImageName:    req.ImageProperties.ImageName,
			OsDistro:     req.ImageProperties.OsDistro,
			OsVersion:    req.ImageProperties.OsVersion,
			Architecture: req.ImageProperties.Architecture,
			BootMode:     req.ImageProperties.BootMode,
			Description:  req.ImageProperties.Description,
			DiskSize:     req.ImageProperties.DiskSize,
			ImageType:    req.ImageProperties.ImageType,
			MaximumRam:   req.ImageProperties.MaximumRam,
			MinimumRam:   req.ImageProperties.MinimumRam,
		},
		ProjectID: req.ProjectId,
	})
	if err != nil {
		return nil, err
	}

	response, err := this.client.RequestToEndpoint(ctx, EndpointNameCtimage, request)
	if err != nil {
		return nil, err
	}

	var resp imageImportRealResponse
	err = response.ParseByStandardModelWithCheck(&resp)
	if err != nil {
		return nil, err
	}
	var images []ImageImportImagesResponse
	for _, image := range resp.Images {
		images = append(images, ImageImportImagesResponse{
			Architecture:     image.Architecture,
			AzName:           image.AzName,
			BootMode:         image.BootMode,
			ContainerFormat:  image.ContainerFormat,
			CreatedTime:      image.CreatedTime,
			Description:      image.Description,
			DestinationUser:  image.DestinationUser,
			DiskFormat:       image.DiskFormat,
			DiskId:           image.DiskID,
			DiskSize:         image.DiskSize,
			ImageClass:       image.ImageClass,
			ImageId:          image.ImageID,
			ImageName:        image.ImageName,
			ImageType:        image.ImageType,
			MaximumRam:       image.MaximumRAM,
			MinimumRam:       image.MinimumRAM,
			OsDistro:         image.OsDistro,
			OsType:           image.OsType,
			OsVersion:        image.OsVersion,
			ProjectId:        image.ProjectID,
			SharedListLength: image.SharedListLength,
			Size:             image.Size,
			SourceServerId:   image.SourceServerID,
			SourceUser:       image.SourceUser,
			Status:           image.Status,
			Tags:             image.Tags,
			UpdatedTime:      image.UpdatedTime,
			Visibility:       image.Visibility,
		})
	}
	return &ImageImportResponse{
		Images: images,
	}, nil
}

type imageImportImagePropertiesRealRequest struct {
	ImageName    string `json:"imageName"`
	OsDistro     string `json:"osDistro"`
	OsVersion    string `json:"osVersion"`
	Architecture string `json:"architecture"`
	BootMode     string `json:"bootMode"`
	Description  string `json:"description"`
	DiskSize     int    `json:"diskSize"`
	ImageType    string `json:"imageType"`
	MaximumRam   int    `json:"maximumRAM"`
	MinimumRam   int    `json:"minimumRAM"`
}

type imageImportRealRequest struct {
	RegionID        string                                `json:"regionID"`
	ImageFileSource string                                `json:"imageFileSource"`
	ImageProperties imageImportImagePropertiesRealRequest `json:"imageProperties"`
	ProjectID       string                                `json:"projectID,omitempty"`
}

type imageImportRealResponse struct {
	Images []struct {
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
		Size             int    `json:"size"`
		SourceServerID   string `json:"sourceServerID"`
		SourceUser       string `json:"sourceUser"`
		Status           string `json:"status"`
		Tags             string `json:"tags"`
		UpdatedTime      int    `json:"updatedTime"`
		Visibility       string `json:"visibility"`
	} `json:"images"`
}

type ImageImportImagePropertiesRequest struct {
	ImageName    string
	OsDistro     string
	OsVersion    string
	Architecture string
	BootMode     string
	Description  string
	DiskSize     int
	ImageType    string
	MaximumRam   int
	MinimumRam   int
}

type ImageImportRequest struct {
	RegionId        string                            // 资源池id
	ImageFileSource string                            // 镜像文件地址，格式应为 {internetEndpoint}/{bucket}/{key}。可使用访问控制 endpoint 查询接口来查询外网访问 endpoint，可使用获取桶列表接口来查询您拥有的桶的列表，可使用查看对象列表接口来查询存储桶内所有对象。
	ImageProperties ImageImportImagePropertiesRequest //  镜像属性。注意：对启动方式、最大和最小内存的指定仅在镜像属性中镜像种类的取值为空或空字符串（系统盘镜像）时生效。此外，当前对启动方式、最大和最小内存的指定在多可用区资源池下不生效。可使用资源池概况信息查询接口来确认所指定的资源池是否是多可用区资源池。
	ProjectId       string                            // 企业项目id
}

type ImageImportImagesResponse struct {
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
	Size             int
	SourceServerId   string
	SourceUser       string
	Status           string
	Tags             string
	UpdatedTime      int
	Visibility       string
}

type ImageImportResponse struct {
	Images []ImageImportImagesResponse
}
