package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2CreateInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2CreateInstanceApi

	// 构造请求
	var autoPay bool = true
	var autoRenew bool = false
	request := &Dcs2CreateInstanceRequest{
		RegionId:          "b342b77ef26b11ecb0ac0242ac110002",
		ProjectID:         "0",
		AutoPay:           &autoPay,
		Period:            1,
		ChargeType:        "PrePaid",
		ZoneName:          "cn-xinan1-1A",
		SecondaryZoneName: "cn-xinan1-2A",
		EngineVersion:     "6.0",
		Version:           "BASIC",
		Edition:           "DirectCluster",
		HostType:          "S",
		DataDiskType:      "SAS",
		ShardMemSize:      "1",
		ShardCount:        3,
		Capacity:          "3",
		CopiesCount:       2,
		InstanceName:      "Open-Go-SDK-1600",
		VpcId:             "vpc-olzdlcv7f8",
		SubnetId:          "subnet-4ic7xre126",
		Secgroups:         "sg-whg3m7t06m",
		Password:          "testABC@158",
		AutoRenew:         &autoRenew,
		Size:              1,
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
