package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateEndpointServiceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateEndpointServiceApi

	// 构造请求
	var rawType string = "interface"
	var instanceType string = "lb"
	var instanceID string = "lb-tfa20qb11w"
	var instanceID6 string = "havip-tfa20qb11w"
	var underlayIP string = "100.126.64.1"
	var underlayIP6 string = "100::7:0:0:647e:4000"
	var subnetID string = "subnet-7owr4do29a"
	var oaType string = "close"
	var serviceCharge bool = false
	var forceEnableDns bool = true
	var dnsName string = "Test"
	var reverseIsUnderlay bool = false
	var transitIP string = "192.168.0.1"
	var transitIP6 string = "100::7:0:0:647e:5665"
	request := &CtvpcCreateEndpointServiceRequest{
		ClientToken:    "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		RegionID:       "81f7728662dd11ec810800155d307d5b",
		VpcID:          "vpc-srsiebllhc",
		IpVersion:      2,
		RawType:        &rawType,
		Name:           "acl11",
		InstanceType:   &instanceType,
		InstanceID:     &instanceID,
		InstanceID6:    &instanceID6,
		UnderlayIP:     &underlayIP,
		UnderlayIP6:    &underlayIP6,
		SubnetID:       &subnetID,
		AutoConnection: false,
		Rules: []*CtvpcCreateEndpointServiceRulesRequest{
			{
				Protocol:     "TCP",
				ServerPort:   1,
				EndpointPort: 1,
			},
		},
		OaType:            &oaType,
		ServiceCharge:     &serviceCharge,
		ForceEnableDns:    &forceEnableDns,
		DnsName:           &dnsName,
		ReverseIsUnderlay: &reverseIsUnderlay,
		TransitIP:         &transitIP,
		TransitIP6:        &transitIP6,
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
