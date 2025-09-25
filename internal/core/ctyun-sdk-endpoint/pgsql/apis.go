package pgsql

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

type Apis struct {
	PgsqlDestroyApi             *TeledbDestroyApi
	PgsqlCreateApi              *PgsqlCreateApi
	PgsqlDetailApi              *PgsqlDetailApi
	PgsqlListApi                *PgsqlListApi
	PgsqlRefundApi              *PgsqlRefundApi
	PgsqlRestartApi             *PgsqlRestartApi
	PgsqlSecurityGroupListApi   *PgsqlSecurityGroupListApi
	PgsqlSpecsApi               *PgsqlSpecsApi
	PgsqlStartApi               *PgsqlStartApi
	PgsqlStopApi                *PgsqlStopApi
	PgsqlUpdateInstanceNameApi  *PgsqlUpdateInstanceNameApi
	PgsqlUpdateSecurityGroupApi *PgsqlUpdateSecurityGroupApi
	PgsqlUpgradeApi             *PgsqlUpgradeApi
	PgsqlBindEipApi             *PgsqlBindEipApi
	PgsqlUnBindEipApi           *PgsqlUnBindEipApi
	PgsqlBoundEipListApi        *PgsqlBoundEipListApi
	PgsqlDeleteSecurityGroupApi *PgsqlDeleteSecurityGroupApi
	PgsqlGetNodeListApi         *PgsqlGetNodeListApi
}

func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	builder := ctyunsdk.NewApiHookBuilder()
	for _, hook := range client.Config.ApiHooks {
		builder.AddHooks(hook)
	}

	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointPgSqlTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointPgSqlTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndPointPgSqlProd)
	return &Apis{
		PgsqlDestroyApi:             NewTeledbDestroyApi(client),
		PgsqlCreateApi:              NewPgsqlCreateApi(client),
		PgsqlDetailApi:              NewPgsqlDetailApi(client),
		PgsqlListApi:                NewPgsqlListApi(client),
		PgsqlRefundApi:              NewPgsqlRefundApi(client),
		PgsqlRestartApi:             NewPgsqlRestartApi(client),
		PgsqlSecurityGroupListApi:   NewPgsqlSecurityGroupListApi(client),
		PgsqlSpecsApi:               NewPgsqlSpecsApi(client),
		PgsqlStartApi:               NewPgsqlStartApi(client),
		PgsqlStopApi:                NewPgsqlStopApi(client),
		PgsqlUpdateInstanceNameApi:  NewPgsqlUpdateInstanceNameApi(client),
		PgsqlUpdateSecurityGroupApi: NewPgsqlUpdateSecurityGroupApi(client),
		PgsqlUpgradeApi:             NewPgsqlUpgradeApi(client),
		PgsqlBindEipApi:             NewPgsqlBindEipApi(client),
		PgsqlUnBindEipApi:           NewPgsqlUnBindEipApi(client),
		PgsqlBoundEipListApi:        NewPgsqlBoundEipListApi(client),
		PgsqlDeleteSecurityGroupApi: NewPgsqlDeleteSecurityGroupApi(client),
		PgsqlGetNodeListApi:         NewPgsqlGetNodeListApi(client),
	}

}
