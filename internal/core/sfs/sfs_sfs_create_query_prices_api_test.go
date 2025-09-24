package sfs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestSfsSfsCreateQueryPricesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.SfsSfsCreateQueryPricesApi

	// 构造请求
	request := &SfsSfsCreateQueryPricesRequest{
		RegionID:   "参考[请求示例]",
		OrderNum:   3,
		CycleType:  "参考[请求示例]",
		SfsSize:    500,
		VolumeType: "参考[请求示例]",
		CycleCnt:   3,
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
