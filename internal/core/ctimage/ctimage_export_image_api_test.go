package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageExportImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageExportImageApi

	// 构造请求
	request := &CtimageExportImageRequest{
		Bucket:          "bucket-fa88",
		Filename:        "CTyunOS-test",
		ImageID:         "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		RegionID:        "88f8888888dd88ec888888888d888d8b",
		ImageFileFormat: "raw",
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
