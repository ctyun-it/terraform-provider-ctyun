package ctebm

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ebm"

type Apis struct {
	EbmResetHostnameApi                      *EbmResetHostnameApi
	EbmBatchDeleteInstancesApi               *EbmBatchDeleteInstancesApi
	EbmDestroyInstanceApi                    *EbmDestroyInstanceApi
	EbmMetadataBatchUpdateApi                *EbmMetadataBatchUpdateApi
	EbmMetadataBatchCreateApi                *EbmMetadataBatchCreateApi
	EbmMetadataDeleteApi                     *EbmMetadataDeleteApi
	EbmMetadataDeleteAllApi                  *EbmMetadataDeleteAllApi
	EbmMetadataUpdateApi                     *EbmMetadataUpdateApi
	EbmMetadataCreateApi                     *EbmMetadataCreateApi
	EbmMetadataListApi                       *EbmMetadataListApi
	EbmCreateInstanceV4plusApi               *EbmCreateInstanceV4plusApi
	EbmListInstanceV4plusApi                 *EbmListInstanceV4plusApi
	EbmDescribeInstanceV4plusApi             *EbmDescribeInstanceV4plusApi
	EbmMonitorHistoryDiskMetricApi           *EbmMonitorHistoryDiskMetricApi
	EbmMonitorHistoryInterfaceMetricApi      *EbmMonitorHistoryInterfaceMetricApi
	EbmMonitorHistoryMemoryMetricApi         *EbmMonitorHistoryMemoryMetricApi
	EbmMonitorHistoryCpuMetricApi            *EbmMonitorHistoryCpuMetricApi
	EbmMonitorLatestDiskMetricApi            *EbmMonitorLatestDiskMetricApi
	EbmMonitorLatestInterfaceMetricApi       *EbmMonitorLatestInterfaceMetricApi
	EbmMonitorLatestMemoryMetricApi          *EbmMonitorLatestMemoryMetricApi
	EbmMonitorLatestCpuMetricApi             *EbmMonitorLatestCpuMetricApi
	EbmInstanceRaidListApi                   *EbmInstanceRaidListApi
	EbmInstanceInterfaceSecurityGroupListApi *EbmInstanceInterfaceSecurityGroupListApi
	EbmVncApi                                *EbmVncApi
	EbmUpdateSecurityGroupApi                *EbmUpdateSecurityGroupApi
	EbmDeleteSecurityGroupApi                *EbmDeleteSecurityGroupApi
	EbmAddSecurityGroupApi                   *EbmAddSecurityGroupApi
	EbmRemoveNicApi                          *EbmRemoveNicApi
	EbmAddNicApi                             *EbmAddNicApi
	EbmDetachVolumeApi                       *EbmDetachVolumeApi
	EbmAttachVolumeApi                       *EbmAttachVolumeApi
	EbmInstanceInterfaceListApi              *EbmInstanceInterfaceListApi
	EbmInstanceAttachedVolumeIdListApi       *EbmInstanceAttachedVolumeIdListApi
	EbmInstanceDeviceTypeApi                 *EbmInstanceDeviceTypeApi
	EbmInstanceImageApi                      *EbmInstanceImageApi
	EbmUpdateInstanceApi                     *EbmUpdateInstanceApi
	EbmRenewInstanceApi                      *EbmRenewInstanceApi
	EbmDeleteInstanceApi                     *EbmDeleteInstanceApi
	EbmCreateInstanceApi                     *EbmCreateInstanceApi
	EbmResetPasswordApi                      *EbmResetPasswordApi
	EbmDescribeInstanceApi                   *EbmDescribeInstanceApi
	EbmReinstallInstanceApi                  *EbmReinstallInstanceApi
	EbmRebootInstanceApi                     *EbmRebootInstanceApi
	EbmStopInstanceApi                       *EbmStopInstanceApi
	EbmStartInstanceApi                      *EbmStartInstanceApi
	EbmListInstanceApi                       *EbmListInstanceApi
	EbmImageListApi                          *EbmImageListApi
	EbmRaidTypeListApi                       *EbmRaidTypeListApi
	EbmDeviceStockListApi                    *EbmDeviceStockListApi
	EbmDeviceTypeListApi                     *EbmDeviceTypeListApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		EbmResetHostnameApi:                      NewEbmResetHostnameApi(client),
		EbmBatchDeleteInstancesApi:               NewEbmBatchDeleteInstancesApi(client),
		EbmDestroyInstanceApi:                    NewEbmDestroyInstanceApi(client),
		EbmMetadataBatchUpdateApi:                NewEbmMetadataBatchUpdateApi(client),
		EbmMetadataBatchCreateApi:                NewEbmMetadataBatchCreateApi(client),
		EbmMetadataDeleteApi:                     NewEbmMetadataDeleteApi(client),
		EbmMetadataDeleteAllApi:                  NewEbmMetadataDeleteAllApi(client),
		EbmMetadataUpdateApi:                     NewEbmMetadataUpdateApi(client),
		EbmMetadataCreateApi:                     NewEbmMetadataCreateApi(client),
		EbmMetadataListApi:                       NewEbmMetadataListApi(client),
		EbmCreateInstanceV4plusApi:               NewEbmCreateInstanceV4plusApi(client),
		EbmListInstanceV4plusApi:                 NewEbmListInstanceV4plusApi(client),
		EbmDescribeInstanceV4plusApi:             NewEbmDescribeInstanceV4plusApi(client),
		EbmMonitorHistoryDiskMetricApi:           NewEbmMonitorHistoryDiskMetricApi(client),
		EbmMonitorHistoryInterfaceMetricApi:      NewEbmMonitorHistoryInterfaceMetricApi(client),
		EbmMonitorHistoryMemoryMetricApi:         NewEbmMonitorHistoryMemoryMetricApi(client),
		EbmMonitorHistoryCpuMetricApi:            NewEbmMonitorHistoryCpuMetricApi(client),
		EbmMonitorLatestDiskMetricApi:            NewEbmMonitorLatestDiskMetricApi(client),
		EbmMonitorLatestInterfaceMetricApi:       NewEbmMonitorLatestInterfaceMetricApi(client),
		EbmMonitorLatestMemoryMetricApi:          NewEbmMonitorLatestMemoryMetricApi(client),
		EbmMonitorLatestCpuMetricApi:             NewEbmMonitorLatestCpuMetricApi(client),
		EbmInstanceRaidListApi:                   NewEbmInstanceRaidListApi(client),
		EbmInstanceInterfaceSecurityGroupListApi: NewEbmInstanceInterfaceSecurityGroupListApi(client),
		EbmVncApi:                                NewEbmVncApi(client),
		EbmUpdateSecurityGroupApi:                NewEbmUpdateSecurityGroupApi(client),
		EbmDeleteSecurityGroupApi:                NewEbmDeleteSecurityGroupApi(client),
		EbmAddSecurityGroupApi:                   NewEbmAddSecurityGroupApi(client),
		EbmRemoveNicApi:                          NewEbmRemoveNicApi(client),
		EbmAddNicApi:                             NewEbmAddNicApi(client),
		EbmDetachVolumeApi:                       NewEbmDetachVolumeApi(client),
		EbmAttachVolumeApi:                       NewEbmAttachVolumeApi(client),
		EbmInstanceInterfaceListApi:              NewEbmInstanceInterfaceListApi(client),
		EbmInstanceAttachedVolumeIdListApi:       NewEbmInstanceAttachedVolumeIdListApi(client),
		EbmInstanceDeviceTypeApi:                 NewEbmInstanceDeviceTypeApi(client),
		EbmInstanceImageApi:                      NewEbmInstanceImageApi(client),
		EbmUpdateInstanceApi:                     NewEbmUpdateInstanceApi(client),
		EbmRenewInstanceApi:                      NewEbmRenewInstanceApi(client),
		EbmDeleteInstanceApi:                     NewEbmDeleteInstanceApi(client),
		EbmCreateInstanceApi:                     NewEbmCreateInstanceApi(client),
		EbmResetPasswordApi:                      NewEbmResetPasswordApi(client),
		EbmDescribeInstanceApi:                   NewEbmDescribeInstanceApi(client),
		EbmReinstallInstanceApi:                  NewEbmReinstallInstanceApi(client),
		EbmRebootInstanceApi:                     NewEbmRebootInstanceApi(client),
		EbmStopInstanceApi:                       NewEbmStopInstanceApi(client),
		EbmStartInstanceApi:                      NewEbmStartInstanceApi(client),
		EbmListInstanceApi:                       NewEbmListInstanceApi(client),
		EbmImageListApi:                          NewEbmImageListApi(client),
		EbmRaidTypeListApi:                       NewEbmRaidTypeListApi(client),
		EbmDeviceStockListApi:                    NewEbmDeviceStockListApi(client),
		EbmDeviceTypeListApi:                     NewEbmDeviceTypeListApi(client),
	}
}
