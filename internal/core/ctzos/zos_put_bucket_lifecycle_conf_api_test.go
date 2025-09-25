package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketLifecycleConfApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketLifecycleConfApi

	// 构造请求
	var expiredObjectDeleteMarker bool = false
	request := &ZosPutBucketLifecycleConfRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		LifecycleConfiguration: &ZosPutBucketLifecycleConfLifecycleConfigurationRequest{
			Rules: []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesRequest{
				{
					ID: "rule1",
					Expiration: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesExpirationRequest{
						Date:                      "2022-10-18T00:00:00Z",
						ExpiredObjectDeleteMarker: &expiredObjectDeleteMarker,
						Days:                      22,
					},
					Status: "Enabled",
					NoncurrentVersionExpiration: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionExpirationRequest{
						NoncurrentDays: 123,
					},
					NoncurrentVersionTransitions: []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesNoncurrentVersionTransitionsRequest{
						{
							NoncurrentDays: 123,
							StorageClass:   "STANDARD_IA",
						},
					},
					Filter: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterRequest{
						And: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndRequest{
							Prefix: "cc_",
							Tags: []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterAndTagsRequest{
								{
									Key:   "key2",
									Value: "value2",
								},
							},
						},
						Prefix: "ss_",
						Tag: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesFilterTagRequest{
							Key:   "key1",
							Value: "value1",
						},
					},
					Prefix: "nn_",
					AbortIncompleteMultipartUpload: &ZosPutBucketLifecycleConfLifecycleConfigurationRulesAbortIncompleteMultipartUploadRequest{
						DaysAfterInitiation: 123,
					},
					Transitions: []*ZosPutBucketLifecycleConfLifecycleConfigurationRulesTransitionsRequest{
						{
							Date:         "2022-10-18T00:00:00Z",
							Days:         123,
							StorageClass: "GLACIER",
						},
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
