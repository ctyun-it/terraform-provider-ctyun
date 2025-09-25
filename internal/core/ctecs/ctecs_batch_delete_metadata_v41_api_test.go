package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsBatchDeleteMetadataV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsBatchDeleteMetadataV41Api

	// 构造请求
	request := &CtecsBatchDeleteMetadataV41Request{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		InstanceIDList: "88f888ea-88ff-88ec-a8bc-888888888fe8,a8f8d8c8-88fd-f88a-888b-c8888adf8da8",
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
