package ctimage

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtimageListImageSharesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	//credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	credential := core.CredentialFromEnv()
	apis := NewApis("https://ctimage-global.ctapi-test.ctyun.cn:21443", client)
	api := apis.CtimageListImageSharesApi

	// 构造请求
	request := &CtimageListImageSharesRequest{
		//ImageID:  "5471b138-9502-46bb-859d-67e5ae48776e",
		ImageID:  "asfdfasdf",
		RegionID: "81f7728662dd11ec810800155d307d5b",
		PageNo:   "1",
		PageSize: "10",
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
