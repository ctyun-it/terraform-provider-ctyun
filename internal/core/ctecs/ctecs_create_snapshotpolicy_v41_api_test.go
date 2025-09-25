package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateSnapshotpolicyV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateSnapshotpolicyV41Api

	// 构造请求
	request := &CtecsCreateSnapshotpolicyV41Request{
		RegionID:             "81f7728662dd11ec810800155d307d5b",
		SnapshotPolicyName:   "api-create01",
		SnapshotTime:         "12,13",
		CycleType:            "day",
		CycleDay:             1,
		CycleWeek:            "0,2,6",
		RetentionType:        "num",
		RetentionDay:         2,
		RetentionNum:         3,
		SnapshotPolicyStatus: 1,
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
