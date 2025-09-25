package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmMetadataBatchCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmMetadataBatchCreateApi

	// 构造请求
	request := &EbmMetadataBatchCreateRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		AzName:       "az1",
		InstanceUUID: "ss-9d4u1yd1jr0a3xeu59fq9svecxql",
		Metadata:     &EbmMetadataBatchCreateMetadataRequest{},
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
