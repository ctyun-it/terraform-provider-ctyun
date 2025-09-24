package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsBatchUpdateSnapshotV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsBatchUpdateSnapshotV41Api

	// 构造请求
	request := &CtecsBatchUpdateSnapshotV41Request{
		RegionID: "bb9fdb42056f11eda1610242ac110002",
		UpdateInfo: []*CtecsBatchUpdateSnapshotV41UpdateInfoRequest{
			{
				SnapshotID:          "d4d4ee2b-1478-0e0b-53e6-10e738cfc58c",
				SnapshotName:        "snapshot_update_batch01",
				SnapshotDescription: "snapshot_update_des",
			},
		},
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
