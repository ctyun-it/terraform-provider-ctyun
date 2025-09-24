package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsKeypairAttachInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsKeypairAttachInstanceV41Api

	// 构造请求
	request := &CtecsKeypairAttachInstanceV41Request{
		RegionID:    "bb9fdb42056f11eda1610242ac110002",
		KeyPairName: "KeyPair-886",
		InstanceID:  "b6e2966d-7b1c-385e-abe4-d940caa273b7",
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
