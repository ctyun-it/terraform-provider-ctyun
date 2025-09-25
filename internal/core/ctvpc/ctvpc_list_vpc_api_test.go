package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcListVpcApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcListVpcApi

	// 构造请求
	var projectID string = "0"
	var vpcID string = "vpc-hfw53u96ku,vpc-hfw53u54du"
	var vpcName string = "test"
	request := &CtvpcListVpcRequest{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		ProjectID:  &projectID,
		VpcID:      &vpcID,
		VpcName:    &vpcName,
		PageNumber: 1,
		PageNo:     1,
		PageSize:   5,
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
