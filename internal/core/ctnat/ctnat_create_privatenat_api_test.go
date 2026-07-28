package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtnatCreatePrivatenatApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtnatCreatePrivatenatApi

	// 构造请求
	var autoRenew bool = true
	request := &CtnatCreatePrivatenatRequest{
		RegionID:        "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		ProjectID:       "2256c561639d4ed9b9fb9009398914ad",
		VpcID:           "vpc-bp1xxxu",
		SubnetID:        "subnet-xxxxx",
		Spec:            "small",
		Name:            "acl11",
		Description:     "acl",
		ClientToken:     "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		CycleType:       "month",
		CycleCount:      3,
		AzName:          "az1",
		AutoRenew:       &autoRenew,
		PayVoucherPrice: "",
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
