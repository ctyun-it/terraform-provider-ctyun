package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageDetailImageApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtimageDetailImageApi

	// 构造请求
	var errorFree bool = false
	request := &CtimageDetailImageRequest{
		ImageID:   "8d8e8888-8ed8-88b8-88cb-888f8b8cf8fa",
		RegionID:  "bb9fdb42056f11eda1610242ac110002",
		ErrorFree: &errorFree,
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
