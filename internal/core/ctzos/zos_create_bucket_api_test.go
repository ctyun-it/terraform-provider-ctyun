package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosCreateBucketApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosCreateBucketApi

	// 构造请求
	var isEncrypted bool = false
	var objectLockEnabledForBucket bool = true
	request := &ZosCreateBucketRequest{
		RegionID:    "332232eb-63aa-465e-9028-52e5123866f0",
		ACL:         "private",
		Bucket:      "bucket1",
		ProjectID:   "437367233",
		CmkUUID:     "wasgdfeb-63aa-465e-9028-52e5123866f0",
		IsEncrypted: &isEncrypted,
		StorageType: "STANDARD",
		AZPolicy:    "single-az",
		Labels: []*ZosCreateBucketLabelsRequest{
			{
				Key:   "format",
				Value: "aac",
			},
		},
		OtherBucketInfo: &ZosCreateBucketOtherBucketInfoRequest{
			ObjectLockEnabledForBucket: &objectLockEnabledForBucket,
		},
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
