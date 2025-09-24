package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketLoggingApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketLoggingApi

	// 构造请求
	request := &ZosPutBucketLoggingRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		BucketLoggingStatus: &ZosPutBucketLoggingBucketLoggingStatusRequest{
			LoggingEnabled: &ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledRequest{
				TargetPrefix: "nn_",
				TargetBucket: "tbk1",
				TargetGrants: []*ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsRequest{
					{
						Grantee: &ZosPutBucketLoggingBucketLoggingStatusLoggingEnabledTargetGrantsGranteeRequest{
							RawType:      "CanonicalUser",
							EmailAddress: "example@ctyun.cn",
							DisplayName:  "dname1",
							ID:           "a0aaa00a-a0a0-0000-00a0-0a0a000aa0a0",
						},
						Permission: "FULL_CONTROL",
					},
				},
			},
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
