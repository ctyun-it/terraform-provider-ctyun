package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2DescribeInstancesApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2DescribeInstancesApi

	// 构造请求
	request := &Dcs2DescribeInstancesRequest{
		RegionId:     "200000001852",
		PageIndex:    1,
		PageSize:     10,
		InstanceName: "tyyTest-HrpRun",
		LabelIds:     []string{},
		ProjectId:    "0",
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
