package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmStopInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmStopInstanceApi

	// 构造请求
	request := &EbmStopInstanceRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		AzName:       "az1",
		InstanceUUID: "ss-9d4u1yd1jr0a3xeu59fq9svecxql",
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
