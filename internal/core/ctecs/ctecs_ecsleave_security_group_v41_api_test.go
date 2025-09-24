package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsEcsleaveSecurityGroupV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsEcsleaveSecurityGroupV41Api

	// 构造请求
	request := &CtecsEcsleaveSecurityGroupV41Request{
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		SecurityGroupID: "sg-bp67axxxxzb4p",
		InstanceID:      "89a2d977-c078-5779-f391-f0ab8c9773b6",
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
