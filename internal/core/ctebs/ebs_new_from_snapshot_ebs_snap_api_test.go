package ctebs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbsNewFromSnapshotEbsSnapApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbsNewFromSnapshotEbsSnapApi

	// 构造请求
	var clientToken string = "cbe3840c-bda4-4102-b68f-98c9d7190d69"
	var multiAttach bool = false
	var projectID string = "0"
	var onDemand bool = true
	var cycleType string = "month"
	request := &EbsNewFromSnapshotEbsSnapRequest{
		SnapshotID:  "3f868846-f47f-4619-a5b4-a02e9714f744",
		ClientToken: &clientToken,
		RegionID:    "41f64827f25f468595ffa3a5deb5d15d",
		MultiAttach: &multiAttach,
		ProjectID:   &projectID,
		DiskMode:    "VBD",
		DiskName:    "mydisk-0001",
		DiskSize:    1024,
		OnDemand:    &onDemand,
		CycleType:   &cycleType,
		CycleCount:  2,
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
