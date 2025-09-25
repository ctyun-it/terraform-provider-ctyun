package hpfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestHpfsResizeSfsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.HpfsResizeSfsApi

	// 构造请求
	request := &HpfsResizeSfsRequest{
		SfsSize:     1024,
		SfsUID:      "e828a6e7-e1ed-5528-968a-9dd7077395d2",
		RegionID:    "81f7728662dd11ec810800155d307d5b",
		ClientToken: "yc-test1234",
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
