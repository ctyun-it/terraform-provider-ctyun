package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageCreateEcsSnapshotSystemDiskImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageCreateEcsSnapshotSystemDiskImageApi

	// 构造请求
	var enableImageIntegrityCheck bool = false
	request := &CtimageCreateEcsSnapshotSystemDiskImageRequest{
		ImageName:                 "CTyunOS-test",
		RegionID:                  "88f8888888dd88ec888888888d888d8b",
		SnapshotID:                "c8a8f88d-fb8f-8d8a-e888-8888888b8b8d",
		Description:               "Test CTyunOS",
		EnableImageIntegrityCheck: &enableImageIntegrityCheck,
		Labels: []*CtimageCreateEcsSnapshotSystemDiskImageLabelsRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		MaximumRAM: 0,
		MinimumRAM: 0,
		ProjectID:  "0",
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
