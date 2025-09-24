package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsCreateEbsSnapApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsCreateEbsSnapApi

	// 构造请求
	var clientToken string = "tt-yy-sss"
	request := &EbsCreateEbsSnapRequest{
		ClientToken:     &clientToken,
		RegionID:        "fc862f71-d629-4a0e-9fe0-b104963b3f05",
		SnapshotName:    "snapshot-001",
		DiskID:          "36e88a58-1ebf-40ac-91b6-a8c0eca38314",
		RetentionPolicy: "forever",
		RetentionTime:   7,
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
