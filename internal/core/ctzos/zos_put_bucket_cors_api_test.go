package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestZosPutBucketCorsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ZosPutBucketCorsApi

	// 构造请求
	request := &ZosPutBucketCorsRequest{
		Bucket:   "bucket1",
		RegionID: "332232eb-63aa-465e-9028-52e5123866f0",
		CORSConfiguration: &ZosPutBucketCorsCORSConfigurationRequest{
			CORSRules: []*ZosPutBucketCorsCORSConfigurationCORSRulesRequest{
				{
					AllowedHeaders: []string{},
					AllowedMethods: []string{},
					AllowedOrigins: []string{},
					ExposeHeaders:  []string{},
					MaxAgeSeconds:  3000,
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
