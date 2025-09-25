package scaling

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestScalingConfigUpdateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.ScalingConfigUpdateApi

	// 构造请求
	var monitorService bool = true
	request := &ScalingConfigUpdateRequest{
		RegionID:            "81f7728662dd11ec810800155d307d5b",
		ConfigID:            427,
		Name:                "as-config-local001",
		ImageID:             "b78812b0-ff50-4816-b58f-5c4fbc230b08",
		SecurityGroupIDList: []string{},
		SpecName:            "c6.large.2",
		Volumes: []*ScalingConfigUpdateVolumesRequest{
			{
				VolumeType: "SATA",
				VolumeSize: 40,
				DiskMode:   "VBD",
				Flag:       1,
			},
		},
		UseFloatings: 2,
		BandWidth:    100,
		LoginMode:    2,
		Username:     "root",
		Password:     "ysdf12dfgGG@",
		KeyPairID:    "539b0666-d667-c71f-62b5-4db7a3cbdd59",
		Tags: []*ScalingConfigUpdateTagsRequest{
			{
				Key:   "key1",
				Value: "value1",
			},
		},
		AzNames:        []string{},
		MonitorService: &monitorService,
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
