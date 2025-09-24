package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcVpcCreateSubnetApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcVpcCreateSubnetApi

	// 构造请求
	var description string = ""
	var enableIpv6 bool = false
	var subnetGatewayIP string = "192.168.1.1"
	var subnetType string = "common"
	var dhcpIP string = "192.168.1.2"
	request := &CtvpcVpcCreateSubnetRequest{
		ClientToken:     "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:        "",
		VpcID:           "",
		Name:            "acl11",
		Description:     &description,
		CIDR:            "192.168.1.0/24",
		EnableIpv6:      &enableIpv6,
		DnsList:         []*string{},
		SubnetGatewayIP: &subnetGatewayIP,
		SubnetType:      &subnetType,
		DhcpIP:          &dhcpIP,
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
