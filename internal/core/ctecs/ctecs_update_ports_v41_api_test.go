package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdatePortsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdatePortsV41Api

	// 构造请求
	request := &CtecsUpdatePortsV41Request{
		ClientToken:        "update-port-test",
		RegionID:           "bb9fdb42056f11eda1610242ac110002",
		NetworkInterfaceID: "port-pja7l0zfvk",
		Name:               "nic-update-name",
		Description:        "nic_update_description",
		SecurityGroupIDs:   []string{},
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
