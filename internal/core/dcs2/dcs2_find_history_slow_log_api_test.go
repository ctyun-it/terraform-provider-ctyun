package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2FindHistorySlowLogApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2FindHistorySlowLogApi

	// 构造请求
	request := &Dcs2FindHistorySlowLogRequest{
		RegionId:     "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:   "60ae39640df659abee5ab46b0c6aff9f",
		NodeName:     "access_24318",
		RedisSetName: "",
		StartTime:    "2023-09-26 17:00:00",
		EndTime:      "2023-09-26 17:15:00",
		Page:         1,
		Rows:         10,
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
