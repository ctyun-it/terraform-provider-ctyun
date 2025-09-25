package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsSetDeletePolicyEbsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsSetDeletePolicyEbsApi

	// 构造请求
	request := &EbsSetDeletePolicyEbsRequest{
		RegionID:          "41f64827f25f468595ffa3a5deb5d15d",
		DiskID:            "4afo4ffa-7fd0-4307-822-69808d2a5cb",
		DeleteSnapWithEbs: true,
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
