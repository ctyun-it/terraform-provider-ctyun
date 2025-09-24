package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsRebuildLiteInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsRebuildLiteInstanceApi

	// 构造请求
	request := &CtecsRebuildLiteInstanceRequest{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "adc614e0-e838-d73f-0618-a6d51d09070a",
		Password:   "rebuildTest01",
		ImageID:    "b1d896e1-c977-4fd4-b6c2-5432549977be",
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
