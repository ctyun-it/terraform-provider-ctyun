package mysql

import (
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type Apis struct {
	TeledbCreateApi             *TeledbCreateApi
	TeledbUpgradeApi            *TeledbUpgradeApi
	TeledbRefundApi             *TeledbRefundApi
	TeledbDestroyApi            *TeledbDestroyApi
	TeledbQueryDetailApi        *TeledbQueryDetailApi
	TeledbGetListApi            *TeledbGetListApi
	TeledbStartApi              *TeledbStartApi
	TeledbStopApi               *TeledbStopApi
	TeledbRestartApi            *TeledbRestartApi
	TeledbUpdateInstanceNameApi *TeledbUpdateInstanceNameApi
	TeledbUpdateWritePortApi    *TeledbUpdateWritePortApi
	TeledbBindEipApi            *TeledbBindEipApi
	TeledbUnbindEipApi          *TeledbUnbindEipApi
	TeledbBoundEipListApi       *TeledbBoundEipListApi
	TeledbMysqlSpecsApi         *TeledbMysqlSpecsApi
	TeledbCreateAccessWhiteList *TeledbCreateAccessWhiteList
	TeledbUpdateAccessWhiteList *TeledbUpdateAccessWhiteList
	TeledbDeleteAccessWhiteList *TeledbDeleteAccessWhiteList
	TeledbGetAccessWhiteList    *TeledbGetAccessWhiteList
	TeledbGetAvailabilityZone   *TeledbGetAvailabilityZone
}

func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}

	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtdasTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtdasTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndPointCtdasProd)
	return &Apis{
		TeledbCreateApi:             NewTeledbCreateApi(client),
		TeledbUpgradeApi:            NewTeledbUpgradeApi(client),
		TeledbRefundApi:             NewTeledbRefundApi(client),
		TeledbDestroyApi:            NewTeledbDestroyApi(client),
		TeledbQueryDetailApi:        NewTeledbQueryDetailApi(client),
		TeledbGetListApi:            NewTeledbGetListApi(client),
		TeledbStartApi:              NewTeledbStartApi(client),
		TeledbStopApi:               NewTeledbStopApi(client),
		TeledbRestartApi:            NewTeledbRestartApi(client),
		TeledbUpdateInstanceNameApi: NewTeledbUpdateInstanceNameApi(client),
		TeledbUpdateWritePortApi:    NewTeledbUpdateWritePortApi(client),
		TeledbBindEipApi:            NewTeledbBindEipApi(client),
		TeledbUnbindEipApi:          NewTeledbUnbindEipApi(client),
		TeledbBoundEipListApi:       NewTeledbBoundEipListApi(client),
		TeledbMysqlSpecsApi:         NewTeledbMysqlSpecsApi(client),
		TeledbCreateAccessWhiteList: NewTeledbCreateAccessWhiteList(client),
		TeledbDeleteAccessWhiteList: NewTeledbDeleteAccessWhiteList(client),
		TeledbGetAccessWhiteList:    NewTeledbGetAccessWhiteList(client),
		TeledbUpdateAccessWhiteList: NewTeledbUpdateAccessWhiteList(client),
		TeledbGetAvailabilityZone:   NewTeledbGetAvailabilityZone(client),
	}
}
