package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestHpfsNewDirectoryApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.HpfsNewDirectoryApi

	// 构造请求
	request := &HpfsNewDirectoryRequest{
		RegionID:         "81f7728662dd11ec810800155d307d5b",
		SfsUID:           "484bd9a4-65fc-57ba-a562-743e6073c14d",
		SfsDirectory:     "/test/test01",
		SfsDirectoryMode: "755",
		SfsDirectoryUID:  0,
		SfsDirectoryGID:  0,
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
