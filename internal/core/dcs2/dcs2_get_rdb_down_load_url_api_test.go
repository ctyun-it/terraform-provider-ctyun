package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2GetRdbDownLoadUrlApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2GetRdbDownLoadUrlApi

	// 构造请求
	request := &Dcs2GetRdbDownLoadUrlRequest{
		RegionId:    "bb9fdb42056f11eda1610242ac110002",
		ProdInstId:  "e64c12eb7f8748f7b23b6d90c593796b",
		RestoreName: "20240606125216",
		IpType:      "publicIp",
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
