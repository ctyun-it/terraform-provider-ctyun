package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsListInstanceV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsListInstanceV41Api

	// 构造请求
	var asc bool = true
	request := &CtecsListInstanceV41Request{
		RegionID:        "bb9fdb42056f11eda1610242ac110002",
		AzName:          "cn-huadong1-jsnj1A-public-ctcloud",
		ProjectID:       "0",
		PageNo:          1,
		PageSize:        10,
		State:           "active",
		Keyword:         "ecs-888",
		InstanceName:    "ecs-1",
		InstanceIDList:  "73f321ea-62ff-11ec-a8bc-005056898fe0,88f888ea-88ff-88ec-a8bc-888888888fe8",
		SecurityGroupID: "sg-tolywxbe1f",
		VpcID:           "vpc-euu7edo58k",
		ResourceID:      "9178e00c6fd148a88d4307950a9468df",
		LabelList: []*CtecsListInstanceV41LabelListRequest{
			{
				LabelKey:   "test-key",
				LabelValue: "test-value",
			},
		},
		Sort: "expiredTime",
		Asc:  &asc,
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
