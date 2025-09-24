package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsAttachVolumeV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsAttachVolumeV41Api

	// 构造请求
	request := &CtecsAttachVolumeV41Request{
		DiskID:     "6f8928c7-f961-4ece-b0ee-9f8d6b4663b5",
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "c54322f1-735d-409d-2a7c-76d611492469",
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
