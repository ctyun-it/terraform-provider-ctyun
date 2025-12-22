package ec

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ec"

type Apis struct {
	EcEcRuleCreateAreaApi                    *EcEcRuleCreateAreaApi
	EcEcRuleListAreasApi                     *EcEcRuleListAreasApi
	EcEcRuleUpdateAreasRateApi               *EcEcRuleUpdateAreasRateApi
	EcEcRuleDeleteAreaApi                    *EcEcRuleDeleteAreaApi
	EcCreateInstanceBandwidthPacketApi       *EcCreateInstanceBandwidthPacketApi
	EcCreateInstanceIimitRateApi             *EcCreateInstanceIimitRateApi
	EcDeleteInstanceBandwidthPacketApi       *EcDeleteInstanceBandwidthPacketApi
	EcDeleteInstanceLimitRateApi             *EcDeleteInstanceLimitRateApi
	EcUpdateInstanceBandwidthPacketApi       *EcUpdateInstanceBandwidthPacketApi
	EcUpdateInstanceLimitRateApi             *EcUpdateInstanceLimitRateApi
	EcLIstInstanceBandwidthPacketApi         *EcLIstInstanceBandwidthPacketApi
	EcListInstanceLimitRateApi               *EcListInstanceLimitRateApi
	EcEcListApi                              *EcEcListApi
	EcEcCreateApi                            *EcEcCreateApi
	EcEcUpdateApi                            *EcEcUpdateApi
	EcEcDeleteApi                            *EcEcDeleteApi
	EcEcCreateGatewayApi                     *EcEcCreateGatewayApi
	EcEcUpdateGatewayApi                     *EcEcUpdateGatewayApi
	EcEcDeleteGatewayApi                     *EcEcDeleteGatewayApi
	EcEcListGatewayApi                       *EcEcListGatewayApi
	EcEcCreateRouteApi                       *EcEcCreateRouteApi
	EcEcDeleteRouteApi                       *EcEcDeleteRouteApi
	EcEcListRouteApi                         *EcEcListRouteApi
	EcEcOrderPacketNewApi                    *EcEcOrderPacketNewApi
	EcEcOrderPacketUpgradeApi                *EcEcOrderPacketUpgradeApi
	EcEcOrderPacketRenewApi                  *EcEcOrderPacketRenewApi
	EcEcOrderPacketRefundApi                 *EcEcOrderPacketRefundApi
	EcEcQueryPacketUpgradePriceApi           *EcEcQueryPacketUpgradePriceApi
	EcEcQueryPacketRenewPriceApi             *EcEcQueryPacketRenewPriceApi
	EcEcQueryPacketNewPriceApi               *EcEcQueryPacketNewPriceApi
	EcEcPacketListPacketApi                  *EcEcPacketListPacketApi
	EcEcAddVPCNetworkApi                     *EcEcAddVPCNetworkApi
	EcEcUpdateVPCNetworkApi                  *EcEcUpdateVPCNetworkApi
	EcEcDeleteVPCNetworkApi                  *EcEcDeleteVPCNetworkApi
	EcEcListVPCNetworkApi                    *EcEcListVPCNetworkApi
	EcEcUpdateInstanceRouteApi               *EcEcUpdateInstanceRouteApi
	EcEcListAuthVPCBindCloudHighApi          *EcEcListAuthVPCBindCloudHighApi
	EcEcAddCDANetworkApi                     *EcEcAddCDANetworkApi
	EcEcUpdateCDANetworkApi                  *EcEcUpdateCDANetworkApi
	EcEcDeleteCDANetworkApi                  *EcEcDeleteCDANetworkApi
	EcEcListCDANetworkApi                    *EcEcListCDANetworkApi
	EcEcCreateRouteTableAutoLearnInstanceApi *EcEcCreateRouteTableAutoLearnInstanceApi
	EcEcCreateSDWANInstanceApi               *EcEcCreateSDWANInstanceApi
	EcEcListSDWANInstanceApi                 *EcEcListSDWANInstanceApi
	EcEcUpdateSDWANInstanceApi               *EcEcUpdateSDWANInstanceApi
	EcEcDeleteSDWANInstanceApi               *EcEcDeleteSDWANInstanceApi
	EcEcUpdateWeightsApi                     *EcEcUpdateWeightsApi
	EcEcListSDWANApi                         *EcEcListSDWANApi
	EcEcCheckSDWANApi                        *EcEcCheckSDWANApi
	EcEcBindSDWANApi                         *EcEcBindSDWANApi
	EcEcListCloudHighSubnetApi               *EcEcListCloudHighSubnetApi
	EcEcQueryRemainQuotaApi                  *EcEcQueryRemainQuotaApi
	// 新增云网关计费API
	EcEcCgwBillNewApi    *EcEcCgwBillNewApi
	EcEcCgwBillRefundApi *EcEcCgwBillRefundApi
	// 按需订单查询API
	EcEcTgwOrderQueryApi    *EcEcTgwOrderQueryApi
	EcCreateRegionPeerApi   *EcCreateRegionPeerApi
	EcDeleteRegionPeerApi   *EcDeleteRegionPeerApi
	EcRegionPeerListApi     *EcRegionPeerListApi
	EcRegionPeerUpdateApi   *EcRegionPeerUpdateApi
	EcEcBindCloudGatewayApi *EcEcBindCloudGatewayApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		EcEcRuleCreateAreaApi:                    NewEcEcRuleCreateAreaApi(client),
		EcEcRuleListAreasApi:                     NewEcEcRuleListAreasApi(client),
		EcEcRuleUpdateAreasRateApi:               NewEcEcRuleUpdateAreasRateApi(client),
		EcEcRuleDeleteAreaApi:                    NewEcEcRuleDeleteAreaApi(client),
		EcCreateInstanceBandwidthPacketApi:       NewEcCreateInstanceBandwidthPacketApi(client),
		EcCreateInstanceIimitRateApi:             NewEcCreateInstanceIimitRateApi(client),
		EcDeleteInstanceBandwidthPacketApi:       NewEcDeleteInstanceBandwidthPacketApi(client),
		EcDeleteInstanceLimitRateApi:             NewEcDeleteInstanceLimitRateApi(client),
		EcUpdateInstanceBandwidthPacketApi:       NewEcUpdateInstanceBandwidthPacketApi(client),
		EcUpdateInstanceLimitRateApi:             NewEcUpdateInstanceLimitRateApi(client),
		EcLIstInstanceBandwidthPacketApi:         NewEcLIstInstanceBandwidthPacketApi(client),
		EcListInstanceLimitRateApi:               NewEcListInstanceLimitRateApi(client),
		EcEcListApi:                              NewEcEcListApi(client),
		EcEcCreateApi:                            NewEcEcCreateApi(client),
		EcEcUpdateApi:                            NewEcEcUpdateApi(client),
		EcEcDeleteApi:                            NewEcEcDeleteApi(client),
		EcEcCreateGatewayApi:                     NewEcEcCreateGatewayApi(client),
		EcEcUpdateGatewayApi:                     NewEcEcUpdateGatewayApi(client),
		EcEcDeleteGatewayApi:                     NewEcEcDeleteGatewayApi(client),
		EcEcListGatewayApi:                       NewEcEcListGatewayApi(client),
		EcEcCreateRouteApi:                       NewEcEcCreateRouteApi(client),
		EcEcDeleteRouteApi:                       NewEcEcDeleteRouteApi(client),
		EcEcListRouteApi:                         NewEcEcListRouteApi(client),
		EcEcOrderPacketNewApi:                    NewEcEcOrderPacketNewApi(client),
		EcEcOrderPacketUpgradeApi:                NewEcEcOrderPacketUpgradeApi(client),
		EcEcOrderPacketRenewApi:                  NewEcEcOrderPacketRenewApi(client),
		EcEcOrderPacketRefundApi:                 NewEcEcOrderPacketRefundApi(client),
		EcEcQueryPacketUpgradePriceApi:           NewEcEcQueryPacketUpgradePriceApi(client),
		EcEcQueryPacketRenewPriceApi:             NewEcEcQueryPacketRenewPriceApi(client),
		EcEcQueryPacketNewPriceApi:               NewEcEcQueryPacketNewPriceApi(client),
		EcEcPacketListPacketApi:                  NewEcEcPacketListPacketApi(client),
		EcEcAddVPCNetworkApi:                     NewEcEcAddVPCNetworkApi(client),
		EcEcUpdateVPCNetworkApi:                  NewEcEcUpdateVPCNetworkApi(client),
		EcEcDeleteVPCNetworkApi:                  NewEcEcDeleteVPCNetworkApi(client),
		EcEcListVPCNetworkApi:                    NewEcEcListVPCNetworkApi(client),
		EcEcUpdateInstanceRouteApi:               NewEcEcUpdateInstanceRouteApi(client),
		EcEcListAuthVPCBindCloudHighApi:          NewEcEcListAuthVPCBindCloudHighApi(client),
		EcEcAddCDANetworkApi:                     NewEcEcAddCDANetworkApi(client),
		EcEcUpdateCDANetworkApi:                  NewEcEcUpdateCDANetworkApi(client),
		EcEcDeleteCDANetworkApi:                  NewEcEcDeleteCDANetworkApi(client),
		EcEcListCDANetworkApi:                    NewEcEcListCDANetworkApi(client),
		EcEcCreateRouteTableAutoLearnInstanceApi: NewEcEcCreateRouteTableAutoLearnInstanceApi(client),
		EcEcCreateSDWANInstanceApi:               NewEcEcCreateSDWANInstanceApi(client),
		EcEcListSDWANInstanceApi:                 NewEcEcListSDWANInstanceApi(client),
		EcEcUpdateSDWANInstanceApi:               NewEcEcUpdateSDWANInstanceApi(client),
		EcEcDeleteSDWANInstanceApi:               NewEcEcDeleteSDWANInstanceApi(client),
		EcEcUpdateWeightsApi:                     NewEcEcUpdateWeightsApi(client),
		EcEcListSDWANApi:                         NewEcEcListSDWANApi(client),
		EcEcCheckSDWANApi:                        NewEcEcCheckSDWANApi(client),
		EcEcBindSDWANApi:                         NewEcEcBindSDWANApi(client),
		EcEcListCloudHighSubnetApi:               NewEcEcListCloudHighSubnetApi(client),
		EcEcQueryRemainQuotaApi:                  NewEcEcQueryRemainQuotaApi(client),
		// 注册云网关计费API
		EcEcCgwBillNewApi:    NewEcEcCgwBillNewApi(client),
		EcEcCgwBillRefundApi: NewEcEcCgwBillRefundApi(client),
		// 注册按需订单查询API
		EcEcTgwOrderQueryApi:    NewEcEcTgwOrderQueryApi(client),
		EcCreateRegionPeerApi:   NewEcCreateRegionPeerApi(client),
		EcDeleteRegionPeerApi:   NewEcDeleteRegionPeerApi(client),
		EcRegionPeerListApi:     NewEcRegionPeerListApi(client),
		EcRegionPeerUpdateApi:   NewEcRegionPeerUpdateApi(client),
		EcEcBindCloudGatewayApi: NewEcEcBindCloudGatewayApi(client),
	}
}
