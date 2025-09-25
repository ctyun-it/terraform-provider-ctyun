package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsInstanceBackupPolicyBindRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsInstanceBackupPolicyBindRepoApi

	// 构造请求
	request := &CtecsInstanceBackupPolicyBindRepoRequest{
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		RepositoryID: "62a8c714-7f2a-4efb-b0b9-5542694c9c24",
		PolicyID:     "c0f853a4a5c311edb27d0242ac110007",
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
