package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceDetachKeypairV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceDetachKeypairV41Api

	// 构造请求
	request := &CtecsInstanceDetachKeypairV41Request{
		RegionID:    "bb9fdb42056f11eda1610242ac110002",
		InstanceID:  "39341024-9c57-4be4-a580-ae809d44bafd",
		KeyPairName: "KeyPair-886",
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
