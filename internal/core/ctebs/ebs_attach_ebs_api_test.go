package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsAttachEbsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsAttachEbsApi

	// 构造请求
	var regionID string = "81f7728662dd11ec810800155d307d5b"
	request := &EbsAttachEbsRequest{
		DiskID:     "65544165-c658-45c3-a31a-426c39929151",
		RegionID:   &regionID,
		InstanceID: "24690060-c475-ed64-fd2c-7e96f9a1df37",
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
