package mysql

import (
	ctyunsdk "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

type Apis struct {
	TeledbCreateApi                        *TeledbCreateApi
	TeledbUpgradeApi                       *TeledbUpgradeApi
	TeledbRefundApi                        *TeledbRefundApi
	TeledbDestroyApi                       *TeledbDestroyApi
	TeledbQueryDetailApi                   *TeledbQueryDetailApi
	TeledbGetListApi                       *TeledbGetListApi
	TeledbStartApi                         *TeledbStartApi
	TeledbStopApi                          *TeledbStopApi
	TeledbRestartApi                       *TeledbRestartApi
	TeledbUpdateInstanceNameApi            *TeledbUpdateInstanceNameApi
	TeledbUpdateWritePortApi               *TeledbUpdateWritePortApi
	TeledbBindEipApi                       *TeledbBindEipApi
	TeledbUnbindEipApi                     *TeledbUnbindEipApi
	TeledbBoundEipListApi                  *TeledbBoundEipListApi
	TeledbMysqlSpecsApi                    *TeledbMysqlSpecsApi
	TeledbCreateAccessWhiteList            *TeledbCreateAccessWhiteList
	TeledbUpdateAccessWhiteList            *TeledbUpdateAccessWhiteList
	TeledbDeleteAccessWhiteList            *TeledbDeleteAccessWhiteList
	TeledbGetAccessWhiteList               *TeledbGetAccessWhiteList
	TeledbGetAvailabilityZone              *TeledbGetAvailabilityZone
	TeledbCreateAccountApi                 *TeledbCreateAccountApi
	TeledbDeleteAccountApi                 *TeledbDeleteAccountApi
	TeledbGetAccountInfoApi                *TeledbGetAccountInfoApi
	TeledbGetPrivilegeDetailApi            *TeledbGetPrivilegeDetailApi
	TeledbGrantPrivilegeApi                *TeledbGrantPrivilegeApi
	TeledbRevokePrivilegeApi               *TeledbRevokePrivilegeApi
	TeledbResetPasswordApi                 *TeledbResetPasswordApi
	TeledbUpdateAccountRemarkApi           *TeledbUpdateAccountRemarkApi
	TeledbCheckAccountAvailableApi         *TeledbCheckAccountAvailableApi
	TeledbRevokeSchemaApi                  *TeledbRevokeSchemaApi
	TeledbCreateBackupApi                  *TeledbCreateBackupApi
	TeledbDeleteBackupApi                  *TeledbDeleteBackupApi
	TeledbCancelBackupApi                  *TeledbCancelBackupApi
	TeledbGetBackupListApi                 *TeledbGetBackupListApi
	TeledbGetBackupRecordDetailApi         *TeledbGetBackupRecordDetailApi
	TeledbCreateRecoveryJobApi             *TeledbCreateRecoveryJobApi
	TeledbGetBackupRecoveryListApi         *TeledbGetBackupRecoveryListApi
	TeledbGetRecoverableTimeRangesApi      *TeledbGetRecoverableTimeRangesApi
	TeledbGetDatabaseSchemaApi             *TeledbGetDatabaseSchemaApi
	TeledbCreateDatabaseApi                *TeledbCreateDatabaseApi
	TeledbDeleteDatabaseApi                *TeledbDeleteDatabaseApi
	TeledbGetCharacterSetApi               *TeledbGetCharacterSetApi
	TeledbCheckDatabaseNameAvailableApi    *TeledbCheckDatabaseNameAvailableApi
	TeledbUpdateDatabaseRemarkApi          *TeledbUpdateDatabaseRemarkApi
	TeledbUpdateBackupSettingApi           *TeledbUpdateBackupSettingApi
	TeledbGetBackupSettingDetailApi        *TeledbGetBackupSettingDetailApi
	TeledbCreateParameterTemplateApi       *TeledbCreateParameterTemplateApi
	TeledbDeleteParameterTemplateApi       *TeledbDeleteParameterTemplateApi
	TeledbGetParameterTemplateDetailApi    *TeledbGetParameterTemplateDetailApi
	TeledbGetParameterTemplateListApi      *TeledbGetParameterTemplateListApi
	TeledbCopyParameterTemplateApi         *TeledbCopyParameterTemplateApi
	TeledbUpdateParameterTemplateApi       *TeledbUpdateParameterTemplateApi
	TeledbResetRootPasswordApi             *TeledbResetRootPasswordApi
	TeledbStartAuditApi                    *TeledbStartAuditApi
	TeledbGetAuditStatusApi                *TeledbGetAuditStatusApi
	TeledbUpdateRdsTemplateParameterApi    *TeledbUpdateRdsTemplateParameterApi
	TeledbGetRdsParameterTemplateDetailApi *TeledbGetRdsParameterTemplateDetailApi
	TeledbGetIDByOrderApi                  *TeledbGetIDByOrderApi
	TeledbGetStatusByIDApi                 *TeledbGetStatusByIDApi
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
		TeledbCreateApi:                        NewTeledbCreateApi(client),
		TeledbUpgradeApi:                       NewTeledbUpgradeApi(client),
		TeledbRefundApi:                        NewTeledbRefundApi(client),
		TeledbDestroyApi:                       NewTeledbDestroyApi(client),
		TeledbQueryDetailApi:                   NewTeledbQueryDetailApi(client),
		TeledbGetListApi:                       NewTeledbGetListApi(client),
		TeledbStartApi:                         NewTeledbStartApi(client),
		TeledbStopApi:                          NewTeledbStopApi(client),
		TeledbRestartApi:                       NewTeledbRestartApi(client),
		TeledbUpdateInstanceNameApi:            NewTeledbUpdateInstanceNameApi(client),
		TeledbUpdateWritePortApi:               NewTeledbUpdateWritePortApi(client),
		TeledbBindEipApi:                       NewTeledbBindEipApi(client),
		TeledbUnbindEipApi:                     NewTeledbUnbindEipApi(client),
		TeledbBoundEipListApi:                  NewTeledbBoundEipListApi(client),
		TeledbMysqlSpecsApi:                    NewTeledbMysqlSpecsApi(client),
		TeledbCreateAccessWhiteList:            NewTeledbCreateAccessWhiteList(client),
		TeledbDeleteAccessWhiteList:            NewTeledbDeleteAccessWhiteList(client),
		TeledbGetAccessWhiteList:               NewTeledbGetAccessWhiteList(client),
		TeledbUpdateAccessWhiteList:            NewTeledbUpdateAccessWhiteList(client),
		TeledbGetAvailabilityZone:              NewTeledbGetAvailabilityZone(client),
		TeledbCreateAccountApi:                 NewTeledbCreateAccountApi(client),
		TeledbDeleteAccountApi:                 NewTeledbDeleteAccountApi(client),
		TeledbGetAccountInfoApi:                NewTeledbGetAccountInfoApi(client),
		TeledbGetPrivilegeDetailApi:            NewTeledbGetPrivilegeDetailApi(client),
		TeledbGrantPrivilegeApi:                NewTeledbGrantPrivilegeApi(client),
		TeledbRevokePrivilegeApi:               NewTeledbRevokePrivilegeApi(client),
		TeledbResetPasswordApi:                 NewTeledbResetPasswordApi(client),
		TeledbUpdateAccountRemarkApi:           NewTeledbUpdateAccountRemarkApi(client),
		TeledbCheckAccountAvailableApi:         NewTeledbCheckAccountAvailableApi(client),
		TeledbCreateBackupApi:                  NewTeledbCreateBackupApi(client),
		TeledbDeleteBackupApi:                  NewTeledbDeleteBackupApi(client),
		TeledbCancelBackupApi:                  NewTeledbCancelBackupApi(client),
		TeledbGetBackupListApi:                 NewTeledbGetBackupListApi(client),
		TeledbGetBackupRecordDetailApi:         NewTeledbGetBackupRecordDetailApi(client),
		TeledbCreateRecoveryJobApi:             NewTeledbCreateRecoveryJobApi(client),
		TeledbGetBackupRecoveryListApi:         NewTeledbGetBackupRecoveryListApi(client),
		TeledbGetRecoverableTimeRangesApi:      NewTeledbGetRecoverableTimeRangesApi(client),
		TeledbRevokeSchemaApi:                  NewTeledbRevokeSchemaApi(client),
		TeledbGetDatabaseSchemaApi:             NewTeledbGetDatabaseSchemaApi(client),
		TeledbCreateDatabaseApi:                NewTeledbCreateDatabaseApi(client),
		TeledbDeleteDatabaseApi:                NewTeledbDeleteDatabaseApi(client),
		TeledbGetCharacterSetApi:               NewTeledbGetCharacterSetApi(client),
		TeledbCheckDatabaseNameAvailableApi:    NewTeledbCheckDatabaseNameAvailableApi(client),
		TeledbUpdateDatabaseRemarkApi:          NewTeledbUpdateDatabaseRemarkApi(client),
		TeledbUpdateBackupSettingApi:           NewTeledbUpdateBackupSettingApi(client),
		TeledbGetBackupSettingDetailApi:        NewTeledbGetBackupSettingDetailApi(client),
		TeledbCreateParameterTemplateApi:       NewTeledbCreateParameterTemplateApi(client),
		TeledbDeleteParameterTemplateApi:       NewTeledbDeleteParameterTemplateApi(client),
		TeledbGetParameterTemplateDetailApi:    NewTeledbGetParameterTemplateDetailApi(client),
		TeledbGetParameterTemplateListApi:      NewTeledbGetParameterTemplateListApi(client),
		TeledbCopyParameterTemplateApi:         NewTeledbCopyParameterTemplateApi(client),
		TeledbUpdateParameterTemplateApi:       NewTeledbUpdateParameterTemplateApi(client),
		TeledbResetRootPasswordApi:             NewTeledbResetRootPasswordApi(client),
		TeledbStartAuditApi:                    NewTeledbStartAuditApi(client),
		TeledbGetAuditStatusApi:                NewTeledbGetAuditStatusApi(client),
		TeledbUpdateRdsTemplateParameterApi:    NewTeledbUpdateRdsTemplateParameterApi(client),
		TeledbGetRdsParameterTemplateDetailApi: NewTeledbGetRdsParameterTemplateDetailApi(client),
		TeledbGetIDByOrderApi:                  NewTeledbGetIDByOrderApi(client),
		TeledbGetStatusByIDApi:                 NewTeledbGetStatusByIDApi(client),
	}
}
