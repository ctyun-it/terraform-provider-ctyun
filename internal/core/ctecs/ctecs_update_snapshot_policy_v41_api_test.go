package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateSnapshotPolicyV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateSnapshotPolicyV41Api

	// 构造请求
	request := &CtecsUpdateSnapshotPolicyV41Request{
		RegionID:           "81f7728662dd11ec810800155d307d5b",
		SnapshotPolicyID:   "4f69f096066011ee9caf0242ac110002",
		SnapshotPolicyName: "update-test01",
		SnapshotTime:       "12,13",
		CycleType:          "week",
		CycleDay:           1,
		CycleWeek:          "0,2,6",
		RetentionType:      "date",
		RetentionDay:       2,
		RetentionNum:       3,
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
