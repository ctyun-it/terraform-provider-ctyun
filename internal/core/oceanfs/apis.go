package oceanfs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "oceanfs"

type Apis struct {
	OceanfsStorageTypeApi           *OceanfsStorageTypeApi
	OceanfsListVpcPermissionApi     *OceanfsListVpcPermissionApi
	OceanfsModifyPermissionRuleApi  *OceanfsModifyPermissionRuleApi
	OceanfsDeletePermissionRuleApi  *OceanfsDeletePermissionRuleApi
	OceanfsNewPermissionRuleApi     *OceanfsNewPermissionRuleApi
	OceanfsVpcChangePermissionApi   *OceanfsVpcChangePermissionApi
	OceanfsVpcUnbindPermissionApi   *OceanfsVpcUnbindPermissionApi
	OceanfsVpcBindPermissionApi     *OceanfsVpcBindPermissionApi
	OceanfsModifyPermissionGroupApi *OceanfsModifyPermissionGroupApi
	OceanfsDeletePermissionGroupApi *OceanfsDeletePermissionGroupApi
	OceanfsNewPermissionGroupApi    *OceanfsNewPermissionGroupApi
	OceanfsRenewSfsApi              *OceanfsRenewSfsApi
	OceanfsInfoByNameSfsApi         *OceanfsInfoByNameSfsApi
	OceanfsInfoSfsApi               *OceanfsInfoSfsApi
	OceanfsListPermissionGroupApi   *OceanfsListPermissionGroupApi
	OceanfsListPermissionRuleApi    *OceanfsListPermissionRuleApi
	OceanfsOpendListSfsApi          *OceanfsOpendListSfsApi
	OceanfsResizeSfsApi             *OceanfsResizeSfsApi
	OceanfsRefundSfsApi             *OceanfsRefundSfsApi
	OceanfsNewSfsApi                *OceanfsNewSfsApi
	OceanfsListSfsApi               *OceanfsListSfsApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		OceanfsStorageTypeApi:           NewOceanfsStorageTypeApi(client),
		OceanfsListVpcPermissionApi:     NewOceanfsListVpcPermissionApi(client),
		OceanfsModifyPermissionRuleApi:  NewOceanfsModifyPermissionRuleApi(client),
		OceanfsDeletePermissionRuleApi:  NewOceanfsDeletePermissionRuleApi(client),
		OceanfsNewPermissionRuleApi:     NewOceanfsNewPermissionRuleApi(client),
		OceanfsVpcChangePermissionApi:   NewOceanfsVpcChangePermissionApi(client),
		OceanfsVpcUnbindPermissionApi:   NewOceanfsVpcUnbindPermissionApi(client),
		OceanfsVpcBindPermissionApi:     NewOceanfsVpcBindPermissionApi(client),
		OceanfsModifyPermissionGroupApi: NewOceanfsModifyPermissionGroupApi(client),
		OceanfsDeletePermissionGroupApi: NewOceanfsDeletePermissionGroupApi(client),
		OceanfsNewPermissionGroupApi:    NewOceanfsNewPermissionGroupApi(client),
		OceanfsRenewSfsApi:              NewOceanfsRenewSfsApi(client),
		OceanfsInfoByNameSfsApi:         NewOceanfsInfoByNameSfsApi(client),
		OceanfsInfoSfsApi:               NewOceanfsInfoSfsApi(client),
		OceanfsListPermissionGroupApi:   NewOceanfsListPermissionGroupApi(client),
		OceanfsListPermissionRuleApi:    NewOceanfsListPermissionRuleApi(client),
		OceanfsOpendListSfsApi:          NewOceanfsOpendListSfsApi(client),
		OceanfsResizeSfsApi:             NewOceanfsResizeSfsApi(client),
		OceanfsRefundSfsApi:             NewOceanfsRefundSfsApi(client),
		OceanfsNewSfsApi:                NewOceanfsNewSfsApi(client),
		OceanfsListSfsApi:               NewOceanfsListSfsApi(client),
	}
}
