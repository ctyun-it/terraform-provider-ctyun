package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsPortsDetachInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsPortsDetachInstanceV41Api

	// 构造请求
	request := &CtecsPortsDetachInstanceV41Request{
		ClientToken:        "ports-detach-01",
		RegionID:           "bb9fdb42056f11eda1610242ac110002",
		NetworkInterfaceID: "port-vibsmse8pl",
		InstanceID:         "a628a7d9-ef97-3b16-8a0a-4a794fcdbc39",
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
