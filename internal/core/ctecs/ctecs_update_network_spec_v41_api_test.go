package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateNetworkSpecV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateNetworkSpecV41Api

	// 构造请求
	request := &CtecsUpdateNetworkSpecV41Request{
		RegionID:    "a6449feab4db11e9a6b40242ac110007",
		InstanceID:  "93366056-b08f-4b9b-8e47-c50d92f2d4fd",
		Bandwidth:   100,
		ClientToken: "ea1b9004-f450-11ec-8d4f-00155de3fd73",
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
