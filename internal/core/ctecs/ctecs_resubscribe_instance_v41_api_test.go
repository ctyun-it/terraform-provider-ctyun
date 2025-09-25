package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsResubscribeInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsResubscribeInstanceV41Api

	// 构造请求
	request := &CtecsResubscribeInstanceV41Request{
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		InstanceID:      "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		CycleCount:      6,
		CycleType:       "MONTH",
		ClientToken:     "4cf2962d-e92c-4c00-9181-cfbb2218636c",
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
