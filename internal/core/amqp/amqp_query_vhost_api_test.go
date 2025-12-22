package amqp

import (
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpQueryVhostApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	yourEndpoint := "<YOUR_ENDPOINT>"
	apis := NewApis(yourEndpoint, client)
	api := apis.AmqpQueryVhostApi

	// 构造请求
	request := &AmqpQueryVhostRequest{
		ProdInstId: "",
	}

	// 发起调用
	api.Do(credential, apis, yourEndpoint, request)
}
