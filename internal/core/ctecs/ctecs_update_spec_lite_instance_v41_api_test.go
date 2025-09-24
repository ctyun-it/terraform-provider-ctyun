package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateSpecLiteInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateSpecLiteInstanceV41Api

	// 构造请求
	request := &CtecsUpdateSpecLiteInstanceV41Request{
		ClientToken:   "bdfse888-8ed8-88b8-88cb-888f8b8cf8fa",
		RegionID:      "bb9fdb42056f11eda1610242ac110002",
		InstanceID:    "adc614e0-e838-d73f-0618-a6d51d09070a",
		FlavorSetType: "fix",
		FlavorName:    "lite1.fix.small.1",
		BootDiskSize:  40,
		Bandwidth:     5,
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
