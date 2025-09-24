package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DescribeNodeHistoryMonitorValuesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DescribeNodeHistoryMonitorValuesApi

	// 构造请求
	request := &Dcs2DescribeNodeHistoryMonitorValuesRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "60ae39640df659abee5ab46b0c6aff9f",
		NodeName:   "redis_30101_51466",
		StartTime:  "2024-05-24 16:15:00",
		EndTime:    "2024-05-24 16:20:00",
		RawType:    "NetworkCardReceivedRate",
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
