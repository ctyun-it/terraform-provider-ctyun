package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsPortsAttachInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsPortsAttachInstanceV41Api

	// 构造请求
	request := &CtecsPortsAttachInstanceV41Request{
		ClientToken:        "attach_test01",
		RegionID:           "bb9fdb42056f11eda1610242ac110002",
		AzName:             "cn-huadong1-jsnj1A-public-ctcloud",
		ProjectID:          "0",
		NetworkInterfaceID: "subnet-y8cofge5uj",
		InstanceID:         "a628a7d9-ef97-3b16-8a0a-4a794fcdbc39",
		InstanceType:       3,
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
