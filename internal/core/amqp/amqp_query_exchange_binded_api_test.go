package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpQueryExchangeBindedApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	yourEndpoint := "<YOUR_ENDPOINT>"
	apis := NewApis(yourEndpoint, client)
	api := apis.AmqpQueryExchangeBindedApi

	// 构造请求
	request := &AmqpQueryExchangeBindedRequest{
		ProdInstId: "",
		Vhost:      "",
		Exchange:   "",
	}

	// 发起调用
	api.Do(credential, apis, yourEndpoint, request)
}
