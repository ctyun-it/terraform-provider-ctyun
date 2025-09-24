package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsEcsSnapshotTaskListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsEcsSnapshotTaskListApi

	// 构造请求
	request := &CtecsEcsSnapshotTaskListRequest{
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		TaskID:       "6b2ace9e-663b-4cb0-848e-41c36965458c",
		TaskType:     "create",
		InstanceID:   "b721c0a9-c6c1-4b46-82bc-e12158e0dda7",
		InstanceName: "ecm-2619",
		SnapshotID:   "f16a0ae2-6e0c-db97-a037-67c586f4aa93",
		SnapshotName: "sp-01",
		StrategyID:   "ab1ac6e657be4d4497ca1690cec84cf1",
		QueryContent: "",
		PageNo:       1,
		PageSize:     10,
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
