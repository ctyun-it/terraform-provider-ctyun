package business

const (
	// 1表示运行中，3表示已注销，4表示已退订，5表示变更中，6表示创建中
	RabbitMqStatusRunning      = 1
	RabbitMqStatusUnsubscribed = 4

	RabbitmqExchangeTypeDirect          = "direct"
	RabbitMqExchangeTypeTopic           = "topic"
	RabbitMqExchangeTypeFanout          = "fanout"
	RabbitMqExchangeTypeHeaders         = "headers"
	RabbitMqExchangeTypeXDelayedMessage = "x-delayed-message"
)

var RabbitMqExchangeTypes = []string{
	RabbitmqExchangeTypeDirect,
	RabbitMqExchangeTypeTopic,
	RabbitMqExchangeTypeFanout,
	RabbitMqExchangeTypeHeaders,
	RabbitMqExchangeTypeXDelayedMessage,
}

var RabbitMqExchangeXDelayedTypes = []string{
	RabbitmqExchangeTypeDirect,
	RabbitMqExchangeTypeTopic,
	RabbitMqExchangeTypeFanout,
	RabbitMqExchangeTypeHeaders,
}
