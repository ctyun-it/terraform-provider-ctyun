package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmAddNicApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmAddNicApi

	// 构造请求
	var ipv4 string = "192.168.0.1"
	request := &EbmAddNicRequest{
		RegionID:       "81f7728662dd11ec810800155d307d5b",
		AzName:         "az1",
		InstanceUUID:   "ss-qyjtenwqho0gfztts5vscc2l5efi",
		SubnetUUID:     "c45512bf-7919-55a9-a106-5f4aa9194c7c",
		SecurityGroups: "71386230-c9e7-5465-85fc-29e26f79806d,71386230-c9e7-5465-85fc-29e26f798062",
		Ipv4:           &ipv4,
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
