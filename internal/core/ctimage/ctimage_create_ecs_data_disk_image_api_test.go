package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageCreateEcsDataDiskImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageCreateEcsDataDiskImageApi

	// 构造请求
	var enableImageIntegrityCheck bool = false
	request := &CtimageCreateEcsDataDiskImageRequest{
		DataDiskID:                "8888a888-b888-8888-a888-baee8d8ce88c",
		ImageName:                 "CTyunOS-test",
		InstanceID:                "88f888ea-88ff-88ec-a8bc-888888888fe8",
		RegionID:                  "88f8888888dd88ec888888888d888d8b",
		Description:               "Test CTyunOS",
		EnableImageIntegrityCheck: &enableImageIntegrityCheck,
		Labels: []*CtimageCreateEcsDataDiskImageLabelsRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		ProjectID: "0",
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
