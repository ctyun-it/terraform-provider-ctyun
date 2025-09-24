package ctecs

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestCtecsUpdateVolumeV41Api_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.CtecsUpdateVolumeV41Api

	// 构造请求
	request := &CtecsUpdateVolumeV41Request{
		RegionID: "bb9fdb42056f11eda1610242ac110002",
		DiskName: "磁盘名字，长度限制2~63字符，不支持中文",
		DiskID:   "eff436e3d44040f1b306ab3a14530f02",
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
