package ctecs

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
)

// Apis api的接口
type Apis struct {
	RegionListApi                     *RegionListApi
	KeypairAttachApi                  *KeypairAttachApi
	KeypairDetachApi                  *KeypairDetachApi
	KeypairCreateApi                  *KeypairCreateApi
	KeypairDeleteApi                  *KeypairDeleteApi
	KeypairDetailApi                  *KeypairDetailApi
	KeypairImportApi                  *KeypairImportApi
	EcsDescribeInstancesApi           *EcsDescribeInstancesApi
	EcsFlavorListApi                  *EcsFlavorListApi
	EcsCreateInstanceApi              *EcsCreateInstanceApi
	EcsJoinSecurityGroupApi           *EcsJoinSecurityGroupApi
	EcsLeaveSecurityGroupApi          *EcsLeaveSecurityGroupApi
	EcsVolumeListApi                  *EcsVolumeListApi
	EcsInstanceDetailsApi             *EcsInstanceDetailsApi
	EcsUnsubscribeInstanceApi         *EcsUnsubscribeInstanceApi
	EcsUpdateFlavorSpecApi            *EcsUpdateFlavorSpecApi
	EcsQueryAsyncResultApi            *EcsQueryAsyncResultApi
	EcsStartInstanceApi               *EcsStartInstanceApi
	EcsStopInstanceApi                *EcsStopInstanceApi
	EcsInstanceStatusListApi          *EcsInstanceStatusListApi
	EcsResetPasswordApi               *EcsResetPasswordApi
	EcsChangeToCycleApi               *EcsChangeToCycleApi
	EcsTagOnDemandApi                 *EcsTagOnDemandApi
	EcsTerminateCycleApi              *EcsTerminateCycleApi
	EcsBatchUpdateInstancesApi        *EcsBatchUpdateInstancesApi
	JobShowApi                        *JobShowApi
	EcsOrderQueryUuidApi              *EcsOrderQueryUuidApi
	SecurityGroupRuleEgressCreateApi  *SecurityGroupRuleEgressCreateApi
	SecurityGroupRuleIngressCreateApi *SecurityGroupRuleIngressCreateApi
	EcsShelveInstanceApi              *EcsShelveInstanceApi
	EcsGetUserdataMonitorApi          *EcsGetUserdataMonitorApi
	EcsUserdataDetailsApi             *EcsUserdataDetailsApi
	EcsDestroyInstanceApi             *EcsDestroyInstanceApi
}

// NewApis 构建
func NewApis(client *ctyunsdk.CtyunClient) *Apis {
	client.RegisterEndpoint(ctyunsdk.EnvironmentDev, EndpointCtecsTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentTest, EndpointCtecsTest)
	client.RegisterEndpoint(ctyunsdk.EnvironmentProd, EndpointCtecsProd)
	return &Apis{
		RegionListApi:                     NewRegionListApi(client),
		KeypairAttachApi:                  NewKeypairAttachApi(client),
		KeypairDetachApi:                  NewKeypairDetachApi(client),
		KeypairCreateApi:                  NewKeypairCreateApi(client),
		KeypairDeleteApi:                  NewKeypairDeleteApi(client),
		KeypairDetailApi:                  NewKeypairDetailApi(client),
		KeypairImportApi:                  NewKeypairImportApi(client),
		EcsDescribeInstancesApi:           NewEcsDescribeInstancesApi(client),
		EcsFlavorListApi:                  NewEcsFlavorListApi(client),
		EcsCreateInstanceApi:              NewEcsCreateInstanceApi(client),
		EcsJoinSecurityGroupApi:           NewEcsJoinSecurityGroupApi(client),
		EcsLeaveSecurityGroupApi:          NewEcsLeaveSecurityGroupApi(client),
		EcsVolumeListApi:                  NewEcsVolumeListApi(client),
		EcsInstanceDetailsApi:             NewEcsInstanceDetailsApi(client),
		EcsUnsubscribeInstanceApi:         NewEcsUnsubscribeInstanceApi(client),
		EcsUpdateFlavorSpecApi:            NewEcsUpdateFlavorSpecApi(client),
		EcsQueryAsyncResultApi:            NewEcsQueryAsyncResultApi(client),
		EcsStartInstanceApi:               NewEcsStartInstanceApi(client),
		EcsStopInstanceApi:                NewEcsStopInstanceApi(client),
		EcsInstanceStatusListApi:          NewEcsInstanceStatusListApi(client),
		EcsResetPasswordApi:               NewEcsResetPasswordApi(client),
		EcsChangeToCycleApi:               NewEcsChangeToCycleApi(client),
		EcsTagOnDemandApi:                 NewEcsTagOnDemandApi(client),
		EcsTerminateCycleApi:              NewEcsTerminateCycleApi(client),
		EcsBatchUpdateInstancesApi:        NewEcsBatchUpdateInstancesApi(client),
		JobShowApi:                        NewJobShowApi(client),
		EcsOrderQueryUuidApi:              NewEcsOrderQueryUuid(client),
		SecurityGroupRuleEgressCreateApi:  NewSecurityGroupRuleEgressCreateApi(client),
		SecurityGroupRuleIngressCreateApi: NewSecurityGroupRuleIngressCreateApi(client),
		EcsShelveInstanceApi:              NewEcsShelveInstanceApi(client),
		EcsGetUserdataMonitorApi:          NewEcsGetUserdataMonitorApi(client),
		EcsUserdataDetailsApi:             NewEcsUserdataDetailsApi(client),
		EcsDestroyInstanceApi:             NewEcsDestroyInstanceApi(client),
	}
}
