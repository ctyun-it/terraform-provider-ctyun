package amqp

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestAmqpExchangeQueryV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.AmqpExchangeQueryV3Api

	// 构造请求
	request := &AmqpExchangeQueryV3Request{
		RegionId:   "6b10c8b962244f1f921d1a48d0f15cca",
		ProdInstId: "6b10c8b962244f1f921d1a48d0f15ccaa",
		Vhost:      "vhost1",
		PageNum:    "1",
		PageSize:   "100",
		Name:       "ex1",
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
