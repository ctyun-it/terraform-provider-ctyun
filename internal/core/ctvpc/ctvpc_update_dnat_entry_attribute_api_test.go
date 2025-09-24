package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateDnatEntryAttributeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateDnatEntryAttributeApi

	// 构造请求
	var externalID string = "eip-3ubjeszr3a"
	var virtualMachineID string = "62dab755-1f54-b961-307e-4a0c84274a09"
	var internalIp string = "0.0.0.0"
	var description string = "test"
	var serverType string = ""
	request := &CtvpcUpdateDnatEntryAttributeRequest{
		RegionID:           "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		DNatID:             "ngwdr-jqunxiapy9",
		ExternalID:         &externalID,
		ExternalPort:       80,
		VirtualMachineID:   &virtualMachineID,
		VirtualMachineType: 1,
		InternalIp:         &internalIp,
		InternalPort:       80,
		Protocol:           "tcp",
		ClientToken:        "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		Description:        &description,
		ServerType:         &serverType,
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
