package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtnatModifyPrivatenatDnatApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtnatModifyPrivatenatDnatApi

	// 构造请求
	request := &CtnatModifyPrivatenatDnatRequest{
		RegionID:     "81f7728662dd11ec810800155d307d5b",
		DnatID:       "natgwdr-1o5sdqb7i2",
		ExternalIP:   "192.168.2.3",
		ExternalPort: 1020,
		InternalIP:   "10.0.1.22",
		PortID:       "port-xxxxxx",
		InternalPort: 8812,
		Protocol:     "tcp",
		Description:  "acl",
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
