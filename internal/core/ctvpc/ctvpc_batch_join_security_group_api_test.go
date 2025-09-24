package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcBatchJoinSecurityGroupApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcBatchJoinSecurityGroupApi

	// 构造请求
	var networkInterfaceID string = "port-l0shxxxrfyg9"
	request := &CtvpcBatchJoinSecurityGroupRequest{
		RegionID:           "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		SecurityGroupIDs:   []string{},
		InstanceID:         "89a2d977-c078-5779-f391-f0ab8c9773b6",
		NetworkInterfaceID: &networkInterfaceID,
		Action:             "joinSecurityGroup",
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
