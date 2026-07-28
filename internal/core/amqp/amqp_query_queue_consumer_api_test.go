package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpQueryQueueConsumerApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	yourEndpoint := "<YOUR_ENDPOINT>"
	apis := NewApis(yourEndpoint, client)
	api := apis.AmqpQueryQueueConsumerApi

	// 构造请求
	request := &AmqpQueryQueueConsumerRequest{
		ProdInstId: "",
		Vhost:      "",
		Name:       "",
	}

	// 发起调用
	api.Do(credential, apis, yourEndpoint, request)
}
