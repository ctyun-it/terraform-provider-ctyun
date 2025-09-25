package scaling

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
)

const EndpointName = "scaling"

type Apis struct {
	ScalingIsOpenApi                             *ScalingIsOpenApi
	ScalingQuotaApi                              *ScalingQuotaApi
	ScalingGroupListApi                          *ScalingGroupListApi
	ScalingGroupCreateApi                        *ScalingGroupCreateApi
	ScalingGroupDeleteApi                        *ScalingGroupDeleteApi
	ScalingGroupUpdateApi                        *ScalingGroupUpdateApi
	ScalingGroupDisableApi                       *ScalingGroupDisableApi
	ScalingGroupEnableApi                        *ScalingGroupEnableApi
	ScalingGroupUpdateInstanceMaxNumApi          *ScalingGroupUpdateInstanceMaxNumApi
	ScalingGroupUpdateInstanceMinNumApi          *ScalingGroupUpdateInstanceMinNumApi
	ScalingGroupQueryLoadBalancerListApi         *ScalingGroupQueryLoadBalancerListApi
	ScalingGroupUpdateAttachLoadBalancersApi     *ScalingGroupUpdateAttachLoadBalancersApi
	ScalingGroupUpdateDetachLoadBalancersApi     *ScalingGroupUpdateDetachLoadBalancersApi
	ScalingGroupUpdateRecoveryModeApi            *ScalingGroupUpdateRecoveryModeApi
	ScalingGroupUpdateHealthModeApi              *ScalingGroupUpdateHealthModeApi
	ScalingGroupUpdateHealthPeriodApi            *ScalingGroupUpdateHealthPeriodApi
	ScalingConfigSecurityGroupsCheckApi          *ScalingConfigSecurityGroupsCheckApi
	ScalingGroupCheckApi                         *ScalingGroupCheckApi
	ScalingUpdateConfigIdApi                     *ScalingUpdateConfigIdApi
	ScalingGroupInstanceMoveInApi                *ScalingGroupInstanceMoveInApi
	ScalingGroupInstanceMoveOutApi               *ScalingGroupInstanceMoveOutApi
	ScalingGroupInstanceMoveOutReleaseApi        *ScalingGroupInstanceMoveOutReleaseApi
	ScalingGroupSetInstancesProtectionApi        *ScalingGroupSetInstancesProtectionApi
	ScalingGroupProtectDisableApi                *ScalingGroupProtectDisableApi
	ScalingGroupProtectEnableApi                 *ScalingGroupProtectEnableApi
	ScalingGroupUpdateInstanceMoveOutStrategyApi *ScalingGroupUpdateInstanceMoveOutStrategyApi
	ScalingGroupInstanceMonitorApi               *ScalingGroupInstanceMonitorApi
	ScalingGroupUnhealthyInstanceListApi         *ScalingGroupUnhealthyInstanceListApi
	ScalingGroupQueryInstanceListApi             *ScalingGroupQueryInstanceListApi
	ScalingGroupInstanceAzApi                    *ScalingGroupInstanceAzApi
	ScalingRuleCreateApi                         *ScalingRuleCreateApi
	ScalingRuleDeleteApi                         *ScalingRuleDeleteApi
	ScalingRuleUpdateApi                         *ScalingRuleUpdateApi
	ScalingRuleExecuteApi                        *ScalingRuleExecuteApi
	ScalingRuleCreateAlarmApi                    *ScalingRuleCreateAlarmApi
	ScalingRuleDeleteAlarmApi                    *ScalingRuleDeleteAlarmApi
	ScalingRuleUpdateAlarmApi                    *ScalingRuleUpdateAlarmApi
	ScalingRuleStartAlarmApi                     *ScalingRuleStartAlarmApi
	ScalingRuleStopAlarmApi                      *ScalingRuleStopAlarmApi
	ScalingRuleCreateCycleApi                    *ScalingRuleCreateCycleApi
	ScalingRuleCreateScheduledApi                *ScalingRuleCreateScheduledApi
	ScalingRuleDeleteScheduledApi                *ScalingRuleDeleteScheduledApi
	ScalingRuleUpdateScheduledApi                *ScalingRuleUpdateScheduledApi
	ScalingRuleStartApi                          *ScalingRuleStartApi
	ScalingRuleStopApi                           *ScalingRuleStopApi
	ScalingRuleListApi                           *ScalingRuleListApi
	ScalingRuleQueryCycleApi                     *ScalingRuleQueryCycleApi
	ScalingRuleQueryAlarmApi                     *ScalingRuleQueryAlarmApi
	ScalingRuleUpdateCycleApi                    *ScalingRuleUpdateCycleApi
	ScalingRuleStartCycleApi                     *ScalingRuleStartCycleApi
	ScalingRuleDeleteCycleApi                    *ScalingRuleDeleteCycleApi
	ScalingRuleStopCycleApi                      *ScalingRuleStopCycleApi
	ScalingConfigListApi                         *ScalingConfigListApi
	ScalingConfigCreateApi                       *ScalingConfigCreateApi
	ScalingConfigDeleteApi                       *ScalingConfigDeleteApi
	ScalingConfigUpdateApi                       *ScalingConfigUpdateApi
	ScalingGroupQueryActivitiesApi               *ScalingGroupQueryActivitiesApi
	ScalingGroupQueryActivityDetailApi           *ScalingGroupQueryActivityDetailApi
	ScalingQueryActivitiesListApi                *ScalingQueryActivitiesListApi
	ScalingGroupEnableProtectionApi              *ScalingGroupEnableProtectionApi
	ScalingGroupDisableProtectionApi             *ScalingGroupDisableProtectionApi
	ScalingGroupUpdateConfiglistApi              *ScalingGroupUpdateConfiglistApi
}

func NewApis(endpointUrl string, client *core.CtyunClient) *Apis {
	client.RegisterEndpoint(core.Endpoint{
		Name: EndpointName,
		Url:  endpointUrl,
	})
	return &Apis{
		ScalingIsOpenApi:                             NewScalingIsOpenApi(client),
		ScalingQuotaApi:                              NewScalingQuotaApi(client),
		ScalingGroupListApi:                          NewScalingGroupListApi(client),
		ScalingGroupCreateApi:                        NewScalingGroupCreateApi(client),
		ScalingGroupDeleteApi:                        NewScalingGroupDeleteApi(client),
		ScalingGroupUpdateApi:                        NewScalingGroupUpdateApi(client),
		ScalingGroupDisableApi:                       NewScalingGroupDisableApi(client),
		ScalingGroupEnableApi:                        NewScalingGroupEnableApi(client),
		ScalingGroupUpdateInstanceMaxNumApi:          NewScalingGroupUpdateInstanceMaxNumApi(client),
		ScalingGroupUpdateInstanceMinNumApi:          NewScalingGroupUpdateInstanceMinNumApi(client),
		ScalingGroupQueryLoadBalancerListApi:         NewScalingGroupQueryLoadBalancerListApi(client),
		ScalingGroupUpdateAttachLoadBalancersApi:     NewScalingGroupUpdateAttachLoadBalancersApi(client),
		ScalingGroupUpdateDetachLoadBalancersApi:     NewScalingGroupUpdateDetachLoadBalancersApi(client),
		ScalingGroupUpdateRecoveryModeApi:            NewScalingGroupUpdateRecoveryModeApi(client),
		ScalingGroupUpdateHealthModeApi:              NewScalingGroupUpdateHealthModeApi(client),
		ScalingGroupUpdateHealthPeriodApi:            NewScalingGroupUpdateHealthPeriodApi(client),
		ScalingConfigSecurityGroupsCheckApi:          NewScalingConfigSecurityGroupsCheckApi(client),
		ScalingGroupCheckApi:                         NewScalingGroupCheckApi(client),
		ScalingUpdateConfigIdApi:                     NewScalingUpdateConfigIdApi(client),
		ScalingGroupInstanceMoveInApi:                NewScalingGroupInstanceMoveInApi(client),
		ScalingGroupInstanceMoveOutApi:               NewScalingGroupInstanceMoveOutApi(client),
		ScalingGroupInstanceMoveOutReleaseApi:        NewScalingGroupInstanceMoveOutReleaseApi(client),
		ScalingGroupSetInstancesProtectionApi:        NewScalingGroupSetInstancesProtectionApi(client),
		ScalingGroupProtectDisableApi:                NewScalingGroupProtectDisableApi(client),
		ScalingGroupProtectEnableApi:                 NewScalingGroupProtectEnableApi(client),
		ScalingGroupUpdateInstanceMoveOutStrategyApi: NewScalingGroupUpdateInstanceMoveOutStrategyApi(client),
		ScalingGroupInstanceMonitorApi:               NewScalingGroupInstanceMonitorApi(client),
		ScalingGroupUnhealthyInstanceListApi:         NewScalingGroupUnhealthyInstanceListApi(client),
		ScalingGroupQueryInstanceListApi:             NewScalingGroupQueryInstanceListApi(client),
		ScalingGroupInstanceAzApi:                    NewScalingGroupInstanceAzApi(client),
		ScalingRuleCreateApi:                         NewScalingRuleCreateApi(client),
		ScalingRuleDeleteApi:                         NewScalingRuleDeleteApi(client),
		ScalingRuleUpdateApi:                         NewScalingRuleUpdateApi(client),
		ScalingRuleExecuteApi:                        NewScalingRuleExecuteApi(client),
		ScalingRuleCreateAlarmApi:                    NewScalingRuleCreateAlarmApi(client),
		ScalingRuleDeleteAlarmApi:                    NewScalingRuleDeleteAlarmApi(client),
		ScalingRuleUpdateAlarmApi:                    NewScalingRuleUpdateAlarmApi(client),
		ScalingRuleStartAlarmApi:                     NewScalingRuleStartAlarmApi(client),
		ScalingRuleStopAlarmApi:                      NewScalingRuleStopAlarmApi(client),
		ScalingRuleCreateCycleApi:                    NewScalingRuleCreateCycleApi(client),
		ScalingRuleCreateScheduledApi:                NewScalingRuleCreateScheduledApi(client),
		ScalingRuleDeleteScheduledApi:                NewScalingRuleDeleteScheduledApi(client),
		ScalingRuleUpdateScheduledApi:                NewScalingRuleUpdateScheduledApi(client),
		ScalingRuleStartApi:                          NewScalingRuleStartApi(client),
		ScalingRuleStopApi:                           NewScalingRuleStopApi(client),
		ScalingRuleListApi:                           NewScalingRuleListApi(client),
		ScalingRuleQueryCycleApi:                     NewScalingRuleQueryCycleApi(client),
		ScalingRuleQueryAlarmApi:                     NewScalingRuleQueryAlarmApi(client),
		ScalingRuleUpdateCycleApi:                    NewScalingRuleUpdateCycleApi(client),
		ScalingRuleStartCycleApi:                     NewScalingRuleStartCycleApi(client),
		ScalingRuleDeleteCycleApi:                    NewScalingRuleDeleteCycleApi(client),
		ScalingRuleStopCycleApi:                      NewScalingRuleStopCycleApi(client),
		ScalingConfigListApi:                         NewScalingConfigListApi(client),
		ScalingConfigCreateApi:                       NewScalingConfigCreateApi(client),
		ScalingConfigDeleteApi:                       NewScalingConfigDeleteApi(client),
		ScalingConfigUpdateApi:                       NewScalingConfigUpdateApi(client),
		ScalingGroupQueryActivitiesApi:               NewScalingGroupQueryActivitiesApi(client),
		ScalingGroupQueryActivityDetailApi:           NewScalingGroupQueryActivityDetailApi(client),
		ScalingQueryActivitiesListApi:                NewScalingQueryActivitiesListApi(client),
		ScalingGroupEnableProtectionApi:              NewScalingGroupEnableProtectionApi(client),
		ScalingGroupDisableProtectionApi:             NewScalingGroupDisableProtectionApi(client),
		ScalingGroupUpdateConfiglistApi:              NewScalingGroupUpdateConfiglistApi(client),
	}
}
