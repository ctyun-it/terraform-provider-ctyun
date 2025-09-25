package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsQuerySnapshotListV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsQuerySnapshotListV41Api

	// 构造请求
	request := &CtecsQuerySnapshotListV41Request{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		ProjectID:      "0",
		PageNo:         1,
		PageSize:       10,
		InstanceID:     "c4cb5146-a148-3427-0425-73dd2b81e60a",
		SnapshotStatus: "restoring",
		SnapshotID:     "73e30bcd-119b-9653-f864-a50150434a90",
		QueryContent:   "c4cb5146-a148-3427-0425-73dd2b81e60a",
		SnapshotName:   "snapshot_for_restore",
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
