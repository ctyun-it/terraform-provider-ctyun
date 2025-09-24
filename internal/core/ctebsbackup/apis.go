package ctebsbackup

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ebsbackup"

type Apis struct {
	EbsbackupListEbsBackupPolicyDisksApi     *EbsbackupListEbsBackupPolicyDisksApi
	EbsbackupShowEbsBackupPolicyTaskApi      *EbsbackupShowEbsBackupPolicyTaskApi
	EbsbackupEbsBackupPolicyBindRepoApi      *EbsbackupEbsBackupPolicyBindRepoApi
	EbsbackupUpdateEbsBackupPolicyApi        *EbsbackupUpdateEbsBackupPolicyApi
	EbsbackupEnableEbsBackupPolicyApi        *EbsbackupEnableEbsBackupPolicyApi
	EbsbackupDeleteEbsBackupPolicyApi        *EbsbackupDeleteEbsBackupPolicyApi
	EbsbackupCreateEbsBackupPolicyApi        *EbsbackupCreateEbsBackupPolicyApi
	EbsbackupEbsBackupPolicyUnbindRepoApi    *EbsbackupEbsBackupPolicyUnbindRepoApi
	EbsbackupExecuteEbsBackupPolicyApi       *EbsbackupExecuteEbsBackupPolicyApi
	EbsbackupCreateRepoApi                   *EbsbackupCreateRepoApi
	EbsbackupUpdateEbsBackupRepoApi          *EbsbackupUpdateEbsBackupRepoApi
	EbsbackupResizeRepoApi                   *EbsbackupResizeRepoApi
	EbsbackupRenewRepoApi                    *EbsbackupRenewRepoApi
	EbsbackupDeleteRepoApi                   *EbsbackupDeleteRepoApi
	EbsbackupDeleteEbsBackupApi              *EbsbackupDeleteEbsBackupApi
	EbsbackupDisableEbsBackupPolicyApi       *EbsbackupDisableEbsBackupPolicyApi
	EbsbackupListEbsBackupPolicyTasksApi     *EbsbackupListEbsBackupPolicyTasksApi
	EbsbackupCreateEbsBackupApi              *EbsbackupCreateEbsBackupApi
	EbsbackupRestoreEbsBackupApi             *EbsbackupRestoreEbsBackupApi
	EbsbackupListEbsBackupApi                *EbsbackupListEbsBackupApi
	EbsbackupShowEbsBackupApi                *EbsbackupShowEbsBackupApi
	EbsbackupShowEbsBackupUsageApi           *EbsbackupShowEbsBackupUsageApi
	EbsbackupListEbsBackupPolicyApi          *EbsbackupListEbsBackupPolicyApi
	EbsbackupEbsBackupPolicyBindVolumesApi   *EbsbackupEbsBackupPolicyBindVolumesApi
	EbsbackupEbsBackupPolicyUnbindVolumesApi *EbsbackupEbsBackupPolicyUnbindVolumesApi
	EbsbackupListEbsBackupRepoApi            *EbsbackupListEbsBackupRepoApi
	EbsbackupListBackupRepoApi               *EbsbackupListBackupRepoApi
	EbsbackupShowBackupApi                   *EbsbackupShowBackupApi
	EbsbackupListBackupPolicyApi             *EbsbackupListBackupPolicyApi
	EbsbackupListBackupApi                   *EbsbackupListBackupApi
	EbsbackupShowBackupUsageApi              *EbsbackupShowBackupUsageApi
	EbsbackupCreateBackupApi                 *EbsbackupCreateBackupApi
	EbsbackupRestoreBackupApi                *EbsbackupRestoreBackupApi
	EbsbackupEbsBackupPolicyBindDisksApi     *EbsbackupEbsBackupPolicyBindDisksApi
	EbsbackupEbsBackupPolicyUnbindDisksApi   *EbsbackupEbsBackupPolicyUnbindDisksApi
	EbsbackupListBackupTaskApi               *EbsbackupListBackupTaskApi
	EbsbackupCancelBackupTaskApi             *EbsbackupCancelBackupTaskApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		EbsbackupListEbsBackupPolicyDisksApi:     NewEbsbackupListEbsBackupPolicyDisksApi(client),
		EbsbackupShowEbsBackupPolicyTaskApi:      NewEbsbackupShowEbsBackupPolicyTaskApi(client),
		EbsbackupEbsBackupPolicyBindRepoApi:      NewEbsbackupEbsBackupPolicyBindRepoApi(client),
		EbsbackupUpdateEbsBackupPolicyApi:        NewEbsbackupUpdateEbsBackupPolicyApi(client),
		EbsbackupEnableEbsBackupPolicyApi:        NewEbsbackupEnableEbsBackupPolicyApi(client),
		EbsbackupDeleteEbsBackupPolicyApi:        NewEbsbackupDeleteEbsBackupPolicyApi(client),
		EbsbackupCreateEbsBackupPolicyApi:        NewEbsbackupCreateEbsBackupPolicyApi(client),
		EbsbackupEbsBackupPolicyUnbindRepoApi:    NewEbsbackupEbsBackupPolicyUnbindRepoApi(client),
		EbsbackupExecuteEbsBackupPolicyApi:       NewEbsbackupExecuteEbsBackupPolicyApi(client),
		EbsbackupCreateRepoApi:                   NewEbsbackupCreateRepoApi(client),
		EbsbackupUpdateEbsBackupRepoApi:          NewEbsbackupUpdateEbsBackupRepoApi(client),
		EbsbackupResizeRepoApi:                   NewEbsbackupResizeRepoApi(client),
		EbsbackupRenewRepoApi:                    NewEbsbackupRenewRepoApi(client),
		EbsbackupDeleteRepoApi:                   NewEbsbackupDeleteRepoApi(client),
		EbsbackupDeleteEbsBackupApi:              NewEbsbackupDeleteEbsBackupApi(client),
		EbsbackupDisableEbsBackupPolicyApi:       NewEbsbackupDisableEbsBackupPolicyApi(client),
		EbsbackupListEbsBackupPolicyTasksApi:     NewEbsbackupListEbsBackupPolicyTasksApi(client),
		EbsbackupCreateEbsBackupApi:              NewEbsbackupCreateEbsBackupApi(client),
		EbsbackupRestoreEbsBackupApi:             NewEbsbackupRestoreEbsBackupApi(client),
		EbsbackupListEbsBackupApi:                NewEbsbackupListEbsBackupApi(client),
		EbsbackupShowEbsBackupApi:                NewEbsbackupShowEbsBackupApi(client),
		EbsbackupShowEbsBackupUsageApi:           NewEbsbackupShowEbsBackupUsageApi(client),
		EbsbackupListEbsBackupPolicyApi:          NewEbsbackupListEbsBackupPolicyApi(client),
		EbsbackupEbsBackupPolicyBindVolumesApi:   NewEbsbackupEbsBackupPolicyBindVolumesApi(client),
		EbsbackupEbsBackupPolicyUnbindVolumesApi: NewEbsbackupEbsBackupPolicyUnbindVolumesApi(client),
		EbsbackupListEbsBackupRepoApi:            NewEbsbackupListEbsBackupRepoApi(client),
		EbsbackupListBackupRepoApi:               NewEbsbackupListBackupRepoApi(client),
		EbsbackupShowBackupApi:                   NewEbsbackupShowBackupApi(client),
		EbsbackupListBackupPolicyApi:             NewEbsbackupListBackupPolicyApi(client),
		EbsbackupListBackupApi:                   NewEbsbackupListBackupApi(client),
		EbsbackupShowBackupUsageApi:              NewEbsbackupShowBackupUsageApi(client),
		EbsbackupCreateBackupApi:                 NewEbsbackupCreateBackupApi(client),
		EbsbackupRestoreBackupApi:                NewEbsbackupRestoreBackupApi(client),
		EbsbackupEbsBackupPolicyBindDisksApi:     NewEbsbackupEbsBackupPolicyBindDisksApi(client),
		EbsbackupEbsBackupPolicyUnbindDisksApi:   NewEbsbackupEbsBackupPolicyUnbindDisksApi(client),
		EbsbackupListBackupTaskApi:               NewEbsbackupListBackupTaskApi(client),
		EbsbackupCancelBackupTaskApi:             NewEbsbackupCancelBackupTaskApi(client),
	}
}
