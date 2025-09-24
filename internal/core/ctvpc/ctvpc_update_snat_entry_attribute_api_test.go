package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcUpdateSnatEntryAttributeApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcUpdateSnatEntryAttributeApi

	// 构造请求
	var sourceSubnetID string = "5fe30709-93ef-522f-a1a0-d8c8f6803e0d"
	var sourceCIDR string = "10.1.1.0/24"
	var description string = ",./;'[]·~！@#￥%……&*（） ——-+={}"
	request := &CtvpcUpdateSnatEntryAttributeRequest{
		RegionID:       "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
		SNatID:         "ngwsr-0mghyey7ar",
		SourceSubnetID: &sourceSubnetID,
		SourceCIDR:     &sourceCIDR,
		Description:    &description,
		ClientToken:    "79fa97e3-c48b-xxxx-9f46-6a13d8163678",
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
