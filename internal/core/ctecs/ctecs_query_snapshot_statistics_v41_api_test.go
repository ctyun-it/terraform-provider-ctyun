package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQuerySnapshotStatisticsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQuerySnapshotStatisticsV41Api

	// 构造请求
	request := &CtecsQuerySnapshotStatisticsV41Request{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		InstanceIDList: "8986fc25-4f0a-4fc6-03ef-71c386b49905, 88a6a862-9030-e91b-3843-0144cd5e5dff",
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
