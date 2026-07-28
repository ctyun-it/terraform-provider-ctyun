package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpQueryQueueApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	yourEndpoint := "<YOUR_ENDPOINT>"
	apis := NewApis(yourEndpoint, client)
	api := apis.AmqpQueryQueueApi

	// 构造请求
	request := &AmqpQueryQueueRequest{
		ProdInstId: "",
		Vhost:      "",
		Name:       "",
		PageNum:    0,
		PageSize:   0,
	}

	// 发起调用
	api.Do(credential, apis, yourEndpoint, request)
}
