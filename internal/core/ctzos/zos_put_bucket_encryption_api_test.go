package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketEncryptionApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketEncryptionApi

	// 构造请求
	request := &ZosPutBucketEncryptionRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		ServerSideEncryptionConfiguration: &ZosPutBucketEncryptionServerSideEncryptionConfigurationRequest{
			Rules: []*ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesRequest{
				{
					ApplyServerSideEncryptionByDefault: &ZosPutBucketEncryptionServerSideEncryptionConfigurationRulesApplyServerSideEncryptionByDefaultRequest{
						SSEAlgorithm:   "AES256",
						KMSMasterKeyID: "5e6c8c96-1d90-4k07-9436-1485c8580ab0:::45d9efdad66f11ec9aab0242ac110002:207l1rh223254948a23d713442b278b1",
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
