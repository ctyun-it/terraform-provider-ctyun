package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsDeleteInstanceBackupRepoApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsDeleteInstanceBackupRepoApi

	// 构造请求
	request := &CtecsDeleteInstanceBackupRepoRequest{
		RegionID:     "bb9fdb42056f11eda1610242ac110002",
		RepositoryID: "829e9ef0-7805-492d-988f-de4b4e705339",
		ClientToken:  "4cf2962d-e92c-4c00-9181-cfbb2218636c",
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
