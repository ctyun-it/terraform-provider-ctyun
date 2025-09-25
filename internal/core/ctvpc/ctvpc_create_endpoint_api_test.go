package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateEndpointApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateEndpointApi

	// 构造请求
	var iP string = "192.168.1.1"
	var iP6 string = "100:1:126:d400:d471:7505:fe0:596b"
	var enableDns bool = false
	var payVoucherPrice string = ""
	var deleteProtection bool = false
	var protectionService string = ""
	request := &CtvpcCreateEndpointRequest{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:          "81f7728662dd11ec810800155d307d5b",
		CycleType:         "",
		EndpointServiceID: "",
		IpVersion:         2,
		EndpointName:      "acl11",
		SubnetID:          "subnet-xxxx",
		VpcID:             "",
		IP:                &iP,
		IP6:               &iP6,
		WhitelistFlag:     1,
		Whitelist:         []*string{},
		Whitelist6:        []*string{},
		EnableDns:         &enableDns,
		PayVoucherPrice:   &payVoucherPrice,
		DeleteProtection:  &deleteProtection,
		ProtectionService: &protectionService,
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
