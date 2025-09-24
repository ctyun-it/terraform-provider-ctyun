package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsCreatePolicyEbsSnapApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsCreatePolicyEbsSnapApi

	// 构造请求
	var isEnabled bool = true
	var projectID string = "0"
	request := &EbsCreatePolicyEbsSnapRequest{
		RegionID:           "41f64827f25f468595ffa3a5deb5d15d",
		SnapshotPolicyName: "policy-1",
		RepeatWeekdays:     "0,1,2",
		RepeatTimes:        "0,1,2",
		RetentionTime:      2,
		IsEnabled:          &isEnabled,
		ProjectID:          &projectID,
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
