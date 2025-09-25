package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseGetClusterEventsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseGetClusterEventsApi

	// 构造请求
	request := &CcseGetClusterEventsRequest{
		ClusterId:  "47281b02f87757478f20b1827c97cadf",
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		EventType:  "目前支持的：node_pool_config_sync; node_pool_kubelet_config_update;node_pool_repair;node_repair",
		TaskId:     "123",
		PageNumber: 1,
		PageSize:   10,
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
