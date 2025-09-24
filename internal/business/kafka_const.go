package business

const (
	KafkaStatusRunning        = 1   // 运行中
	KafkaStatusExpired        = 2   // 已过期
	KafkaStatusDeregister     = 3   // 已注销
	KafkaStatusChanging       = 4   // 变更中
	KafkaStatusUnsubscribed   = 5   // 已退订
	KafkaStatusActivating     = 6   // 开通中
	KafkaStatusCanceled       = 7   // 已取消
	KafkaStatusStopped        = 8   // 已停止
	KafkaStatusEIPProcessing  = 9   // 弹性IP处理中
	KafkaStatusRestarting     = 10  // 重启中
	KafkaStatusRestartFailed  = 11  // 重启失败
	KafkaStatusUpgrading      = 12  // 升级中
	KafkaStatusArrears        = 13  // 已欠费
	KafkaStatusActivateFailed = 101 // 开通失败

	KafkaVersion28 = "2.8"
	KafkaVersion36 = "3.6"
)
