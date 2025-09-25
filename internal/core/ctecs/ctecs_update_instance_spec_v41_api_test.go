package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateInstanceSpecV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateInstanceSpecV41Api

	// 构造请求
	request := &CtecsUpdateInstanceSpecV41Request{
		RegionID:    "88f8888888dd88ec888888888d888d8b",
		InstanceID:  "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		Bandwidth:   100,
		FlavorID:    "00ebe3aa-aac0-1d99-0b9e-4d391c5e06d5",
		ClientToken: "bdfse888-8ed8-88b8-88cb-888f8b8cf8fa",
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
