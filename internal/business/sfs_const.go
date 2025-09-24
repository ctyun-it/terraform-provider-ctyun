package business

const (
	SfsMonthCycleType    = "month"
	SfsYearCycleType     = "year"
	SfsOnDemandCycleType = "on_demand"

	SfsPermissionGroupRuleUserPermissionNoRootSquash = "no_root_squash" //nfs 访问用户映射:不匿名root用户
)

var SfsCycleType = []string{
	SfsMonthCycleType,
	SfsYearCycleType,
	SfsOnDemandCycleType,
}
