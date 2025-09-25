package ccse

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCcseUpdateClusterApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CcseUpdateClusterApi

	// 构造请求
	request := &CcseUpdateClusterRequest{
		ClusterId:       "47281b02f87757478f20b1827c97cadf",
		RegionId:        "bb9fdb42056f11eda1610242ac110002",
		ClusterDesc:     "",
		ClusterAlias:    "cce-demo",
		StartPort:       20106,
		EndPort:         32767,
		SecurityGroupId: "sg-dbd1fhzzzz",
		CustomSan: &CcseUpdateClusterCustomSanRequest{
			Action: "overwrite",
			Values: []string{},
		},
		Cubecni: &CcseUpdateClusterCubecniRequest{
			MinPoolSize:   0,
			MaxPoolSize:   5,
			AppendSubnets: []string{},
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
