package rocketmq

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestRocketmqInstQueryDetailV3Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.RocketmqInstQueryDetailV3Api

	// 构造请求
	request := &RocketmqInstQueryDetailV3Request{
		RegionId:   "80a4c94407304e1c8bc1ab15faef6e12",
		ProdInstId: "80a4c94407304e1c8bc1ab15faef6e12",
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
