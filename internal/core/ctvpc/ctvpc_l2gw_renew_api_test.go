package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcL2gwRenewApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcL2gwRenewApi

	// 构造请求
	request := &CtvpcL2gwRenewRequest{
		ClientToken: "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:    "",
		CycleType:   "month",
		CycleCount:  1,
		L2gwID:      "l2gw-xxxx",
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
