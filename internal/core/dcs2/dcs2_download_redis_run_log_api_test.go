package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DownloadRedisRunLogApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DownloadRedisRunLogApi

	// 构造请求
	request := &Dcs2DownloadRedisRunLogRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "60bbf68d37c440049cc20fc5a0b16fbf",
		NodeName:   "redis_6379_59889",
		LogType:    "INFO",
		Date:       "2024-11-11",
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
