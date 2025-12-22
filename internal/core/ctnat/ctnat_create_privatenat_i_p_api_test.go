package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtnatCreatePrivatenatIPApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtnatCreatePrivatenatIPApi

	// 构造请求
	request := &CtnatCreatePrivatenatIPRequest{
		RegionID:     "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		NatGatewayID: "natgw-aldjlfsfd",
		Address:      "IP地址，必须在中转网段指定的网络范围内",
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
