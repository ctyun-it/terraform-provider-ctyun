package ctebsbackup

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsbackupCreateEbsBackupPolicyApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsbackupCreateEbsBackupPolicyApi

	// 构造请求
	var remainFirstOfCurMonth bool = false
	var advRetentionStatus bool = true
	request := &EbsbackupCreateEbsBackupPolicyRequest{
		RegionID:              "81f7728662dd11ec810800155d307d5b",
		PolicyName:            "test-policy",
		Status:                0,
		CycleType:             "day",
		CycleDay:              1,
		CycleWeek:             "0,2,6",
		Time:                  "1,20",
		RetentionType:         "num",
		RetentionNum:          1,
		RetentionDay:          1,
		RemainFirstOfCurMonth: &remainFirstOfCurMonth,
		ProjectID:             "0",
		FullBackupInterval:    1,
		AdvRetentionStatus:    &advRetentionStatus,
		AdvRetention: &EbsbackupCreateEbsBackupPolicyAdvRetentionRequest{
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
