package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmBatchDeleteInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmBatchDeleteInstancesApi

	// 构造请求
	request := &EbmBatchDeleteInstancesRequest{
		RegionID:         "81f7728662dd11ec810800155d307d5b",
		AzName:           "az1",
		InstanceUUIDList: "ss-eztbrc1j541lb2pim0lxr3uctdfo,ss-0bmsde9loclfnzhck4mj4g9ck9hg",
		ClientToken:      "ea1b9004-f450-11ec-8d4f-00155de3fd73",
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
