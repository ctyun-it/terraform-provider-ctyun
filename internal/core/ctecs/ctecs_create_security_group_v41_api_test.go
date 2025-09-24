package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateSecurityGroupV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateSecurityGroupV41Api

	// 构造请求
	request := &CtecsCreateSecurityGroupV41Request{
		ClientToken: "create-sg-test-01",
		RegionID:    "bb9fdb42056f11eda1610242ac110002",
		ProjectID:   "6732237e53bc4591b0e67d750030ebe3",
		VpcID:       "4797e8a1-722d-4996-9362-458001813e41",
		Name:        "sg-bp67axxxxzb4p",
		Description: "acl",
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
