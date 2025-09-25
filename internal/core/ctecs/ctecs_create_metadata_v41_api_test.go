package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateMetadataV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateMetadataV41Api

	// 构造请求
	request := &CtecsCreateMetadataV41Request{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "b67b7f1f-095b-1249-b379-8dd5cc542a05",
		Metadata:   &CtecsCreateMetadataV41MetadataRequest{},
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
