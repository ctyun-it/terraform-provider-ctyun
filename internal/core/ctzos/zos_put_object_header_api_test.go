package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutObjectHeaderApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutObjectHeaderApi

	// 构造请求
	request := &ZosPutObjectHeaderRequest{
		Bucket:   "bucket1",
		Key:      "obj1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		Headers: &ZosPutObjectHeaderHeadersRequest{
			CacheControl:       "max-age=604800",
			ContentDisposition: "attachment; filename=example.pdf",
			ContentEncoding:    "gzip",
			ContentLanguage:    "en-US",
			ContentType:        "application/json",
			Expires:            "Sat, 31 Dec 2023 23:59:59 GMT",
		},
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
