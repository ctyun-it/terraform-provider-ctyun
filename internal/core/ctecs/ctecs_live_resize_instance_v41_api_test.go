package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsLiveResizeInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsLiveResizeInstanceV41Api

	// 构造请求
	request := &CtecsLiveResizeInstanceV41Request{
		RegionID:        "41f64827f25f468595ffa3a5deb5d15d",
		InstanceID:      "285010af-16f1-137e-06c0-920d4bdd0026",
		FlavorID:        "s2.small.1或00ebe3aa-aac0-1d99-0b9e-4d391c5e06d5",
		ClientToken:     "resize3003",
		PayVoucherPrice: 20.55,
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
