package ctimage

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ctimage"

type Apis struct {
	CtimageListImagesApi                       *CtimageListImagesApi
	CtimageDetailImageApi                      *CtimageDetailImageApi
	CtimageCreateEcsSystemDiskImageApi         *CtimageCreateEcsSystemDiskImageApi
	CtimageDeleteImageApi                      *CtimageDeleteImageApi
	CtimageImportImageApi                      *CtimageImportImageApi
	CtimageDeactivatePrivateImageApi           *CtimageDeactivatePrivateImageApi
	CtimageCreateEcsDataDiskImageApi           *CtimageCreateEcsDataDiskImageApi
	CtimageListImageSharesApi                  *CtimageListImageSharesApi
	CtimageDeleteImageImportTaskApi            *CtimageDeleteImageImportTaskApi
	CtimageListImageImportTasksApi             *CtimageListImageImportTasksApi
	CtimageReactivateImageApi                  *CtimageReactivateImageApi
	CtimageExportImageApi                      *CtimageExportImageApi
	CtimageValidateImageFileSrcApi             *CtimageValidateImageFileSrcApi
	CtimageUnshareImageApi                     *CtimageUnshareImageApi
	CtimageRejectImageApi                      *CtimageRejectImageApi
	CtimageAcceptImageApi                      *CtimageAcceptImageApi
	CtimageShareImageApi                       *CtimageShareImageApi
	CtimageUpdateImageApi                      *CtimageUpdateImageApi
	CtimageCreateEcsSnapshotSystemDiskImageApi *CtimageCreateEcsSnapshotSystemDiskImageApi
	CtimageCreateFullEcsImageApi               *CtimageCreateFullEcsImageApi
	CtimageListSupportedDestCrrRegionsApi      *CtimageListSupportedDestCrrRegionsApi
	CtimageCopyImageAcrossRegionsApi           *CtimageCopyImageAcrossRegionsApi
	CtimageCopyImageApi                        *CtimageCopyImageApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CtimageListImagesApi:                       NewCtimageListImagesApi(client),
		CtimageDetailImageApi:                      NewCtimageDetailImageApi(client),
		CtimageCreateEcsSystemDiskImageApi:         NewCtimageCreateEcsSystemDiskImageApi(client),
		CtimageDeleteImageApi:                      NewCtimageDeleteImageApi(client),
		CtimageImportImageApi:                      NewCtimageImportImageApi(client),
		CtimageDeactivatePrivateImageApi:           NewCtimageDeactivatePrivateImageApi(client),
		CtimageCreateEcsDataDiskImageApi:           NewCtimageCreateEcsDataDiskImageApi(client),
		CtimageListImageSharesApi:                  NewCtimageListImageSharesApi(client),
		CtimageDeleteImageImportTaskApi:            NewCtimageDeleteImageImportTaskApi(client),
		CtimageListImageImportTasksApi:             NewCtimageListImageImportTasksApi(client),
		CtimageReactivateImageApi:                  NewCtimageReactivateImageApi(client),
		CtimageExportImageApi:                      NewCtimageExportImageApi(client),
		CtimageValidateImageFileSrcApi:             NewCtimageValidateImageFileSrcApi(client),
		CtimageUnshareImageApi:                     NewCtimageUnshareImageApi(client),
		CtimageRejectImageApi:                      NewCtimageRejectImageApi(client),
		CtimageAcceptImageApi:                      NewCtimageAcceptImageApi(client),
		CtimageShareImageApi:                       NewCtimageShareImageApi(client),
		CtimageUpdateImageApi:                      NewCtimageUpdateImageApi(client),
		CtimageCreateEcsSnapshotSystemDiskImageApi: NewCtimageCreateEcsSnapshotSystemDiskImageApi(client),
		CtimageCreateFullEcsImageApi:               NewCtimageCreateFullEcsImageApi(client),
		CtimageListSupportedDestCrrRegionsApi:      NewCtimageListSupportedDestCrrRegionsApi(client),
		CtimageCopyImageAcrossRegionsApi:           NewCtimageCopyImageAcrossRegionsApi(client),
		CtimageCopyImageApi:                        NewCtimageCopyImageApi(client),
	}
}
