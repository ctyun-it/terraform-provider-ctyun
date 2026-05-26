package business

const (
	// RocketMQ实例状态定义
	RocketMqStatusRunning      = 1 // 运行中
	RocketMqStatusUnsubscribed = 4 // 已退订

	// 计费模式
	RocketMqBillModePrepaid  = 1 // 预付费
	RocketMqBillModePostpaid = 2 // 后付费
)

// RocketMqInstance RocketMQ实例结构体
type RocketMqInstance struct {
	ProdInstId  string // 实例ID
	ClusterName string // 集群名称
	SpecName    string // 规格名称
	NodeNum     int32  // 节点数
	DiskSize    string // 磁盘大小
	BillMode    string // 计费模式
	Status      int32  // 状态
	CreateTime  string // 创建时间
	ExpireTime  string // 到期时间
	Endpoint    string // 接入点
}
