package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsDeleteSecurityGroupV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsDeleteSecurityGroupV41Api

	// 构造请求
	request := &CtecsDeleteSecurityGroupV41Request{
		ClientToken:     "del-sg-test-01",
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		ProjectID:       "",
		SecurityGroupID: "sg-bp67axxxxzb4p",
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
