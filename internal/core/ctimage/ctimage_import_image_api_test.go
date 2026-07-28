package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageImportImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageImportImageApi

	// 构造请求
	var enableImageIntegrityCheck bool = false
	request := &CtimageImportImageRequest{
		ImageFileSource: "https://xxx.zos.ctyun.cn/bucket-xxx/image-xxx",
		ImageProperties: &CtimageImportImageImagePropertiesRequest{
			ImageName:    "CTyunOS-test",
			Architecture: "x86_64",
			BootMode:     "bios",
			Description:  "Test CTyunOS",
			DiskSize:     40,
			ImageType:    "data_disk_image",
			MaximumRAM:   0,
			MinimumRAM:   0,
			OsDistro:     "CTyunOS",
			OsVersion:    "23.01",
		},
		RegionID:                  "88f8888888dd88ec888888888d888d8b",
		EnableImageIntegrityCheck: &enableImageIntegrityCheck,
		Labels: []*CtimageImportImageLabelsRequest{
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
