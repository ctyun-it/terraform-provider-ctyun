package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcVpcListPortApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcVpcListPortApi

	// 构造请求
	var vpcID string = "vpc-r5i4zghgvq"
	var deviceID string = "b96afb7b971660db375d0963ad10c192"
	var subnetID string = "subnet-r5i4zghgvq"
	request := &CtvpcVpcListPortRequest{
		RegionID:   "81f7728662dd11ec810800155d307d5b",
		VpcID:      &vpcID,
		DeviceID:   &deviceID,
		SubnetID:   &subnetID,
		PageNumber: 1,
		PageNo:     1,
		PageSize:   10,
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
