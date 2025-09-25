package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsExtendVolumeV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsExtendVolumeV41Api

	// 构造请求
	request := &CtecsExtendVolumeV41Request{
		DiskSize:    20,
		DiskID:      "0ae97ef5-6ee2-44af-9d05-1a509b0a1be6",
		RegionID:    "bb9fdb42056f11eda1610242ac110002",
		ClientToken: "resize0211v3",
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
