package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceBackupPolicyBindInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceBackupPolicyBindInstancesApi

	// 构造请求
	request := &CtecsInstanceBackupPolicyBindInstancesRequest{
		RegionID:       "bb9fdb42056f11eda1610242ac110002",
		PolicyID:       "e14f067ea2d111ed98fd0242ac110007",
		InstanceIDList: "73f321ea-62ff-11ec-a8bc-005056898fe0,285010af-16f1-137e-06c0-920d4bdd0026",
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
