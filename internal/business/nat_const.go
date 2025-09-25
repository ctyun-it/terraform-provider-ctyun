package business

const (
	// 订购周期常量
	MonthCycleType    = "month"
	YearCycleType     = "year"
	OnDemandCycleType = "on_demand"

	//资源状态常量
	NatStatusStarted    = "started"     //启用
	NatStatusRenewed    = "renewed"     //续订
	NatStatusRefunded   = "refunded"    //退订
	NatStatusDestroyed  = "destroyed"   //销毁
	NatStatusFailed     = "failed"      //失败
	NatStatusStarting   = "starting"    //正在启动
	NatStatusChanged    = "changed"     //变配
	NatStatusExpired    = "expired"     //过期
	NatStatusUnknown    = "unknown"     //未知
	NatStatusInProgress = "in_progress" //

	//  Nat规格
	SpecSmall      = 1
	SpecMedium     = 2
	SpecLarge      = 3
	SpecExtraLarge = 4

	// protocol
	ProtocolTcp = "tcp"
	ProtocolUdp = "udp"

	ProtocolTCP      = "TCP"
	ProtocolUDP      = "UDP"
	NatStateRunning  = "running"  // 运行中
	NatStateCreating = "creating" // 创建中
	NatStateExpired  = "expired"  // 已过期
	NatStateFreeze   = "freeze"   //已冻结

	//DNAT运行状态
	DNatStateActive   = "active"
	DNatStateFreezing = "freezing"
	DNatStateCreating = "creating"

	//DNAT运行状态
	DNatStateACTIVE   = "ACTIVE"
	DNatStateFREEZING = "FREEZING"
	DNatStateCREATING = "CREATING"

	//SNAT状态
	SNatStatusACTIVE   = "ACTIVE"
	SNatStatusCreating = "Creating"
	SNatStatusFreezing = "Freezing"

	SubnetTypeVPC    = 1
	SubnetTypeCustom = 2

	// SNAT创建状态
	NatCreateStatusING  = "in_progress"
	NatCreateStatusDone = "done"

	// DNAT类型
	VirtualMachineTypeCloud  = "instance" //服务器
	VirtualMachineTypeCustom = "custom"   //自定义

	// DNAT serverType
	ServerTypeVM = "VM"
	ServerTypeBM = "BM"
)

var NatOrderCycleTypes = []string{
	MonthCycleType,
	YearCycleType,
	OnDemandCycleType,
}

var NatStatus = []string{
	NatStatusStarted,
	NatStatusRenewed,
	NatStatusRefunded,
	NatStatusDestroyed,
	NatStatusFailed,
	NatStatusStarting,
	NatStatusChanged,
	NatStatusExpired,
	NatStatusUnknown,
}

var NatSpecs = []int64{
	SpecSmall,
	SpecMedium,
	SpecLarge,
	SpecExtraLarge,
}

var DNatProtocols = []string{
	ProtocolTcp,
	ProtocolUdp,
}
var DNatStatus = []string{
	DNatStateActive,
	DNatStateFreezing,
	DNatStateCreating,
}
var DNatStates = []string{
	DNatStateACTIVE,
	DNatStateCREATING,
	DNatStateFREEZING,
}

var SNatSubnetTypes = []int32{
	SubnetTypeVPC,
	SubnetTypeCustom,
}

var ServerTypes = []string{
	ServerTypeVM,
	ServerTypeBM,
}

var NatStates = []string{
	NatStateRunning,
	NatStateCreating,
	NatStateExpired,
	NatStateFreeze,
}

var SNatStatus = []string{
	SNatStatusACTIVE,
	SNatStatusCreating,
	SNatStatusFreezing,
}

var SNatProtocols = []string{
	ProtocolTCP,
	ProtocolUDP,
}
var DnatStatus = []string{
	NatCreateStatusING,
	NatCreateStatusDone,
}
