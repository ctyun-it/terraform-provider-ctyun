package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2AddValueApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2AddValueApi

	// 构造请求
	request := &Dcs2AddValueRequest{
		RegionId:   "bb9fdb42056f11eda1610242ac110002",
		ProdInstId: "60ae39640df659abee5ab46b0c6aff9f",
		GroupId:    "0",
		Key:        "hkey",
		RawType:    "hash",
		NewKey:     "aaa",
		Score:      "1.0",
		NewValue:   "v2-1",
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
