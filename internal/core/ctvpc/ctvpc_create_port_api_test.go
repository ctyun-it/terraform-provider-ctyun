package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreatePortApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreatePortApi

	// 构造请求
	var primaryPrivateIp string = ""
	var name string = "acl11"
	var description string = ""
	request := &CtvpcCreatePortRequest{
		ClientToken:             "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:                "",
		SubnetID:                "",
		PrimaryPrivateIp:        &primaryPrivateIp,
		Ipv6Addresses:           []*string{},
		SecurityGroupIds:        []*string{},
		SecondaryPrivateIpCount: 1,
		SecondaryPrivateIps:     []*string{},
		Name:                    &name,
		Description:             &description,
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
