package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsTotalInstanceBackupVolumeSizeV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsTotalInstanceBackupVolumeSizeV41Api

	// 构造请求
	request := &CtecsTotalInstanceBackupVolumeSizeV41Request{
		RegionID:   "bb9fdb42056f11eda1610242ac110002",
		InstanceID: "69aac66c-78e8-e704-e6e1-311b3f40a278",
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
