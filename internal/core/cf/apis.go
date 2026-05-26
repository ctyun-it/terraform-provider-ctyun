package cf

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "cf"

type Apis struct {
	CfDeleteFunctionV2024Api     *CfDeleteFunctionV2024Api
	CfGetFunctionV2024Api        *CfGetFunctionV2024Api
	CfListFunctionsV2024Api      *CfListFunctionsV2024Api
	CfGetFunctionCodeApi         *CfGetFunctionCodeApi
	CfUpdateFunctionApi          *CfUpdateFunctionApi
	CfListVpcBindingsApi         *CfListVpcBindingsApi
	CfListInstancesApi           *CfListInstancesApi
	CfListFunctionVersionsApi    *CfListFunctionVersionsApi
	CfPublishFunctionVersionApi  *CfPublishFunctionVersionApi
	CfDeleteFunctionVersionApi   *CfDeleteFunctionVersionApi
	CfDeleteAliasApi             *CfDeleteAliasApi
	CfGetAliasApi                *CfGetAliasApi
	CfUpdateAliasApi             *CfUpdateAliasApi
	CfCreateAliasApi             *CfCreateAliasApi
	CfCreateTriggerV2024Api      *CfCreateTriggerV2024Api
	CfDeleteTriggerV2024Api      *CfDeleteTriggerV2024Api
	CfGetTriggerApi              *CfGetTriggerApi
	CfListTriggersV2024Api       *CfListTriggersV2024Api
	CfUpdateTriggerV2024Api      *CfUpdateTriggerV2024Api
	CfDeleteProvisionConfigApi   *CfDeleteProvisionConfigApi
	CfGetProvisionConfigApi      *CfGetProvisionConfigApi
	CfListProvisionConfigsApi    *CfListProvisionConfigsApi
	CfPutProvisionConfigApi      *CfPutProvisionConfigApi
	CfDeleteConcurrencyConfigApi *CfDeleteConcurrencyConfigApi
	CfGetConcurrencyConfigApi    *CfGetConcurrencyConfigApi
	CfListConcurrencyConfigsApi  *CfListConcurrencyConfigsApi
	CfPutConcurrencyConfigApi    *CfPutConcurrencyConfigApi
	CfCreateLayerVersionApi      *CfCreateLayerVersionApi
	CfDeleteLayerVersionApi      *CfDeleteLayerVersionApi
	CfGetLayerVersionApi         *CfGetLayerVersionApi
	CfGetLayerVersionByCtrnApi   *CfGetLayerVersionByCtrnApi
	CfListLayersApi              *CfListLayersApi
	CfPutLayerACLApi             *CfPutLayerACLApi
	CfDeleteAsyncInvokeConfigApi *CfDeleteAsyncInvokeConfigApi
	CfGetAsyncInvokeConfigApi    *CfGetAsyncInvokeConfigApi
	CfListAsyncInvokeConfigsApi  *CfListAsyncInvokeConfigsApi
	CfPutAsyncInvokeConfigApi    *CfPutAsyncInvokeConfigApi
	CfGetAsyncTaskApi            *CfGetAsyncTaskApi
	CfListAsyncTasksApi          *CfListAsyncTasksApi
	CfStopAsyncTaskApi           *CfStopAsyncTaskApi
	CfInvokeFunctionApi          *CfInvokeFunctionApi
	CfCreateCustomDomainApi      *CfCreateCustomDomainApi
	CfDeleteCustomDomainApi      *CfDeleteCustomDomainApi
	CfGetCustomDomainApi         *CfGetCustomDomainApi
	CfListCustomDomainsApi       *CfListCustomDomainsApi
	CfUpdateCustomDomainApi      *CfUpdateCustomDomainApi
	CfCreateFunctionV2024Api     *CfCreateFunctionV2024Api
	CfGetWorkflowApi             *CfGetWorkflowApi
	CfCreateWorkflowApi          *CfCreateWorkflowApi
	CfUpdateWorkflowApi          *CfUpdateWorkflowApi
	CfDeleteWorkflowApi          *CfDeleteWorkflowApi
	CfQueryWorkflowApi           *CfQueryWorkflowApi
	CfExecuteWorkflowSyncApi     *CfExecuteWorkflowSyncApi
	CfExecuteWorkflowApi         *CfExecuteWorkflowApi
	CfGetExecutionHistoryApi     *CfGetExecutionHistoryApi
	CfGetExecutionDetailApi      *CfGetExecutionDetailApi
	CfGetTaskResultApi           *CfGetTaskResultApi
	CfGetTaskEventTraceApi       *CfGetTaskEventTraceApi
	CfReportTaskSuccessApi       *CfReportTaskSuccessApi
	CfReportTaskFailureApi       *CfReportTaskFailureApi
	CfStopExecutionApi           *CfStopExecutionApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		CfDeleteFunctionV2024Api:     NewCfDeleteFunctionV2024Api(client),
		CfGetFunctionV2024Api:        NewCfGetFunctionV2024Api(client),
		CfListFunctionsV2024Api:      NewCfListFunctionsV2024Api(client),
		CfGetFunctionCodeApi:         NewCfGetFunctionCodeApi(client),
		CfUpdateFunctionApi:          NewCfUpdateFunctionApi(client),
		CfListVpcBindingsApi:         NewCfListVpcBindingsApi(client),
		CfListInstancesApi:           NewCfListInstancesApi(client),
		CfListFunctionVersionsApi:    NewCfListFunctionVersionsApi(client),
		CfPublishFunctionVersionApi:  NewCfPublishFunctionVersionApi(client),
		CfDeleteFunctionVersionApi:   NewCfDeleteFunctionVersionApi(client),
		CfDeleteAliasApi:             NewCfDeleteAliasApi(client),
		CfGetAliasApi:                NewCfGetAliasApi(client),
		CfUpdateAliasApi:             NewCfUpdateAliasApi(client),
		CfCreateAliasApi:             NewCfCreateAliasApi(client),
		CfCreateTriggerV2024Api:      NewCfCreateTriggerV2024Api(client),
		CfDeleteTriggerV2024Api:      NewCfDeleteTriggerV2024Api(client),
		CfGetTriggerApi:              NewCfGetTriggerApi(client),
		CfListTriggersV2024Api:       NewCfListTriggersV2024Api(client),
		CfUpdateTriggerV2024Api:      NewCfUpdateTriggerV2024Api(client),
		CfDeleteProvisionConfigApi:   NewCfDeleteProvisionConfigApi(client),
		CfGetProvisionConfigApi:      NewCfGetProvisionConfigApi(client),
		CfListProvisionConfigsApi:    NewCfListProvisionConfigsApi(client),
		CfPutProvisionConfigApi:      NewCfPutProvisionConfigApi(client),
		CfDeleteConcurrencyConfigApi: NewCfDeleteConcurrencyConfigApi(client),
		CfGetConcurrencyConfigApi:    NewCfGetConcurrencyConfigApi(client),
		CfListConcurrencyConfigsApi:  NewCfListConcurrencyConfigsApi(client),
		CfPutConcurrencyConfigApi:    NewCfPutConcurrencyConfigApi(client),
		CfCreateLayerVersionApi:      NewCfCreateLayerVersionApi(client),
		CfDeleteLayerVersionApi:      NewCfDeleteLayerVersionApi(client),
		CfGetLayerVersionApi:         NewCfGetLayerVersionApi(client),
		CfGetLayerVersionByCtrnApi:   NewCfGetLayerVersionByCtrnApi(client),
		CfListLayersApi:              NewCfListLayersApi(client),
		CfPutLayerACLApi:             NewCfPutLayerACLApi(client),
		CfDeleteAsyncInvokeConfigApi: NewCfDeleteAsyncInvokeConfigApi(client),
		CfGetAsyncInvokeConfigApi:    NewCfGetAsyncInvokeConfigApi(client),
		CfListAsyncInvokeConfigsApi:  NewCfListAsyncInvokeConfigsApi(client),
		CfPutAsyncInvokeConfigApi:    NewCfPutAsyncInvokeConfigApi(client),
		CfGetAsyncTaskApi:            NewCfGetAsyncTaskApi(client),
		CfListAsyncTasksApi:          NewCfListAsyncTasksApi(client),
		CfStopAsyncTaskApi:           NewCfStopAsyncTaskApi(client),
		CfInvokeFunctionApi:          NewCfInvokeFunctionApi(client),
		CfCreateCustomDomainApi:      NewCfCreateCustomDomainApi(client),
		CfDeleteCustomDomainApi:      NewCfDeleteCustomDomainApi(client),
		CfGetCustomDomainApi:         NewCfGetCustomDomainApi(client),
		CfListCustomDomainsApi:       NewCfListCustomDomainsApi(client),
		CfUpdateCustomDomainApi:      NewCfUpdateCustomDomainApi(client),
		CfCreateFunctionV2024Api:     NewCfCreateFunctionV2024Api(client),
		CfGetWorkflowApi:             NewCfGetWorkflowApi(client),
		CfCreateWorkflowApi:          NewCfCreateWorkflowApi(client),
		CfUpdateWorkflowApi:          NewCfUpdateWorkflowApi(client),
		CfDeleteWorkflowApi:          NewCfDeleteWorkflowApi(client),
		CfQueryWorkflowApi:           NewCfQueryWorkflowApi(client),
		CfExecuteWorkflowSyncApi:     NewCfExecuteWorkflowSyncApi(client),
		CfExecuteWorkflowApi:         NewCfExecuteWorkflowApi(client),
		CfGetExecutionHistoryApi:     NewCfGetExecutionHistoryApi(client),
		CfGetExecutionDetailApi:      NewCfGetExecutionDetailApi(client),
		CfGetTaskResultApi:           NewCfGetTaskResultApi(client),
		CfGetTaskEventTraceApi:       NewCfGetTaskEventTraceApi(client),
		CfReportTaskSuccessApi:       NewCfReportTaskSuccessApi(client),
		CfReportTaskFailureApi:       NewCfReportTaskFailureApi(client),
		CfStopExecutionApi:           NewCfStopExecutionApi(client),
	}
}
