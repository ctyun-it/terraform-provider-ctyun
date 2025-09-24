package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcAssociatePrivateZoneToVpcApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcAssociatePrivateZoneToVpcApi

	// 构造请求
	var clientToken string = "79fa97e3-c48b-xxxx-9f46-6a13d8163678"
	request := &CtvpcAssociatePrivateZoneToVpcRequest{
		ClientToken: &clientToken,
		RegionID:    "",
		ZoneID:      "zone-r5i4zghgvq",
		VpcIDList:   "vpc-tuid8d646e, vpc-tuid8d646e",
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
