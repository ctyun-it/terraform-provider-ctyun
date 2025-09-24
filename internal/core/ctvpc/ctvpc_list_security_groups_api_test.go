package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListSecurityGroupsApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListSecurityGroupsApi

	// 构造请求
	var vpcID string = ""
	var queryContent string = ""
	var projectID string = "0"
	var instanceID string = ""
	request := &CtvpcListSecurityGroupsRequest{
		RegionID:     "",
		VpcID:        &vpcID,
		QueryContent: &queryContent,
		ProjectID:    &projectID,
		InstanceID:   &instanceID,
		PageNumber:   1,
		PageNo:       1,
		PageSize:     10,
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
