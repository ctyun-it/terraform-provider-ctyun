package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosGenerateObjectDownloadLinkApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosGenerateObjectDownloadLinkApi

	// 构造请求
	request := &ZosGenerateObjectDownloadLinkRequest{
		RegionID:  "332232eb-63aa-465e-9028-52e5123866f0",
		Bucket:    "bucket1",
		Key:       "obj1",
		VersionID: "USzD.sN0vODAsQ84ncdT20oRiY2lFCD",
		ExpiresIn: 3300,
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
