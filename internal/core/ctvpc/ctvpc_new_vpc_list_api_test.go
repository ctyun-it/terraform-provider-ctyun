package ctvpc

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtvpcNewVpcListApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtvpcNewVpcListApi

	// 构造请求
	var vpcID string = "vpc-hfw53u96ku,vpc-hfw53u54du"
	var vpcName string = "test"
	var projectID string = "0"
	request := &CtvpcNewVpcListRequest{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		VpcID:      &vpcID,
		VpcName:    &vpcName,
		PageNumber: 1,
		PageNo:     1,
		PageSize:   5,
		ProjectID:  &projectID,
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
