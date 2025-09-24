package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateDnatEntryApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateDnatEntryApi

	// 构造请求
	var virtualMachineID string = "62dab755-1f54-b961-307e-4a0c84274a09"
	var internalIp string = "0.0.0.0"
	var serverType string = "VM"
	var description string = "test"
	request := &CtvpcCreateDnatEntryRequest{
		RegionID:           "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		NatGatewayID:       "nat 网关 ID",
		ExternalID:         "eip-3ubjeszr3a",
		ExternalPort:       80,
		VirtualMachineID:   &virtualMachineID,
		VirtualMachineType: 1,
		InternalIp:         &internalIp,
		ServerType:         &serverType,
		InternalPort:       80,
		Protocol:           "tcp",
		ClientToken:        "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		Description:        &description,
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
