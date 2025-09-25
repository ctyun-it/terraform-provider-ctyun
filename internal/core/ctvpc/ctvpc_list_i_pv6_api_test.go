package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListIPv6Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListIPv6Api

	// 构造请求
	var vpcID string = ""
	var subnetID string = ""
	var ipAddress string = ""
	request := &CtvpcListIPv6Request{
		RegionID:  "",
		VpcID:     &vpcID,
		SubnetID:  &subnetID,
		IpAddress: &ipAddress,
		Page:      1,
		PageNo:    1,
		PageSize:  1,
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
