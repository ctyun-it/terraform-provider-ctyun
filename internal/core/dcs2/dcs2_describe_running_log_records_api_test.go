package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DescribeRunningLogRecordsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DescribeRunningLogRecordsApi

	// 构造请求
	request := &Dcs2DescribeRunningLogRecordsRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "a0784f6c27cf0139e79e055f89f03f1d",
		NodeName:   "redis_30101_49609",
		PageSize:   10,
		PageIndex:  1,
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
