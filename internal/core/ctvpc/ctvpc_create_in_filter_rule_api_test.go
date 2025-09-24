package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateInFilterRuleApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateInFilterRuleApi

	// 构造请求
	request := &CtvpcCreateInFilterRuleRequest{
		RegionID:         "81f7728662dd11ec810800155d307d5b",
		MirrorFilterID:   "mrfilter-2ck692kf8o",
		DestCidr:         "0.0.0.0/0",
		SrcCidr:          "0.0.0.0/0",
		DestPort:         "-",
		SrcPort:          "-",
		Protocol:         "all",
		EnableCollection: true,
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
