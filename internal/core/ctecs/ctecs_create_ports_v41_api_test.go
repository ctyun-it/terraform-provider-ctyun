package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreatePortsV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreatePortsV41Api

	// 构造请求
	request := &CtecsCreatePortsV41Request{
		ClientToken:             "ports_create07231529",
		RegionID:                "bb9fdb42056f11eda1610242ac110002",
		SubnetID:                "subnet-y8cofge5uj",
		PrimaryPrivateIp:        "172.16.0.141",
		Ipv6Addresses:           []string{},
		SecurityGroupIds:        []string{},
		SecondaryPrivateIpCount: 1,
		SecondaryPrivateIps:     []string{},
		Name:                    "nic-test01",
		Description:             "new-Description",
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
