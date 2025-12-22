package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpCreateQueueApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.AmqpCreateQueueApi

	// 构造请求
	var durable bool = false
	var auto_delete bool = false
	request := &AmqpCreateQueueRequest{
		ProdInstId:            "",
		Vhost:                 "",
		Name:                  "",
		Durable:               &durable,
		Auto_delete:           &auto_delete,
		XMessageTtl:           0,
		XExpires:              0,
		XMaxLength:            0,
		XDeadLetterExchange:   "",
		XDeadLetterRoutingKey: "",
		XMaxPriority:          0,
	}

	// 发起调用
	response, err := api.Do(context.Background(), *credential, request)
	if err != nil {
		t.Log("request error:", err)
		t.Fail()
		return
	}
	t.Logf("%+v\n", *response)
}
