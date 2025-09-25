package ctebs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "ebs"

type Apis struct {
	EbsUpdateIopsEbsApi             *EbsUpdateIopsEbsApi
	EbsSetDeletePolicyEbsApi        *EbsSetDeletePolicyEbsApi
	EbsCreateOrderEbsSnapApi        *EbsCreateOrderEbsSnapApi
	EbsModifyPolicyStatusEbsSnapApi *EbsModifyPolicyStatusEbsSnapApi
	EbsQueryPolicyEbsSnapApi        *EbsQueryPolicyEbsSnapApi
	EbsDeletePolicyEbsSnapApi       *EbsDeletePolicyEbsSnapApi
	EbsCancelPolicyEbsSnapApi       *EbsCancelPolicyEbsSnapApi
	EbsApplyPolicyEbsSnapApi        *EbsApplyPolicyEbsSnapApi
	EbsModifyPolicyEbsSnapApi       *EbsModifyPolicyEbsSnapApi
	EbsCreatePolicyEbsSnapApi       *EbsCreatePolicyEbsSnapApi
	EbsNewFromSnapshotEbsSnapApi    *EbsNewFromSnapshotEbsSnapApi
	EbsBatchRollbackEbsSnapApi      *EbsBatchRollbackEbsSnapApi
	EbsRollbackEbsSnapApi           *EbsRollbackEbsSnapApi
	EbsQuerySizeEbsSnapApi          *EbsQuerySizeEbsSnapApi
	EbsListEbsSnapApi               *EbsListEbsSnapApi
	EbsDeleteEbsSnapApi             *EbsDeleteEbsSnapApi
	EbsCreateEbsSnapApi             *EbsCreateEbsSnapApi
	EbsDetachEbsApi                 *EbsDetachEbsApi
	EbsAttachEbsApi                 *EbsAttachEbsApi
	EbsUpdateEbsNameApi             *EbsUpdateEbsNameApi
	EbsQueryEbsListApi              *EbsQueryEbsListApi
	EbsQueryEbsByNameApi            *EbsQueryEbsByNameApi
	EbsQueryEbsByIDApi              *EbsQueryEbsByIDApi
	EbsRenewEbsApi                  *EbsRenewEbsApi
	EbsResizeEbsApi                 *EbsResizeEbsApi
	EbsRefundEbsApi                 *EbsRefundEbsApi
	EbsNewEbsApi                    *EbsNewEbsApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		EbsUpdateIopsEbsApi:             NewEbsUpdateIopsEbsApi(client),
		EbsSetDeletePolicyEbsApi:        NewEbsSetDeletePolicyEbsApi(client),
		EbsCreateOrderEbsSnapApi:        NewEbsCreateOrderEbsSnapApi(client),
		EbsModifyPolicyStatusEbsSnapApi: NewEbsModifyPolicyStatusEbsSnapApi(client),
		EbsQueryPolicyEbsSnapApi:        NewEbsQueryPolicyEbsSnapApi(client),
		EbsDeletePolicyEbsSnapApi:       NewEbsDeletePolicyEbsSnapApi(client),
		EbsCancelPolicyEbsSnapApi:       NewEbsCancelPolicyEbsSnapApi(client),
		EbsApplyPolicyEbsSnapApi:        NewEbsApplyPolicyEbsSnapApi(client),
		EbsModifyPolicyEbsSnapApi:       NewEbsModifyPolicyEbsSnapApi(client),
		EbsCreatePolicyEbsSnapApi:       NewEbsCreatePolicyEbsSnapApi(client),
		EbsNewFromSnapshotEbsSnapApi:    NewEbsNewFromSnapshotEbsSnapApi(client),
		EbsBatchRollbackEbsSnapApi:      NewEbsBatchRollbackEbsSnapApi(client),
		EbsRollbackEbsSnapApi:           NewEbsRollbackEbsSnapApi(client),
		EbsQuerySizeEbsSnapApi:          NewEbsQuerySizeEbsSnapApi(client),
		EbsListEbsSnapApi:               NewEbsListEbsSnapApi(client),
		EbsDeleteEbsSnapApi:             NewEbsDeleteEbsSnapApi(client),
		EbsCreateEbsSnapApi:             NewEbsCreateEbsSnapApi(client),
		EbsDetachEbsApi:                 NewEbsDetachEbsApi(client),
		EbsAttachEbsApi:                 NewEbsAttachEbsApi(client),
		EbsUpdateEbsNameApi:             NewEbsUpdateEbsNameApi(client),
		EbsQueryEbsListApi:              NewEbsQueryEbsListApi(client),
		EbsQueryEbsByNameApi:            NewEbsQueryEbsByNameApi(client),
		EbsQueryEbsByIDApi:              NewEbsQueryEbsByIDApi(client),
		EbsRenewEbsApi:                  NewEbsRenewEbsApi(client),
		EbsResizeEbsApi:                 NewEbsResizeEbsApi(client),
		EbsRefundEbsApi:                 NewEbsRefundEbsApi(client),
		EbsNewEbsApi:                    NewEbsNewEbsApi(client),
	}
}
