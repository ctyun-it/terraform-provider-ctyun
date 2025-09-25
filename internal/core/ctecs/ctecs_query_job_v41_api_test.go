package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQueryJobV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQueryJobV41Api

	// 构造请求
	request := &CtecsQueryJobV41Request{
		RegionID: "bb9fdb42056f11eda1610242ac110002",
		JobIDs:   "a8e88ab8-888e-8888-8b88-c8f88a88e8bf,a8e88ab8-888e-8888-8b88-c8f88a88e8bk",
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
