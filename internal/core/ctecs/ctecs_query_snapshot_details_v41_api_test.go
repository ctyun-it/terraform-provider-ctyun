package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQuerySnapshotDetailsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQuerySnapshotDetailsV41Api

	// 构造请求
	request := &CtecsQuerySnapshotDetailsV41Request{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		SnapshotID: "c7a7f06d-fb0f-8d5a-e710-9262995b6b6d",
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
