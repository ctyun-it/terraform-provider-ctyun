package pgsql

import ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"

type Apis struct {
	PgsqlDestroyApi                       *TeledbDestroyApi
	PgsqlCreateApi                        *PgsqlCreateApi
	PgsqlDetailApi                        *PgsqlDetailApi
	PgsqlListApi                          *PgsqlListApi
	PgsqlRefundApi                        *PgsqlRefundApi
	PgsqlRestartApi                       *PgsqlRestartApi
	PgsqlSecurityGroupListApi             *PgsqlSecurityGroupListApi
	PgsqlSpecsApi                         *PgsqlSpecsApi
	PgsqlStartApi                         *PgsqlStartApi
	PgsqlStopApi                          *PgsqlStopApi
	PgsqlUpdateInstanceNameApi            *PgsqlUpdateInstanceNameApi
	PgsqlUpdateSecurityGroupApi           *PgsqlUpdateSecurityGroupApi
	PgsqlUpgradeApi                       *PgsqlUpgradeApi
	PgsqlBindEipApi                       *PgsqlBindEipApi
	PgsqlUnBindEipApi                     *PgsqlUnBindEipApi
	PgsqlBoundEipListApi                  *PgsqlBoundEipListApi
	PgsqlDeleteSecurityGroupApi           *PgsqlDeleteSecurityGroupApi
	PgsqlGetNodeListApi                   *PgsqlGetNodeListApi
	PostgresqlCollectorPolicyApi          *PostgresqlCollectorPolicyApi
	PostgresqlGetCollectorPolicyApi       *PostgresqlGetCollectorPolicyApi
	PostgresqlAddSecurityGroupApi         *PostgresqlAddSecurityGroupApi
	PostgresqlChangeSecurityGroupApi      *PostgresqlChangeSecurityGroupApi
	PgsqlCreateAccountApi                 *PgsqlCreateAccountApi
	PgsqlDeleteAccountApi                 *PgsqlDeleteAccountApi
	PgsqlGrantPrivilegeApi                *PgsqlGrantPrivilegeApi
	PgsqlRevokePrivilegeApi               *PgsqlRevokePrivilegeApi
	PgsqlLockAccountApi                   *PgsqlLockAccountApi
	PgsqlUnLockAccountApi                 *PgsqlUnLockAccountApi
	PgsqlResetPasswordApi                 *PgsqlResetPasswordApi
	PgsqlResetRootPasswordApi             *PgsqlResetRootPasswordApi
	PgsqlUpdateAccountRemarkApi           *PgsqlUpdateAccountRemarkApi
	PgsqlGetAccountListApi                *PgsqlGetAccountListApi
	PgsqlCreateDatabaseApi                *PgsqlCreateDatabaseApi
	PgsqlDeleteDatabaseApi                *PgsqlDeleteDatabaseApi
	PgsqlGetCharacterSetApi               *PgsqlGetCharacterSetApi
	PgsqlGetCollationTimeZoneApi          *PgsqlGetCollationTimeZoneApi
	PgsqlUpdateDatabaseRemarkApi          *PgsqlUpdateDatabaseRemarkApi
	PgsqlGetDatabaseSchemaApi             *PgsqlGetDatabaseSchemaApi
	PgsqlGetDatabaseSchemaListApi         *PgsqlGetDatabaseSchemaListApi
	PgsqlCreateBackupApi                  *PgsqlCreateBackupApi
	PgsqlDeleteBackupApi                  *PgsqlDeleteBackupApi
	PgsqlGetRecoverableBackupListApi      *PgsqlGetRecoverableBackupListApi
	PgsqlGetBackupTaskListApi             *PgsqlGetBackupTaskListApi
	PgsqlGetBackupListApi                 *PgsqlGetBackupListApi
	PgsqlCreateParameterTemplateApi       *PgsqlCreateParameterTemplateApi
	PgsqlDeleteParameterTemplateApi       *PgsqlDeleteParameterTemplateApi
	PgsqlGetParameterTemplateDetailApi    *PgsqlGetParameterTemplateDetailApi
	PgsqlGetParameterTemplateListApi      *PgsqlGetParameterTemplateListApi
	PgsqlUpdateParameterTemplateApi       *PgsqlUpdateParameterTemplateApi
	PgsqlUpdateParameterTemplateRemarkApi *PgsqlUpdateParameterTemplateRemarkApi
	PgsqlUpdateWhiteListApi               *PgsqlUpdateWhiteListApi
	PgsqlGetWhiteListApi                  *PgsqlGetWhiteListApi
	PgsqlGetIDByOrderApi                  *PgsqlGetIDByOrderApi
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
		PgsqlDestroyApi:                       NewTeledbDestroyApi(client),
		PgsqlCreateApi:                        NewPgsqlCreateApi(client),
		PgsqlDetailApi:                        NewPgsqlDetailApi(client),
		PgsqlListApi:                          NewPgsqlListApi(client),
		PgsqlRefundApi:                        NewPgsqlRefundApi(client),
		PgsqlRestartApi:                       NewPgsqlRestartApi(client),
		PgsqlSecurityGroupListApi:             NewPgsqlSecurityGroupListApi(client),
		PgsqlSpecsApi:                         NewPgsqlSpecsApi(client),
		PgsqlStartApi:                         NewPgsqlStartApi(client),
		PgsqlStopApi:                          NewPgsqlStopApi(client),
		PgsqlUpdateInstanceNameApi:            NewPgsqlUpdateInstanceNameApi(client),
		PgsqlUpdateSecurityGroupApi:           NewPgsqlUpdateSecurityGroupApi(client),
		PgsqlUpgradeApi:                       NewPgsqlUpgradeApi(client),
		PgsqlBindEipApi:                       NewPgsqlBindEipApi(client),
		PgsqlUnBindEipApi:                     NewPgsqlUnBindEipApi(client),
		PgsqlBoundEipListApi:                  NewPgsqlBoundEipListApi(client),
		PgsqlDeleteSecurityGroupApi:           NewPgsqlDeleteSecurityGroupApi(client),
		PgsqlGetNodeListApi:                   NewPgsqlGetNodeListApi(client),
		PostgresqlCollectorPolicyApi:          NewPostgresqlCollectorPolicyApi(client),
		PostgresqlGetCollectorPolicyApi:       NewPostgresqlGetCollectorPolicyApi(client),
		PostgresqlAddSecurityGroupApi:         NewPostgresqlAddSecurityGroupApi(client),
		PostgresqlChangeSecurityGroupApi:      NewPostgresqlChangeSecurityGroupApi(client),
		PgsqlCreateAccountApi:                 NewPgsqlCreateAccountApi(client),
		PgsqlDeleteAccountApi:                 NewPgsqlDeleteAccountApi(client),
		PgsqlGrantPrivilegeApi:                NewPgsqlGrantPrivilegeApi(client),
		PgsqlRevokePrivilegeApi:               NewPgsqlRevokePrivilegeApi(client),
		PgsqlLockAccountApi:                   NewPgsqlLockAccountApi(client),
		PgsqlUnLockAccountApi:                 NewPgsqlUnLockAccountApi(client),
		PgsqlResetPasswordApi:                 NewPgsqlResetPasswordApi(client),
		PgsqlResetRootPasswordApi:             NewPgsqlResetRootPasswordApi(client),
		PgsqlUpdateAccountRemarkApi:           NewPgsqlUpdateAccountRemarkApi(client),
		PgsqlGetAccountListApi:                NewPgsqlGetAccountListApi(client),
		PgsqlCreateDatabaseApi:                NewPgsqlCreateDatabaseApi(client),
		PgsqlDeleteDatabaseApi:                NewPgsqlDeleteDatabaseApi(client),
		PgsqlGetCharacterSetApi:               NewPgsqlGetCharacterSetApi(client),
		PgsqlGetCollationTimeZoneApi:          NewPgsqlGetCollationTimeZoneApi(client),
		PgsqlUpdateDatabaseRemarkApi:          NewPgsqlUpdateDatabaseRemarkApi(client),
		PgsqlGetDatabaseSchemaApi:             NewPgsqlGetDatabaseSchemaApi(client),
		PgsqlGetDatabaseSchemaListApi:         NewPgsqlGetDatabaseSchemaListApi(client),
		PgsqlCreateBackupApi:                  NewPgsqlCreateBackupApi(client),
		PgsqlDeleteBackupApi:                  NewPgsqlDeleteBackupApi(client),
		PgsqlGetRecoverableBackupListApi:      NewPgsqlGetRecoverableBackupListApi(client),
		PgsqlGetBackupTaskListApi:             NewPgsqlGetBackupTaskListApi(client),
		PgsqlGetBackupListApi:                 NewPgsqlGetBackupListApi(client),
		PgsqlCreateParameterTemplateApi:       NewPgsqlCreateParameterTemplateApi(client),
		PgsqlDeleteParameterTemplateApi:       NewPgsqlDeleteParameterTemplateApi(client),
		PgsqlGetParameterTemplateDetailApi:    NewPgsqlGetParameterTemplateDetailApi(client),
		PgsqlGetParameterTemplateListApi:      NewPgsqlGetParameterTemplateListApi(client),
		PgsqlUpdateParameterTemplateApi:       NewPgsqlUpdateParameterTemplateApi(client),
		PgsqlUpdateParameterTemplateRemarkApi: NewPgsqlUpdateParameterTemplateRemarkApi(client),
		PgsqlUpdateWhiteListApi:               NewPgsqlUpdateWhiteListApi(client),
		PgsqlGetWhiteListApi:                  NewPgsqlGetWhiteListApi(client),
		PgsqlGetIDByOrderApi:                  NewPgsqlGetIDByOrderApi(client),
	}

}
