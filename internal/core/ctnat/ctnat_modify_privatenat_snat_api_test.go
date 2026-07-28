package ctnat

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtnatModifyPrivatenatSnatApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtnatModifyPrivatenatSnatApi

	// 构造请求
	request := &CtnatModifyPrivatenatSnatRequest{
		RegionID:       "81f7728662dd11ec810800155d307d5b",
		SnatID:         "natgwsr-1o5sdqb7i2",
		SnatIps:        []string{},
		SourceSubnetID: "subnet-xxxxxx",
		Description:    "acl",
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
