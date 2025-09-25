package business

import "github.com/ctyun-it/terraform-provider-ctyun/internal/utils"

const (
	OrderStatusWaitingToPay    = 1   // 待支付
	OrderStatusPayed           = 2   // 已支付
	OrderStatusFinish          = 3   // 完成
	OrderStatusCanceled        = 4   // 取消
	OrderStatusFail            = 5   // 施工失败
	OrderStatusPaying          = 7   // 正在支付中
	OrderStatusWaitingToAduit  = 8   // 待审核
	OrderStatusAduitedSuccess  = 9   // 审核通过
	OrderStatusAuditedFail     = 10  // 审核未通过
	OrderStatusRevoked         = 11  // 撤单完成
	OrderStatusUnsubscribing   = 12  // 退订中
	OrderStatusUnsubscribed    = 13  // 退订完成
	OrderStatusOpening         = 14  // 开通中
	OrderStatusChangeRemoved   = 15  // 变更移除
	OrderStatusAutoRevoking    = 16  // 自动撤单中
	OrderStatusManualRevoking  = 17  // 手动撤单中
	OrderStatusAborting        = 18  // 终止中
	OrderStatusPayFail         = 22  // 支付失败
	OrderStatusWaitingToRevoke = -2  // 待撤回
	OrderStatusUnknown         = -1  // 未知
	OrderStatusError           = 0   // 错误
	OrderStatusDeleted         = 999 // 逻辑删除

	OrderCycleTypeMonth    = "month"
	OrderCycleTypeYear     = "year"
	OrderCycleTypeOnDemand = "on_demand"
)

var OrderStatusName = map[int]string{
	1:   "待支付",
	2:   "已支付",
	3:   "完成",
	4:   "取消",
	5:   "施工失败",
	7:   "正在支付中",
	8:   "待审核",
	9:   "审核通过",
	10:  "审核未通过",
	11:  "撤单完成",
	12:  "退订中",
	13:  "退订完成",
	14:  "开通中",
	15:  "变更移除",
	16:  "自动撤单中",
	17:  "手动撤单中",
	18:  "终止中",
	22:  "支付失败",
	-2:  "待撤回",
	-1:  "未知",
	0:   "错误",
	999: "逻辑删除",
}

const (
	OrderCycleTypeMapScene1 = iota
)

var OrderCycleTypes = []string{
	OrderCycleTypeMonth,
	OrderCycleTypeYear,
	OrderCycleTypeOnDemand,
}

var OrderCycleTypeMap = utils.Must(
	[]any{
		OrderCycleTypeMonth,
		OrderCycleTypeYear,
		OrderCycleTypeOnDemand,
	},
	map[utils.Scene][]any{
		OrderCycleTypeMapScene1: {
			"MONTH",
			"YEAR",
			"",
		},
	},
)
