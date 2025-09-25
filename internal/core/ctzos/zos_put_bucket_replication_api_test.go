package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketReplicationApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketReplicationApi

	// 构造请求
	var plot bool = true
	var history bool = true
	request := &ZosPutBucketReplicationRequest{
		Bucket:         "bucket1",
		RegionID:       "332232eb-63aa-465e-9028-52e5123866f0",
		TargetBucket:   "bucket2",
		TargetRegionID: "3a585beb-63aa-465e-9028-52e5123866f0",
		Prefixes:       []string{},
		Plot:           &plot,
		History:        &history,
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
