package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcVniaCreateNetworkPathApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcVniaCreateNetworkPathApi

	// 构造请求
	request := &CtvpcVniaCreateNetworkPathRequest{
		RegionID:   "81f7728662dd11ec810800155d307d5b",
		Name:       "test",
		SourceID:   "64851e6e-1f92-fe42-3d71-9a1632a45e62",
		SourceType: "ecs",
		SourcePort: 80,
		TargetType: "ecs",
		TargetPort: 80,
		TargetID:   "64851e6e-1f92-fe42-3d71-9a1632a45e62",
		SourceIP:   "192.168.0.1",
		TargetIP:   "192.168.0.2",
		Protocol:   "ICMP",
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
