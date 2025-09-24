package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsCreateInstanceBackupPolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsCreateInstanceBackupPolicyApi

	// 构造请求
	request := &CtecsCreateInstanceBackupPolicyRequest{
		RegionID:      "bb9fdb42056f11eda1610242ac110002",
		PolicyName:    "test-bak",
		CycleType:     "day",
		CycleDay:      1,
		CycleWeek:     "0,2,6",
		Time:          "1,20",
		Status:        1,
		RetentionType: "date",
		RetentionDay:  30,
		RetentionNum:  20,
		ProjectID:     "0",
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
