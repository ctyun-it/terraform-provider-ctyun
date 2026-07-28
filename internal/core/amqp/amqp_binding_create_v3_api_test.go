package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpBindingCreateV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.AmqpBindingCreateV3Api

	// 构造请求
	request := &AmqpBindingCreateV3Request{
		RegionId:         "6b10c8b962244f1f921d1a48d0f15cca",
		ProdInstId:       "6b10c8b962244f1f921d1a48d0f15cca",
		Source:           "ex1",
		Destination_type: "q",
		Destination:      "qu1",
		Routing_key:      "key1",
		Vhost:            "vhost1",
		Arguments:        "",
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
