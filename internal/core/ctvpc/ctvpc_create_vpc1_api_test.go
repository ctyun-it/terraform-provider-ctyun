package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcCreateVpc1Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcCreateVpc1Api

	// 构造请求
	var projectID string = "0"
	var enableIpv6 bool = false
	var ipv6SegmentPoolID string = ""
	request := &CtvpcCreateVpc1Request{
		ClientToken:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		ProjectID:         &projectID,
		RegionID:          "",
		Name:              "acl11",
		CIDR:              "192.168.0.0/16",
		EnableIpv6:        &enableIpv6,
		Ipv6SegmentPoolID: &ipv6SegmentPoolID,
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
