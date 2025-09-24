package business

const (
	ScalingHealthModeServer = 1 // 1：云服务器健康检查
	ScalingHealthModeElb    = 2 // 2：弹性负载均衡健康检查

	ScalingHealthModeServerStr = "server" //云服务器健康检查
	ScalingHealthModeElbStr    = "lb"     // 弹性负载均衡健康检查

	ScalingMoveOutStrategyEarlierConfig = 1 // 1：较早创建的配置较早创建的云主机
	ScalingMoveOutStrategyLaterConfig   = 2 // 2：较晚创建的配置较晚创建的云主机
	ScalingMoveOutStrategyEarlierVM     = 3 // 3：较早创建的云主机
	ScalingMoveOutStrategyLaterVM       = 4 // 4：较晚创建的云主机

	ScalingMoveOutStrategyEarlierConfigStr = "earlier_config" // 较早创建的配置较早创建的云主机
	ScalingMoveOutStrategyLaterConfigStr   = "later_config"   // 较晚创建的配置较晚创建的云主机
	ScalingMoveOutStrategyEarlierVMStr     = "earlier_vm"     // 较早创建的云主机
	ScalingMoveOutStrategyLaterVMStr       = "later_vm"       // 较晚创建的云主机

	ScalingAzStrategyUniformDistribution  = 1 // 均衡分布
	ScalingAzStrategyPriorityDistribution = 2 // 优先级分布

	ScalingAzStrategyUniformDistributionStr  = "uniform_distribution"
	ScalingAzStrategyPriorityDistributionStr = "priority_distribution"

	ScalingControlStatusStart = 1
	ScalingControlStatusStop  = 2

	ScalingControlStatusStartStr = "enable"
	ScalingControlStatusStopStr  = "disable"

	ScalingControlProtectionEnable  = true
	ScalingControlProtectionDisable = false

	ScalingControlProtectionEnableStr  = "enable"
	ScalingControlProtectionDisableStr = "disable"

	ScalingLoginModePassword = 1
	ScalingLoginModeKeyPair  = 2

	ScalingLoginModePasswordStr = "password"
	ScalingLoginModeKeyPairStr  = "key_pair"

	ScalingUseFloatingsDisable = 1 // 1. 不使用
	ScalingUseFloatingsAuto    = 2 // 2. 自动分配

	ScalingUseFloatingsDisableStr = "diable" //  不使用
	ScalingUseFloatingsAutoStr    = "auto"   //  自动分配

	ScalingVolumeFlagOS      = 1      // 系统盘
	ScalingVolumeFlagData    = 2      // 数据盘
	ScalingVolumeFlagOSStr   = "OS"   // 系统盘
	ScalingVolumeFlagDataStr = "DATA" // 数据盘

	ScalingVisibilityPublic  = 1
	ScalingVisibilityPrivate = 2

	ScalingVisibilityPublicStr  = "public"  // 共有镜像
	ScalingVisibilityPrivateStr = "private" // 私有镜像

	ScalingOsTypeLinux      = 1
	ScalingOsTypeWindows    = 2
	ScalingOsTypeLinuxStr   = "linux"
	ScalingOsTypeWindowsStr = "windows"

	ScalingPolicyAlert   = 1
	ScalingPolicyRegular = 2
	ScalingPolicyPeriod  = 3
	ScalingPolicyTarget  = 4

	ScalingPolicyAlertStr   = "alert"
	ScalingPolicyRegularStr = "regular"
	ScalingPolicyPeriodStr  = "period"
	ScalingPolicyTargetStr  = "target"

	ScalingPolicyOperateUnitCount      = 1 // 个数
	ScalingPolicyOperateUnitPercent    = 2 // 百分比
	ScalingPolicyOperateUnitNull       = 0
	ScalingPolicyOperateUnitCountStr   = "count"
	ScalingPolicyOperateUnitPercentStr = "percent"

	ActionIncrease    = 1 // 增加
	ActionDecrease    = 2 // 减少
	ActionSet         = 3 // 设置为
	ActionNull        = 0
	ActionIncreaseStr = "increase"
	ActionDecreaseStr = "decrease"
	ActionSetStr      = "set"
	defaultStr        = ""

	CycleMonthly    = 1 // 按月循环
	CycleWeekly     = 2 // 按周循环
	CycleDaily      = 3 // 按天循环
	CycleNull       = 0
	CycleMonthlyStr = "monthly"
	CycleWeeklyStr  = "weekly"
	CycleDailyStr   = "daily"

	StatusEnabled     = 1 // 启用
	StatusDisabled    = 2 // 停用
	StatusEnabledStr  = "enable"
	StatusDisabledStr = "disable"

	ScalingActivityEnabled   = 1 // 已启用
	ScalingActivityMovingIn  = 2 // 正在移入
	ScalingActivityMovingOut = 3 // 正在移出

	ScalingActivityEnabledStr   = "enabled"
	ScalingActivityMovingInStr  = "moving_in"
	ScalingActivityMovingOutStr = "moving_out"

	ExecutionModeAutoStrategy          = 1
	ExecutionModeManualStrategy        = 2
	ExecutionModeManualAddInstances    = 3
	ExecutionModeManualRemoveInstances = 4
	ExecutionModeSatisfyMinSize        = 5
	ExecutionModeAdjustToLimits        = 6
	ExecutionModeHealthCheckAdd        = 7
	ExecutionModeHealthCheckRemove     = 8

	ExecutionModeAutoStrategyStr          = "auto_strategy"           // 自动执行策略
	ExecutionModeManualStrategyStr        = "manual_strategy"         // 手动执行策略
	ExecutionModeManualAddInstancesStr    = "manual_add_instances"    // 手动移入实例
	ExecutionModeManualRemoveInstancesStr = "manual_remove_instances" // 手动移出实例
	ExecutionModeSatisfyMinSizeStr        = "satisfy_min_size"        // 新建伸缩组满足最小数
	ExecutionModeAdjustToLimitsStr        = "adjust_to_limits"        // 修改伸缩组满足最大最小限制
	ExecutionModeHealthCheckAddStr        = "health_check_add"        // 健康检查移入
	ExecutionModeHealthCheckRemoveStr     = "health_check_remove"     // 健康检查移出

	HealthStatusNormal       = 1 // 正常
	HealthStatusAbnormal     = 2 // 异常
	HealthStatusInitializing = 3 // 初始化

	HealthStatusNormalStr       = "normal"       // 正常
	HealthStatusAbnormalStr     = "abnormal"     // 异常
	HealthStatusInitializingStr = "initializing" // 初始化

	ProtectStatusProtected   = 1 // 已保护
	ProtectStatusUnprotected = 2 // 未保护

	ProtectStatusProtectedStr   = "enable"  // 已保护
	ProtectStatusUnprotectedStr = "disable" // 未保护
)

var ScalingHealthMode = []string{
	ScalingHealthModeServerStr,
	ScalingHealthModeElbStr,
}

var ScalingHealthModeDict = map[string]int32{
	ScalingHealthModeServerStr: ScalingHealthModeServer,
	ScalingHealthModeElbStr:    ScalingHealthModeElb,
}

var ScalingHealthModeDictRev = map[int32]string{
	ScalingHealthModeServer: ScalingHealthModeServerStr,
	ScalingHealthModeElb:    ScalingHealthModeElbStr,
}

var ScalingMoveOutStrategy = []string{
	ScalingMoveOutStrategyEarlierConfigStr,
	ScalingMoveOutStrategyLaterConfigStr,
	ScalingMoveOutStrategyEarlierVMStr,
	ScalingMoveOutStrategyLaterVMStr,
}
var ScalingMoveOutStrategyDict = map[string]int32{
	ScalingMoveOutStrategyEarlierConfigStr: ScalingMoveOutStrategyEarlierConfig,
	ScalingMoveOutStrategyLaterConfigStr:   ScalingMoveOutStrategyLaterConfig,
	ScalingMoveOutStrategyEarlierVMStr:     ScalingMoveOutStrategyEarlierVM,
	ScalingMoveOutStrategyLaterVMStr:       ScalingMoveOutStrategyLaterVM,
}

var ScalingMoveOutStrategyDictRev = map[int32]string{
	ScalingMoveOutStrategyEarlierConfig: ScalingMoveOutStrategyEarlierConfigStr,
	ScalingMoveOutStrategyLaterConfig:   ScalingMoveOutStrategyLaterConfigStr,
	ScalingMoveOutStrategyEarlierVM:     ScalingMoveOutStrategyEarlierVMStr,
	ScalingMoveOutStrategyLaterVM:       ScalingMoveOutStrategyLaterVMStr,
}

var ScalingAzStrategy = []string{
	ScalingAzStrategyUniformDistributionStr,
	ScalingAzStrategyPriorityDistributionStr,
}

var ScalingAzStrategyDict = map[string]int32{
	ScalingAzStrategyUniformDistributionStr:  ScalingAzStrategyUniformDistribution,
	ScalingAzStrategyPriorityDistributionStr: ScalingAzStrategyPriorityDistribution,
}

var ScalingAzStrategyDictRev = map[int32]string{
	ScalingAzStrategyUniformDistribution:  ScalingAzStrategyUniformDistributionStr,
	ScalingAzStrategyPriorityDistribution: ScalingAzStrategyPriorityDistributionStr,
}

var ScalingControlStatus = []string{
	ScalingControlStatusStartStr,
	ScalingControlStatusStopStr,
}

var ScalingControlStatusDict = map[string]int32{
	ScalingControlStatusStartStr: ScalingControlStatusStart,
	ScalingControlStatusStopStr:  ScalingControlStatusStop,
}

var ScalingControlStatusDictRev = map[int32]string{
	ScalingControlStatusStart: ScalingControlStatusStartStr,
	ScalingControlStatusStop:  ScalingControlStatusStopStr,
}
var ScalingControlProtectionStatus = []string{
	ScalingControlProtectionEnableStr,
	ScalingControlProtectionDisableStr,
}

var ScalingControlProtectionDict = map[string]bool{
	ScalingControlProtectionEnableStr:  ScalingControlProtectionEnable,
	ScalingControlProtectionDisableStr: ScalingControlProtectionDisable,
}

var ScalingControlProtectionDictRev = map[bool]string{
	ScalingControlProtectionEnable:  ScalingControlProtectionEnableStr,
	ScalingControlProtectionDisable: ScalingControlProtectionDisableStr,
}

var ScalingLoginMode = []string{
	ScalingLoginModePasswordStr,
	ScalingLoginModeKeyPairStr,
}

var ScalingLoginModeDict = map[string]int32{
	ScalingLoginModePasswordStr: ScalingLoginModePassword,
	ScalingLoginModeKeyPairStr:  ScalingLoginModeKeyPair,
}

var ScalingLoginModeDictRev = map[int32]string{
	ScalingLoginModePassword: ScalingLoginModePasswordStr,
	ScalingLoginModeKeyPair:  ScalingLoginModeKeyPairStr,
}

var ScalingUseFloatings = []string{
	ScalingUseFloatingsDisableStr,
	ScalingUseFloatingsAutoStr,
}

var ScalingUseFloatingsDict = map[string]int32{
	ScalingUseFloatingsDisableStr: ScalingUseFloatingsDisable,
	ScalingUseFloatingsAutoStr:    ScalingUseFloatingsAuto,
}

var ScalingUseFloatingsDictRev = map[int32]string{
	ScalingUseFloatingsDisable: ScalingUseFloatingsDisableStr,
	ScalingUseFloatingsAuto:    ScalingUseFloatingsAutoStr,
}

var ScalingVolumeFlag = []string{
	ScalingVolumeFlagOSStr,
	ScalingVolumeFlagDataStr,
}
var ScalingVolumeFlagDict = map[string]int32{
	ScalingVolumeFlagOSStr:   ScalingVolumeFlagOS,
	ScalingVolumeFlagDataStr: ScalingVolumeFlagData,
}

var ScalingVolumeFlagDictRev = map[int32]string{
	ScalingVolumeFlagOS:   ScalingVolumeFlagOSStr,
	ScalingVolumeFlagData: ScalingVolumeFlagDataStr,
}

var ScalingVisibility = []string{
	ScalingVisibilityPublicStr,
	ScalingVisibilityPrivateStr,
}

var ScalingVisibilityDictRev = map[int32]string{
	ScalingVisibilityPublic:  ScalingVisibilityPublicStr,
	ScalingVisibilityPrivate: ScalingVisibilityPrivateStr,
}

var ScalingOsTypeDictRev = map[int32]string{
	ScalingOsTypeLinux:   ScalingOsTypeLinuxStr,
	ScalingOsTypeWindows: ScalingOsTypeWindowsStr,
}

// 策略类型数组和映射
var ScalingPolicyTypes = []string{
	ScalingPolicyAlertStr,
	ScalingPolicyRegularStr,
	ScalingPolicyPeriodStr,
	ScalingPolicyTargetStr,
}
var ScalingPolicyTypeDict = map[string]int32{
	ScalingPolicyAlertStr:   ScalingPolicyAlert,
	ScalingPolicyRegularStr: ScalingPolicyRegular,
	ScalingPolicyPeriodStr:  ScalingPolicyPeriod,
	ScalingPolicyTargetStr:  ScalingPolicyTarget,
}

var ScalingPolicyTypeDictRev = map[int32]string{
	ScalingPolicyAlert:   ScalingPolicyAlertStr,
	ScalingPolicyRegular: ScalingPolicyRegularStr,
	ScalingPolicyPeriod:  ScalingPolicyPeriodStr,
	ScalingPolicyTarget:  ScalingPolicyTargetStr,
}

// 操作单位数组和映射
var OperateUnits = []string{
	ScalingPolicyOperateUnitCountStr,
	ScalingPolicyOperateUnitPercentStr,
}

var OperateUnitDict = map[string]int32{
	ScalingPolicyOperateUnitCountStr:   ScalingPolicyOperateUnitCount,
	ScalingPolicyOperateUnitPercentStr: ScalingPolicyOperateUnitPercent,
}

var OperateUnitDictRev = map[int32]string{
	ScalingPolicyOperateUnitCount:   ScalingPolicyOperateUnitCountStr,
	ScalingPolicyOperateUnitPercent: ScalingPolicyOperateUnitPercentStr,
	ScalingPolicyOperateUnitNull:    defaultStr,
}

// 执行动作数组和映射
var Actions = []string{
	ActionIncreaseStr,
	ActionDecreaseStr,
	ActionSetStr,
}

var ActionDict = map[string]int32{
	ActionIncreaseStr: ActionIncrease,
	ActionDecreaseStr: ActionDecrease,
	ActionSetStr:      ActionSet,
}

var ActionDictRev = map[int32]string{
	ActionIncrease: ActionIncreaseStr,
	ActionDecrease: ActionDecreaseStr,
	ActionSet:      ActionSetStr,
	ActionNull:     defaultStr,
}

// 循环方式数组和映射
var Cycles = []string{
	CycleMonthlyStr,
	CycleWeeklyStr,
	CycleDailyStr,
}

var CycleDict = map[string]int32{
	CycleMonthlyStr: CycleMonthly,
	CycleWeeklyStr:  CycleWeekly,
	CycleDailyStr:   CycleDaily,
}

var CycleDictRev = map[int32]string{
	CycleMonthly: CycleMonthlyStr,
	CycleWeekly:  CycleWeeklyStr,
	CycleDaily:   CycleDailyStr,
	CycleNull:    defaultStr,
}

var ScalingPolicyStatuses = []string{
	StatusEnabledStr,
	StatusDisabledStr,
}

var ScalingPolicyStatusDict = map[string]int32{
	StatusEnabledStr:  StatusEnabled,
	StatusDisabledStr: StatusDisabled,
}

var ScalingPolicyStatusDictRev = map[int32]string{
	StatusEnabled:  StatusEnabledStr,
	StatusDisabled: StatusDisabledStr,
}

//var ScalingActivityStatus = []string{
//	ScalingActivityEnabledStr,
//	ScalingActivityMovingInStr,
//	ScalingActivityMovingOutStr,
//}

var ScalingActivityStatusMapRev = map[int32]string{
	ScalingActivityEnabled:   ScalingActivityEnabledStr,
	ScalingActivityMovingIn:  ScalingActivityMovingInStr,
	ScalingActivityMovingOut: ScalingActivityMovingOutStr,
}

var ScalingActivityStatusMap = map[string]int32{
	ScalingActivityEnabledStr:   ScalingActivityEnabled,
	ScalingActivityMovingInStr:  ScalingActivityMovingIn,
	ScalingActivityMovingOutStr: ScalingActivityMovingOut,
}

var ExecutionModes = []string{
	ExecutionModeAutoStrategyStr,
	ExecutionModeManualStrategyStr,
	ExecutionModeManualAddInstancesStr,
	ExecutionModeManualRemoveInstancesStr,
	ExecutionModeSatisfyMinSizeStr,
	ExecutionModeAdjustToLimitsStr,
	ExecutionModeHealthCheckAddStr,
	ExecutionModeHealthCheckRemoveStr,
}

// 从字符串到代码值的映射
var ExecutionModeToCode = map[string]int32{
	ExecutionModeAutoStrategyStr:          ExecutionModeAutoStrategy,
	ExecutionModeManualStrategyStr:        ExecutionModeManualStrategy,
	ExecutionModeManualAddInstancesStr:    ExecutionModeManualAddInstances,
	ExecutionModeManualRemoveInstancesStr: ExecutionModeManualRemoveInstances,
	ExecutionModeSatisfyMinSizeStr:        ExecutionModeSatisfyMinSize,
	ExecutionModeAdjustToLimitsStr:        ExecutionModeAdjustToLimits,
	ExecutionModeHealthCheckAddStr:        ExecutionModeHealthCheckAdd,
	ExecutionModeHealthCheckRemoveStr:     ExecutionModeHealthCheckRemove,
}

// 从代码值到字符串的映射
var ExecutionModeToString = map[int32]string{
	ExecutionModeAutoStrategy:          ExecutionModeAutoStrategyStr,
	ExecutionModeManualStrategy:        ExecutionModeManualStrategyStr,
	ExecutionModeManualAddInstances:    ExecutionModeManualAddInstancesStr,
	ExecutionModeManualRemoveInstances: ExecutionModeManualRemoveInstancesStr,
	ExecutionModeSatisfyMinSize:        ExecutionModeSatisfyMinSizeStr,
	ExecutionModeAdjustToLimits:        ExecutionModeAdjustToLimitsStr,
	ExecutionModeHealthCheckAdd:        ExecutionModeHealthCheckAddStr,
	ExecutionModeHealthCheckRemove:     ExecutionModeHealthCheckRemoveStr,
}

// 从字符串表示到代码值的映射
var HealthStatusStrToCode = map[string]int32{
	HealthStatusNormalStr:       HealthStatusNormal,
	HealthStatusAbnormalStr:     HealthStatusAbnormal,
	HealthStatusInitializingStr: HealthStatusInitializing,
}

// 从代码值到字符串表示的映射
var HealthStatusCodeToStr = map[int32]string{
	HealthStatusNormal:       HealthStatusNormalStr,
	HealthStatusAbnormal:     HealthStatusAbnormalStr,
	HealthStatusInitializing: HealthStatusInitializingStr,
}

// 从字符串表示到代码值的映射
var ProtectStatusStrToCode = map[string]int32{
	ProtectStatusProtectedStr:   ProtectStatusProtected,
	ProtectStatusUnprotectedStr: ProtectStatusUnprotected,
}

// 从代码值到字符串表示的映射
var ProtectStatusCodeToStr = map[int32]string{
	ProtectStatusProtected:   ProtectStatusProtectedStr,
	ProtectStatusUnprotected: ProtectStatusUnprotectedStr,
}
