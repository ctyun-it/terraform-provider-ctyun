package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsModifyPolicyStatusEbsSnapApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsModifyPolicyStatusEbsSnapApi

	// 构造请求
	request := &EbsModifyPolicyStatusEbsSnapRequest{
		RegionID:         "41f64827f25f468595ffa3a5deb5d15d",
		SnapshotPolicyID: "0ae97ef5-6ee2-44af-9d05-1a509b0a1be6",
		TargetStatus:     "activated",
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
