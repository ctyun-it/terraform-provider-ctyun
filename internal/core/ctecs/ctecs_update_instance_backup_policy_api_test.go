package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateInstanceBackupPolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateInstanceBackupPolicyApi

	// 构造请求
	var totalBackup bool = false
	var advRetentionStatus bool = false
	request := &CtecsUpdateInstanceBackupPolicyRequest{
		RegionID:           "bb9fdb42056f11eda1610242ac110002",
		PolicyID:           "d58bc64aa3b411edaf600242ac110009",
		PolicyName:         "test-bak",
		CycleType:          "day",
		CycleDay:           1,
		CycleWeek:          "0,2,6",
		Time:               "1,20",
		Status:             1,
		RetentionType:      "date",
		RetentionDay:       30,
		RetentionNum:       20,
		TotalBackup:        &totalBackup,
		AdvRetentionStatus: &advRetentionStatus,
		AdvRetention: &CtecsUpdateInstanceBackupPolicyAdvRetentionRequest{
			AdvDay:   1,
			AdvWeek:  1,
			AdvMonth: 1,
			AdvYear:  1,
		},
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
