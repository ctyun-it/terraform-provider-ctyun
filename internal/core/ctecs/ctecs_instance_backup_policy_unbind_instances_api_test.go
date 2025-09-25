package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceBackupPolicyUnbindInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceBackupPolicyUnbindInstancesApi

	// 构造请求
	request := &CtecsInstanceBackupPolicyUnbindInstancesRequest{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		PolicyID:       "e14f067ea2d111ed98fd0242ac110007",
		InstanceIDList: "73f321ea-62ff-11ec-a8bc-005056898fe0,543074a6-75a0-5e31-2c77-1501e3354f69",
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
