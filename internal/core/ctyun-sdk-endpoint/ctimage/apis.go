package ctimage

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// Apis api的接口
type Apis struct {
	ImageListApi            *ImageListApi
	ImageDetailApi          *ImageDetailApi
	ImageImportApi          *ImageImportApi
	ImageListImportTasksApi *ImageListImportTasksApi
	ImageDeleteApi          *ImageDeleteApi
	ImageUpdateApi          *ImageUpdateApi
	ImageShareCreateApi     *ImageShareCreateApi
	ImageShareDeleteApi     *ImageShareDeleteApi
	ImageShareAcceptApi     *ImageShareAcceptApi
	ImageShareRejectApi     *ImageShareRejectApi
	ImageShareListApi       *ImageShareListApi
}

// NewApis 构建
func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtimageTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentTest, EndpointCtimageTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndpointCtimageProd)
	return &Apis{
		ImageListApi:            NewImageListApi(client),
		ImageDetailApi:          NewImageDetailApi(client),
		ImageImportApi:          NewImageImportApi(client),
		ImageListImportTasksApi: NewImageListImportTasksApi(client),
		ImageDeleteApi:          NewImageDeleteApi(client),
		ImageUpdateApi:          NewImageUpdateApi(client),
		ImageShareCreateApi:     NewImageShareCreateApi(client),
		ImageShareDeleteApi:     NewImageShareDeleteApi(client),
		ImageShareAcceptApi:     NewimageShareAcceptApi(client),
		ImageShareRejectApi:     NewimageShareRejectApi(client),
		ImageShareListApi:       NewimageShareListApi(client),
	}
}
