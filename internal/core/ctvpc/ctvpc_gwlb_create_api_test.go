package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcGwlbCreateApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcGwlbCreateApi

	// 构造请求
	var projectID string = ""
	var privateIpAddress string = ""
	var ipv6Address string = ""
	var deleteProtection bool = false
	var ipv6Enabled bool = false
	var payVoucherPrice string = ""
	request := &CtvpcGwlbCreateRequest{
		RegionID:         "81f7728662dd11ec810800155d307d5b",
		ClientToken:      "fewfrf",
		ProjectID:        &projectID,
		SubnetID:         "",
		Name:             "acl11",
		PrivateIpAddress: &privateIpAddress,
		Ipv6Address:      &ipv6Address,
		DeleteProtection: &deleteProtection,
		Ipv6Enabled:      &ipv6Enabled,
		CycleType:        "on_demand",
		PayVoucherPrice:  &payVoucherPrice,
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
