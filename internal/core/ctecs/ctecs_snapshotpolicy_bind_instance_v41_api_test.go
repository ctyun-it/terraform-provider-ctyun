package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsSnapshotpolicyBindInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsSnapshotpolicyBindInstanceV41Api

	// 构造请求
	request := &CtecsSnapshotpolicyBindInstanceV41Request{
		RegionID:         "81f7728662dd11ec810800155d307d5b",
		SnapshotPolicyID: "4f69f096066011ee9caf0242ac110002",
		InstanceIDs:      "8d130fba-a3f3-c434-2855-283c96782545",
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
