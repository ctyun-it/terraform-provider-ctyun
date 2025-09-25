package hpfs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "hpfs"

type Apis struct {
	HpfsNewSfsApi                *HpfsNewSfsApi
	HpfsInfoSfsApi               *HpfsInfoSfsApi
	HpfsInfoByNameSfsApi         *HpfsInfoByNameSfsApi
	HpfsRefundSfsApi             *HpfsRefundSfsApi
	HpfsResizeSfsApi             *HpfsResizeSfsApi
	HpfsListClusterApi           *HpfsListClusterApi
	HpfsListBaselineApi          *HpfsListBaselineApi
	HpfsListRegionApi            *HpfsListRegionApi
	HpfsListZoneApi              *HpfsListZoneApi
	HpfsCountQuotaApi            *HpfsCountQuotaApi
	HpfsCapacityQuotaApi         *HpfsCapacityQuotaApi
	HpfsDataflowQuotaApi         *HpfsDataflowQuotaApi
	HpfsRenameSfsApi             *HpfsRenameSfsApi
	HpfsInfoDirectoryApi         *HpfsInfoDirectoryApi
	HpfsNewDirectoryApi          *HpfsNewDirectoryApi
	HpfsListSfsByClusterApi      *HpfsListSfsByClusterApi
	HpfsListSfsBySfstypeApi      *HpfsListSfsBySfstypeApi
	HpfsListClusterByDeviceApi   *HpfsListClusterByDeviceApi
	HpfsListSfsApi               *HpfsListSfsApi
	HpfsListDataflowApi          *HpfsListDataflowApi
	HpfsInfoDataflowApi          *HpfsInfoDataflowApi
	HpfsNewDataflowApi           *HpfsNewDataflowApi
	HpfsUpdateDataflowApi        *HpfsUpdateDataflowApi
	HpfsDeleteDataflowApi        *HpfsDeleteDataflowApi
	HpfsNewDataflowtaskApi       *HpfsNewDataflowtaskApi
	HpfsListDataflowtaskApi      *HpfsListDataflowtaskApi
	HpfsInfoDataflowtaskApi      *HpfsInfoDataflowtaskApi
	HpfsNewProtocolServiceApi    *HpfsNewProtocolServiceApi
	HpfsDeleteProtocolServiceApi *HpfsDeleteProtocolServiceApi
	HpfsListProtocolServiceApi   *HpfsListProtocolServiceApi
	HpfsInfoProtocolServiceApi   *HpfsInfoProtocolServiceApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		HpfsNewSfsApi:                NewHpfsNewSfsApi(client),
		HpfsInfoSfsApi:               NewHpfsInfoSfsApi(client),
		HpfsInfoByNameSfsApi:         NewHpfsInfoByNameSfsApi(client),
		HpfsRefundSfsApi:             NewHpfsRefundSfsApi(client),
		HpfsResizeSfsApi:             NewHpfsResizeSfsApi(client),
		HpfsListClusterApi:           NewHpfsListClusterApi(client),
		HpfsListBaselineApi:          NewHpfsListBaselineApi(client),
		HpfsListRegionApi:            NewHpfsListRegionApi(client),
		HpfsListZoneApi:              NewHpfsListZoneApi(client),
		HpfsCountQuotaApi:            NewHpfsCountQuotaApi(client),
		HpfsCapacityQuotaApi:         NewHpfsCapacityQuotaApi(client),
		HpfsDataflowQuotaApi:         NewHpfsDataflowQuotaApi(client),
		HpfsRenameSfsApi:             NewHpfsRenameSfsApi(client),
		HpfsInfoDirectoryApi:         NewHpfsInfoDirectoryApi(client),
		HpfsNewDirectoryApi:          NewHpfsNewDirectoryApi(client),
		HpfsListSfsByClusterApi:      NewHpfsListSfsByClusterApi(client),
		HpfsListSfsBySfstypeApi:      NewHpfsListSfsBySfstypeApi(client),
		HpfsListClusterByDeviceApi:   NewHpfsListClusterByDeviceApi(client),
		HpfsListSfsApi:               NewHpfsListSfsApi(client),
		HpfsListDataflowApi:          NewHpfsListDataflowApi(client),
		HpfsInfoDataflowApi:          NewHpfsInfoDataflowApi(client),
		HpfsNewDataflowApi:           NewHpfsNewDataflowApi(client),
		HpfsUpdateDataflowApi:        NewHpfsUpdateDataflowApi(client),
		HpfsDeleteDataflowApi:        NewHpfsDeleteDataflowApi(client),
		HpfsNewDataflowtaskApi:       NewHpfsNewDataflowtaskApi(client),
		HpfsListDataflowtaskApi:      NewHpfsListDataflowtaskApi(client),
		HpfsInfoDataflowtaskApi:      NewHpfsInfoDataflowtaskApi(client),
		HpfsNewProtocolServiceApi:    NewHpfsNewProtocolServiceApi(client),
		HpfsDeleteProtocolServiceApi: NewHpfsDeleteProtocolServiceApi(client),
		HpfsListProtocolServiceApi:   NewHpfsListProtocolServiceApi(client),
		HpfsInfoProtocolServiceApi:   NewHpfsInfoProtocolServiceApi(client),
	}
}
