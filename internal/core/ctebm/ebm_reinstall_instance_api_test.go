package ctebm

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"testing"
)

func TestEbmReinstallInstanceApi_Do(t *testing.T) {
	// 初始化
	client := core.DefaultClient()
	credential := core.NewCredential("<YOUR_AK>", "<YOUR_SK>")
	// credential := core.CredentialFromEnv()
	apis := NewApis("<YOUR_ENDPOINT>", client)
	api := apis.EbmReinstallInstanceApi

	// 构造请求
	var systemVolumeRaidUUID string = "r-wtzluqacgzzxgunnabdkpnpjew3d"
	var dataVolumeRaidUUID string = "r-qytwf9r5h0yn9x4evjkyr0n1cwyb"
	var userData string = "ZWNobyBoZWxsbyBnb3N0YWNrIQ=="
	var keyName string = "my-keyname"
	request := &EbmReinstallInstanceRequest{
		RegionID:             "81f7728662dd11ec810800155d307d5b",
		AzName:               "az1",
		InstanceUUID:         "ss-9d4u1yd1jr0a3xeu59fq9svecxql",
		Hostname:             "host-pm-3301",
		Password:             "*************",
		ImageUUID:            "im-as6g7uju3cesx8n7qru8vqn2iqkf",
		SystemVolumeRaidUUID: &systemVolumeRaidUUID,
		DataVolumeRaidUUID:   &dataVolumeRaidUUID,
		RedoRaid:             false,
		UserData:             &userData,
		KeyName:              &keyName,
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
