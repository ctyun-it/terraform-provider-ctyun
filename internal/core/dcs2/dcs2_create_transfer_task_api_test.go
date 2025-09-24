package dcs2

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestDcs2CreateTransferTaskApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.Dcs2CreateTransferTaskApi

	// 构造请求
	var originalCluster bool = false
	var originalCluster1 bool = false
	request := &Dcs2CreateTransferTaskRequest{
		RegionId: "bb9fdb42056f11eda1610242ac110002",
		SourceDbInfo: &Dcs2CreateTransferTaskSourceDbInfoRequest{
			SpuInstId:       "",
			IpAddr:          "连接地址（ip:port）",
			OriginalCluster: &originalCluster,
			AccountName:     "",
			Password:        "",
		},
		TargetDbInfo: &Dcs2CreateTransferTaskTargetDbInfoRequest{
			SpuInstId:       "",
			IpAddr:          "连接地址（ip:port）",
			OriginalCluster: &originalCluster1,
			AccountName:     "",
			Password:        "",
		},
		SyncMode:     1,
		ConflictMode: 1,
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
